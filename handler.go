package twig

import (
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
)

// HandlerFunc Twig的Handler方法
type HandlerFunc func(Ctx) error

// Mount Mount当前Handler到注册器
func (h HandlerFunc) Mount(reg Register, method, path string, m ...MiddlewareFunc) Router {
	return reg.AddHandler(method, path, h, m...)
}

type MiddlewareFunc func(HandlerFunc) HandlerFunc

func (m MiddlewareFunc) UsedBy(reg Register) {
	reg.Use(m)
}

// WrapHttpHandler 包装http.Handler 为HandlerFunc
func WrapHttpHandler(h http.Handler) HandlerFunc {
	return func(c Ctx) error {
		h.ServeHTTP(c.Resp(), c.Req())
		return nil
	}
}

// Merge 中间件包装器
func Merge(handler HandlerFunc, m []MiddlewareFunc) HandlerFunc {
	if m == nil {
		return handler
	}
	h := handler
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

// NotFoundHandler 全局404处理方法
var NotFoundHandler = func(c Ctx) error {
	return ErrNotFound
}

// MethodNotAllowedHandler 全局405处理方法
var MethodNotAllowedHandler = func(c Ctx) error {
	return ErrMethodNotAllowed
}

// 获取handler的名称
func HandlerName(h HandlerFunc) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}

// Static 处理静态文件的HandlerFunc
func Static(r string) HandlerFunc {
	root := path.Clean(r)
	return func(c Ctx) error {
		p, err := url.PathUnescape(c.Param("*"))
		if err != nil {
			return err
		}
		name := filepath.Join(root, path.Clean("/"+p)) // 安全考虑 + "/"
		return c.File(name)
	}
}

// ServerInfo ServerInfo 中间件将Twig#Name()设置 Server 头
// Debug状态下，返回 x-powerd-by 为Twig#Type()
func ServerInfo() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Ctx) error {
			w := c.Resp()
			w.Header().Set(HeaderServer, c.Twig().Name())
			if c.Twig().Debug {
				w.Header().Set(HeaderXPoweredBy, c.Twig().Type())
			}
			return next(c)
		}
	}
}
