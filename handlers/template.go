package handlers

import (
	"embed"
	"html/template"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/chai2010/gettext-go"
)

var (
	// TemplateCache is a cache of templates
	TemplateCache = make(map[string]*template.Template)
	WebFS         embed.FS
	LocaleByte    []byte
)

// 						      //
// --- Template data structs  //
// 						      //

// TemplateData data for the base template
type TemplateData struct {
	Favtype   string
	FavIcon   string
	Wiki      string
	Namespace string

	HomeTemplateData
	WikiViewTemplateData
	SearchTemplateData
}

// HomeCards contain data for the
// cards displayed on the home site
type HomeCards struct {
	Image string
	Title string
	Text  string
	Link  string
}

// HomeTemplateData contain data
// for the home template
type HomeTemplateData struct {
	Cards     []HomeCards
	Version   string
	BuildTime string
}

// WikiViewTemplateData data for wiki view page
type WikiViewTemplateData struct {
	Source string
}

// SearchTemplateData data for search
type SearchTemplateData struct {
	SearchSource string
	QueryText    string
	Results      []SearchResult
	NbResults    int
	TimeTook     time.Duration
	ActualPageNb int
	NbPages      int
	PreviousPage int
	NextPage     int
}

// SearchResult represents a single result of a search query
type SearchResult struct {
	Img   string
	Wiki  string
	Title string
	Link  string
}

func translate(input string) string {
	return gettext.PGettext("", input)
}

// 						    //
// --- Template functions   //
// 						    //

// Load and execute
func serveTemplate(tmpFile string, w http.ResponseWriter, r *http.Request, btd TemplateData) error {
	var err error
	tmplName := path.Base(tmpFile)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Find in cache
	tmpl, has := TemplateCache[tmplName]

	// Get the Accept-Language header from the HTTP request
	headerLang := r.Header.Get("Accept-Language")

	// Keep only the first element of the array
	lang := strings.Split(headerLang, ",")[0]

	// e.g. en-GB -> en
	if strings.Contains(lang, "-") {
		lang = strings.Split(lang, "-")[0]
	}

	// Get locale
	gettext.BindLocale(gettext.New("ZimWiki", "locale", LocaleByte))

	funcMap := template.FuncMap{
		"gettext": translate,
	}

	if !has {
		// Parse if not in cache
		tmpl, err = template.New(tmplName).Funcs(funcMap).ParseFS(WebFS, BaseTemplate, tmpFile)
		if err != nil {
			return err
		}

		// Cache template
		TemplateCache[tmplName] = tmpl
	}

	if len(btd.Wiki) == 0 {
		btd.Wiki = "-"
	}

	gettext.SetLanguage(lang)

	// Execute template
	return tmpl.ExecuteTemplate(w, path.Base(BaseTemplate), btd)
}
