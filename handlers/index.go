package handlers

import (
	"net/http"
)

// Index handle index route
func Index(w http.ResponseWriter, r *http.Request, hd *HandlerData) error {
	var cards []HomeCards

	// Create cards
	for _, file := range hd.ZimService.GetFiles() {
		info := file.GetInfos()
		cards = append(cards, HomeCards{
			Text:  info.GetDescription(),
			Title: info.Title,
			Link:  "lol",
		})
	}

	return serveTemplate(HomeTemplate, HomeTemplateData{
		Cards: cards,
	}, w)
}
