package web

import (
	"bytes"
	"context"
	"html/template"
)

type TemplateEngine  interface {
	// Render 渲染页面
	// tplName 模板名称, 按名索引
	// data 渲染数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)

	// 不需要，让具体实现自己管自己的模版
	// AddTemplate(tplName string, tpl []string)
}

type GoTemplateEngine struct {
	// 模版引擎
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	bs := bytes.NewBuffer(nil)
	err := g.T.ExecuteTemplate(bs, tplName, data)
	if err != nil {
		return nil, err
	}
	return bs.Bytes(), nil
}