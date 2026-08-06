package templates

import (
	"io"
	"text/template"
)

func ParseFile(filename string, w io.Writer) error {
	tmpl, err := template.New("").
		Funcs(FuncMap).
		ParseFiles(filename)

	if err != nil {
		return err
	}

	for _, t := range tmpl.Templates() {
		if t.Name() == "" {
			continue
		}
		if err := tmpl.ExecuteTemplate(w, t.Name(), nil); err != nil {
			return err
		}

	}
	return nil
}
