package landingparser

import (
	"github.com/ITLab-Projects/pkg/models/landing"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

type LandingParser struct {
	parser *parser.Parser
}

func New() *LandingParser {
	lp := &LandingParser{
		parser: parser.New(),
	}

	return lp
}

func (lp *LandingParser) Parse(data []byte) *landing.Landing {
	if lp.parser == nil {
		lp.parser = parser.New()
	}
	
	node := lp.parser.Parse(data)
	
	l := &landing.Landing{
	}

	for _, child := range node.GetChildren() {
		switch n := child.(type) {
		case *ast.Heading:
			if childs := n.GetChildren(); len(childs) != 0 {
				switch string(childs[0].(*ast.Text).Literal) {
				case "Title":
					p := ast.GetNextNode(child).(*ast.Paragraph)
					text := ast.GetFirstChild(p).(*ast.Text)
					if text != nil {
						l.Title = string(text.Literal)
					}
				case "Description":
					next := ast.GetNextNode(child)
					switch node := next.(type) {
					case *ast.CodeBlock:
						l.Description = string(node.Literal)
					case *ast.Paragraph:
						child := ast.GetFirstChild(node).(*ast.Text)
						l.Description = string(child.Literal)
					}
				case "Images":
					l.Image = getImagesURL(ast.GetNextNode(child))
				case "Videos":
					l.Videos = getLinkList(ast.GetNextNode(child))
				case "Tags":
					l.Tags = getStringList(ast.GetNextNode(child))
				case "Tech":
					l.Tech = getStringList(ast.GetNextNode(child))
				case "Developers":
					l.Developers = getStringList(ast.GetNextNode(child))
				case "Site":
					l.Site = landing.Site(getLink(ast.GetNextNode(child)))
				case "SourceCode":
					l.SourceCode = getSourceCode(ast.GetNextNode(child))
				}
				
			}
		}
	}

	return l
}

func getImagesURL(
	node ast.Node,
) []string {
	var urls []string

	ast.WalkFunc(
		node,
		func(node ast.Node, entering bool) ast.WalkStatus {
			if !entering {
				return ast.GoToNext
			}

			switch n := node.(type) {
			case *ast.Image:
				urls = append(urls, string(n.Destination))
				return ast.SkipChildren
			}

			return ast.GoToNext
		},
	)

	return urls
}

func getLinkList(
	node ast.Node,
) []string {
	links := make([]string, 0)

	ast.WalkFunc(
		node,
		func(node ast.Node, entering bool) ast.WalkStatus {
			if !entering {
				return ast.GoToNext
			}

			switch n := node.(type) {
			case *ast.Link:
				links = append(links, string(n.Destination))
			}

			return ast.GoToNext
		},
	)

	return links
}

func getStringList(
	node ast.Node,
) []string {
	var list []string

	ast.WalkFunc(
		node,
		func(node ast.Node, entering bool) ast.WalkStatus {
			if !entering {
				return ast.GoToNext
			}

			switch n := node.(type) {
			case *ast.Text:
				list = append(list, string(n.Literal))
			}

			return ast.GoToNext
		},
	)

	return list
}

func getLink(
	node ast.Node,
) string {
	var link string

	ast.WalkFunc(
		node,
		func(node ast.Node, entering bool) ast.WalkStatus {
			if !entering {
				return ast.GoToNext
			}

			switch n := node.(type) {
			case *ast.Link:
				link = string(n.Destination)
				return ast.Terminate
			}

			return ast.GoToNext
		},
	)

	return link
}

func getText(
	node ast.Node,
) string {
	var text string

	ast.WalkFunc(
		node,
		func(node ast.Node, entering bool) ast.WalkStatus {
			if !entering {
				return ast.GoToNext
			}

			switch n := node.(type) {
			case *ast.Text:
				text = string(n.Literal)
				return ast.Terminate
			}

			return ast.GoToNext
		},
	)

	return text
}

func getSourceCode(
	node ast.Node,
) []*landing.SourceCode {
	var sourceCode []*landing.SourceCode

	ast.WalkFunc(
		node,
		func(node ast.Node, entering bool) ast.WalkStatus {
			if !entering {
				return ast.GoToNext
			}

			switch n := node.(type) {
			case *ast.TableHeader:
				return ast.SkipChildren
			case *ast.TableRow:
				s := &landing.SourceCode{}
				childs := n.Children
				if len(childs) != 2 {
					return ast.GoToNext
				}
				s.Name = getText(childs[0])
				s.Value = getLink(childs[1])

				sourceCode = append(sourceCode, s)
			}

			return ast.GoToNext
		},
	)

	return sourceCode
}