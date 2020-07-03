package handlers

import (
	"net/http"

	"git.jojii.de/jojii/zimserver/zim"
)

// Index handle index route
func Index(w http.ResponseWriter, r *http.Request, hd *HandlerData) error {
	var cards []HomeCards

	// Create cards
	for _, file := range hd.ZimService.GetFiles() {
		info := file.GetInfos()

		// Get Faviconlink
		var favurl string
		favIcon, err := file.Favicon()
		if err == nil {
			favurl = zim.GetRawWikiURL(file, favIcon)
		}

		cards = append(cards, HomeCards{
			Text:  info.GetDescription(),
			Title: info.Title,
			Image: favurl,
			Link:  zim.GetMainpageURL(file),
		})
	}

	return serveTemplate(HomeTemplate, HomeTemplateData{
		Cards: cards,
	}, w)
}
