package web

import (
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{
	"Cut": Trim,
}

func Trim(s template.HTML) template.HTML {
	return s[:180] + "..."
}

type TemplateData struct {
	FirstLocation  string
	SecondLocation string
	Articles       interface{}
}

func TemplateRender(w http.ResponseWriter, tmpl string, td *TemplateData) error {
	var t *template.Template
	if false {
		//	var ok = true
		//	t, ok = AppCfg.Cache[tmpl]
		//	if !ok {
		//		log.Fatal().Msg("cache is not working")
		//	}
	} else {
		cache, err := TemplateRenderCache()

		if err != nil {
			return err
		}

		t = cache[tmpl]
	}

	buf := new(bytes.Buffer)
	err := t.Execute(buf, td)

	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)

	if err != nil {
		return err
	}

	return nil
}

func TemplateRenderCache() (map[string]*template.Template, error) {
	pages, err := filepath.Glob("./platform/templates/*.page.tmpl")
	var templateCache = make(map[string]*template.Template)

	if err != nil {
		return templateCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return templateCache, err
		}

		matches, err := filepath.Glob("./platform/templates/*.layout.tmpl")

		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./platform/templates/*.layout.tmpl")
			if err != nil {
				return templateCache, err
			}
			templateCache[name] = ts
		}
	}

	return templateCache, nil
}
