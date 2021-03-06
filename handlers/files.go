package handlers

// ...
const (
	// Global paths
	root = "html/static/"

	// Paths
	AssetsPath    = root + "assets/"
	TemplatesPath = root + "templates/"

	// Templates
	BaseTemplate     = root + "index.html"
	HomeTemplate     = TemplatesPath + "home.html.tpl"
	WikiPageTemplate = TemplatesPath + "viewPage.html.tpl"
	SearchTemplate   = TemplatesPath + "search.html.tpl"
)
