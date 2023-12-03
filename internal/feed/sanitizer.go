package feed

import (
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// Sanitizer is the key transformation for all the content that will be saved in the database.
type Sanitizer interface {
	// Apply will execute the callbacks created in the struct. Returns in the same order as the args (title, description, content).
	Apply(title, desc, content string) (string, string, string, string)
}

type stdSanitize struct {
	htmlScape   func(string) string
	Sanitize    func(string) string
	HomeContent func(string) string
	TrimSpace   func(string) string
}

func newSanitizer() *stdSanitize {
	p := bluemonday.NewPolicy()
	p.AllowElements([]string{"p", "h1", "h2", "h3", "h4", "h5", "h6", "div", "span", "hr", "br", "b", "i", "strong", "em", "ol", "ul", "li", "pre", "code", "blockquote", "article", "section"}...)

	p2 := bluemonday.NewPolicy()
	p.AllowElements()
	return &stdSanitize{
		htmlScape:   html.UnescapeString,
		Sanitize:    p.Sanitize,
		HomeContent: p2.Sanitize,
		TrimSpace:   strings.TrimSpace,
	}
}

func (s *stdSanitize) Apply(title, desc, content string) (string, string, string, string) {
	title = s.apply(title)
	desc = s.apply(desc)
	content = s.apply(content)
	rawContent := s.applyHomeContent(content)

	return title, desc, content, rawContent
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

func (s *stdSanitize) applyHomeContent(str string) string {
	if str == "" {
		return ""
	}

	str = s.HomeContent(str)

	return str
}
