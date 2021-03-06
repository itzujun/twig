package twig

import (
	"io"

	"github.com/twiglab/twig/internal/uuid"
)

// Plugger 定义了Twig的外部插件
// 如果插件需要生命周期管理，请实现Cycler接口
// 如果插件需要访问Twig本身，请实现Attacher接口
type Plugger interface {
	ID() string
}

// GetPlugin 从当前Ctx中获取Plugin
func GetPlugger(id string, c Ctx) (p Plugger, ok bool) {
	t := c.Twig()
	p, ok = t.GetPlugger(id)
	return
}

// Binder 数据绑定接口
// Binder 作为一个插件集成到Twig中,请实现Plugin接口
type Binder interface {
	Bind(interface{}, Ctx) error
}

// GetBinder 获取绑定接口
func GetBinder(id string, c Ctx) (binder Binder, ok bool) {
	var plugger Plugger
	if plugger, ok = GetPlugger(id, c); ok {
		binder, ok = plugger.(Binder)
	}
	return
}

type Renderer interface {
	Render(io.Writer, string, interface{}, Ctx) error
}

func GetRenderer(id string, c Ctx) (r Renderer, ok bool) {
	var plugger Plugger
	if plugger, ok = GetPlugger(id, c); ok {
		r, ok = plugger.(Renderer)
	}
	return
}

// IdGenerator ID发生器接口
type IdGenerator interface {
	NextID() string
}

func GetIdGenerator(id string, c Ctx) (gen IdGenerator, ok bool) {
	var plugger Plugger
	if plugger, ok = GetPlugger(id, c); ok {
		gen, ok = plugger.(IdGenerator)
	}
	return
}

const uuidPluginID = "_twig_uuid_plugin_id_"

type uuidGen struct {
}

func (id uuidGen) ID() string {
	return uuidPluginID
}

func (id uuidGen) NextID() string {
	return uuid.NewV1().String32()
}

func GenID(c Ctx) string {
	idgen, _ := GetIdGenerator(uuidPluginID, c)
	return idgen.NextID()
}
