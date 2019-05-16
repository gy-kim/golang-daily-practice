package advantages

import "io"

func NewGeneratorV2(template io.Reader) *Generator {
	return &Generator{
		template: template,
	}
}

func (g *Generator) getStorge() Storage {
	if g.storage == nil {
		g.storage = &DefaultStorage{}
	}
	return g.storage
}

func (g *Generator) getRenderer() Renderer {
	if g.renderer == nil {
		g.renderer = &DefaultRenderer{}
	}
	return g.renderer
}

type DefaultStorage struct{}

func (d *DefaultStorage) Load() []interface{} {
	return nil
}

type DefaultRenderer struct{}

func (d *DefaultRenderer) Render(template io.Reader, params ...interface{}) []byte {
	return nil
}
