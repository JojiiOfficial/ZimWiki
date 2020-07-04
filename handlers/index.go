package handlers

import (
	"net/http"

	"git.jojii.de/jojii/zimserver/zim"
)

// Index handle index route
func Index(w http.ResponseWriter, r *http.Request, hd *HandlerData) error {
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

	return serveTemplate(HomeTemplate, HomeTemplateData{
		Cards: cards,
	}, w)
}
