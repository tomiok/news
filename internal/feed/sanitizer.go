package feed

import (
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// Sanitizer is the key transformation for all the content that will be saved in the database.
type Sanitizer interface {
	// Apply will execute the callbacks created in the struct. Returns in the same order as the args (title, description, content).
	Apply(title, desc, content string) (string, string, string)
}

type stdSanitize struct {
	htmlScape func(s string) string
	Sanitize  func(s string) string
	TrimSpace func(s string) string
}

func newSanitizer() *stdSanitize {
	p := bluemonday.NewPolicy()
	p.AllowElements([]string{"h1", "h2", "h3", "h4", "h5", "h6", "div", "span", "hr", "br", "b", "i", "strong", "em", "ol", "ul", "li", "pre", "code", "blockquote", "article", "section"}...)
	return &stdSanitize{
		htmlScape: html.UnescapeString,
		Sanitize:  p.Sanitize,
		TrimSpace: strings.TrimSpace,
	}
}

func (s *stdSanitize) Apply(title, desc, content string) (string, string, string) {
	title = s.apply(title)
	desc = s.apply(desc)
	content = s.apply(content)

	return title, desc, content
}

func (s *stdSanitize) apply(str string) string {
	if str == "" {
		return ""
	}

	str = s.TrimSpace(str)
	str = s.Sanitize(str)
	str = s.htmlScape(str)

	return str
}
