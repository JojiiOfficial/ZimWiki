package handlers

import (
	"net/http"

	"github.com/JojiiOfficial/ZimWiki/zim"
)

var (
	version   string
	buildTime string
)

// Index handle index route
func Index(w http.ResponseWriter, r *http.Request, hd HandlerData) error {
	var cards []HomeCards

	// Create cards
	for i := range hd.ZimService.GetFiles() {
		file := &hd.ZimService.GetFiles()[i]

		info := file.GetInfos()

		// Get Faviconlink
		var favurl string
		favIcon, err := file.Favicon()
		if err == nil {
			favurl = zim.GetRawWikiURL(file, favIcon)
		}

		// Create homeCard
		cards = append(cards, HomeCards{
			Text:  info.GetDescription(),
			Title: info.Title,
			Image: favurl,
			Link:  zim.GetMainpageURL(file),
		})
	}

	// Cache files
	w.Header().Set("Cache-Control", "max-age=31536000, public")

	return serveTemplate(HomeTemplate, w, r, TemplateData{
		HomeTemplateData: HomeTemplateData{
			Cards:     cards,
			Version:   version,
			BuildTime: buildTime,
		},
	})
}
