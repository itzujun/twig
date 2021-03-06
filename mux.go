package twig

import (
	"net/http"
)

// Register 接口
// 实现路由注册
type Register interface {
	AddHandler(string, string, HandlerFunc, ...MiddlewareFunc) Router
	Use(...MiddlewareFunc)
}

// Lookuper 接口
// 实现路由查找
type Lookuper interface {
	Lookup(string, string, *http.Request) MuxerCtx
}

// Muxer Muxer 描述一个具体的路由
type Muxer interface {
	Lookuper
	Register
}

// Router 接口，Route接口用于描述一个已经加入Register的路由，由Register的AddHandler方法返回
// Router 提供命名路由的方法，被命名的路由可以用于Ctx.URL方法查找
type Router interface {
	Identifier
	Method() string
	Path() string
	Namer
}

type NamedRoute struct {
	N string // name
	P string // path
	M string // method
}

// ID 命名路由的ID
func (r *NamedRoute) ID() string {
	return r.M + r.P
}

// Name 命名路由的名称
func (r *NamedRoute) Name() string {
	return r.N
}

// Method 命名路由的Http方法
func (r *NamedRoute) Method() string {
	return r.M
}

// Path 命名路由的注册路径
func (r *NamedRoute) Path() string {
	return r.P
}

// SetName  Namer#SetName
func (r *NamedRoute) SetName(name string) {
	r.N = name
}

// Type 类型
func (r *NamedRoute) Type() string {
	return "handler"
}
