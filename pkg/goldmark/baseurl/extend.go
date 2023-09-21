// Package replacer is a extension for the goldmark
// (http://github.com/yuin/goldmark).
//
// This extension adds support for authomaticaly replacing text in markdowns.
package baseurl

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

var extender *Extender

type Extender struct {
	BaseURL        string
	PublicFilePath string
}

// New return initialized image render with source url replacing support.
func NewExtender(baseURL string, publicFilePath string) goldmark.Extender {
	extender = &Extender{
		BaseURL:        baseURL,
		PublicFilePath: publicFilePath,
	}
	return extender
}

func (e *Extender) Extend(m goldmark.Markdown) {
	// m.Parser().AddOptions(
	// 	parser.WithASTTransformers(
	// util.Prioritized(NewTransformer(), 500),
	// ),
	// )
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(NewTransformer(), 500),
		),
	)
	/*
		m.Renderer().AddOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(NewRenderer(), 500),
			),
		)
	*/
}
