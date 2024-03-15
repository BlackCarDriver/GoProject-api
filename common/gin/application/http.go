package application

import (
	"context"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"time"
)

var log = logs.GetBeeLogger()

func HttpServe(host, port string) error {
	var errLog = logs.GetLogger()
	httpServer := http.Server{Addr: net.JoinHostPort(host, port), ErrorLog: errLog}
	return startHttpServer(&httpServer, nil)
}

func startHttpServer(httpServer *http.Server, listener net.Listener) (err error) {
	AtExit(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()
		err := httpServer.Shutdown(ctx)
		log.Info("http server byebye: addr=%s err=%v", httpServer.Addr, err)
	})

	if listener == nil {
		err = httpServer.ListenAndServe()
	} else {
		httpServer.Serve(listener)
	}

	if err == http.ErrServerClosed {
		log.Info("http wait shutdown: addr=%s err=%v ", httpServer.Addr, err)
		select { // fix: wait Shutdown
		}
	} else if err != nil {
		log.Emergency("http server crashed: addr=%s err=%v ", httpServer.Addr, err)
	}
	return err
}

// HTTPGinServe 启动gin http server
func HTTPGinServe(host, port string, engine *gin.Engine, listener net.Listener) error {
	return startHttpServer(&http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: engine,
	}, listener)
}
