package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type logOptions struct {
	enableLog bool            // 是否开启日志
	extraPath map[string]bool // 日志开启后，不需要打印日志的URL path
}

type initOptions struct {
	CheckWhiteUIDS bool
}

// OptionItem 其他选项
type OptionItem func(*initOptions)

// LogOption metrics middleware日志选项
type LogOption func(options *logOptions)

var OptionOther initOptions

// Recovery panic拦截
func Recovery(c *gin.Context, recovered interface{}) {
	if recovered != nil {
		path, raw := c.Request.URL.Path, c.Request.URL.RawQuery
		if len(raw) != 0 {
			path = path + "?" + raw
		}
		s := fmt.Sprintf("%v", recovered)
		_ = c.Error(errors.New(s))
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}

// NewGinServe 创建gin engine
func NewGinServe(args ...LogOption) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.CustomRecovery(Recovery))
	return router
}

// NewGinServeOther 创建gin engine
func NewGinServeOther(option OptionItem, args ...LogOption) *gin.Engine {
	option(&OptionOther)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.CustomRecovery(Recovery))
	return router
}
