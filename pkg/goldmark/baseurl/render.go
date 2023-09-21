package baseurl

import (
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type Transformer struct {
}

func NewTransformer() parser.ASTTransformer {
	return &Transformer{}
}

func (c *Transformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(n ast.Node, entering bool) (status ast.WalkStatus, err error) {
		if image, ok := n.(*ast.Image); ok {
			if entering {
				image.SetAttributeString("class", []byte("img-content"))
				src := util.BytesToReadOnlyString(image.Destination)

				re := regexp.MustCompile(`(\.\.\/)*\.data/assets/`)
				src = re.ReplaceAllString(src, extender.BaseURL+extender.PublicFilePath)

				image.Destination = util.StringToReadOnlyBytes(src)

			}
		}
		return ast.WalkContinue, nil
	})
}

// Renderer struct is a renderer.NodeRenderer implementation for the extension.
type Renderer struct {
	html.Config
}

// NewRenderer builds a new Renderer with given options and returns it.
func NewRenderer() renderer.NodeRenderer {
	return &Renderer{
		Config: html.NewConfig(),
	}
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs interface.
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	//reg.Register(ast.KindText, r.renderText)
	//reg.Register(ast.KindString, r.renderString)
	//reg.Register(ast.KindImage, r.renderImage)
}

func (r *Renderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	_, _ = w.WriteString("<img src=\"")
	if r.Unsafe || !html.IsDangerousURL(n.Destination) {

		_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
	}
	_, _ = w.WriteString(`" alt="`)
	_, _ = w.Write(n.Text(source))
	_ = w.WriteByte('"')
	if n.Title != nil {
		_, _ = w.WriteString(` title="`)
		r.Writer.Write(w, n.Title)
		_ = w.WriteByte('"')
	}
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ImageAttributeFilter)
	}
	if r.XHTML {
		_, _ = w.WriteString(" />")
	} else {
		_, _ = w.WriteString(">")
	}
	return ast.WalkSkipChildren, nil

}
