package tpl

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"path"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/toughstruct/peaedge/log"
)

type CommonTemplate struct {
	Templates *template.Template
	AssetsFs  embed.FS
}

func NewCommonTemplate(fs embed.FS, dirs []string, funcMap map[string]interface{}) *CommonTemplate {
	var templates = template.New("GlobalTemplate").Funcs(funcMap)
	var ct = &CommonTemplate{Templates: templates, AssetsFs: fs}
	for _, d := range dirs {
		ct.parseDir(d)
	}
	return ct
}

func (ct *CommonTemplate) parseDir(dir string) {
	fss, _ := ct.AssetsFs.ReadDir(dir)
	for _, item := range fss {
		if item.IsDir() {
			continue
		}
		c, err := ct.AssetsFs.ReadFile(path.Join(dir, item.Name()))
		if err == nil {
			ct.parseItem(item, c, ct.Templates)
		}
	}
}

func (ct *CommonTemplate) ParseLocalDir(dir string, nameprifix string) error {
	fss, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, item := range fss {
		if item.IsDir() {
			continue
		}
		c, err := ioutil.ReadFile(path.Join(dir, item.Name()))
		if err != nil {
			log.Error(err)
			continue
		}
		name := nameprifix + strings.TrimSuffix(item.Name(), path.Ext(item.Name()))
		// if ct.Templates.Lookup(name) != nil {
		// 	continue
		// }
		tplstr := fmt.Sprintf(`{{define "%s"}}%s{{end}}`, name, c)
		t, err := ct.Templates.Parse(tplstr)
		if err != nil {
			log.Error(err)
			continue
		}
		ct.Templates = t
		return nil
	}
	return nil
}

func (ct *CommonTemplate) parseItem(item fs.DirEntry, c []byte, templates *template.Template) {
	name := strings.TrimSuffix(item.Name(), path.Ext(item.Name()))
	if templates.Lookup(name) != nil {
		return
	}
	tplstr := fmt.Sprintf(`{{define "%s"}}%s{{end}}`, name, c)
	ct.Templates = template.Must(templates.Parse(tplstr))
}

func (ct *CommonTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return ct.Templates.ExecuteTemplate(w, name, data)
}
