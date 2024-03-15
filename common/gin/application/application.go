package application

import (
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var app = struct {
	shutdownSig  chan os.Signal
	exitHandlers []func()
	sync.Mutex
}{}

func init() {}

func process() {
	defer signal.Stop(app.shutdownSig)

	for {
		select {
		case sig := <-app.shutdownSig:
			log.Info("receive signal: sgi=%v handler=%v", sig, len(app.exitHandlers))
			var wg sync.WaitGroup
			wg.Add(len(app.exitHandlers))

			for _, h := range app.exitHandlers {
				go func(fn func()) {
					fn()
					wg.Done()
				}(h)
			}

			allDone := make(chan struct{})
			go func() {
				wg.Wait()
				allDone <- struct{}{}
			}()

			select {
			case <-time.After(6 * time.Second):
				log.Info("force to shutdown due to waitting timeout, byebye")
			case <-allDone:
				log.Info("all 'atexit' handlers done gracefully, byebye")
			}
			os.Exit(0)
		}
	}
}

func lazyInit() {
	app.Lock()
	if app.shutdownSig == nil {
		log.Info("all 'atexit' init...")
		app.shutdownSig = make(chan os.Signal, 16)
		signal.Notify(app.shutdownSig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go process()
	}
	app.Unlock()
}

func AtExit(handler func()) {
	log.Info("lazy init: name=%s", runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
	lazyInit() // 用到时才初始化

	app.Lock()
	app.exitHandlers = append(app.exitHandlers, handler)
	app.Unlock()
}
