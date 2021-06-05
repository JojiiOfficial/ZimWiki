package handlers

import "github.com/JojiiOfficial/ZimWiki/zim"

// HandlerData data for handler funcs
type HandlerData struct {
	ZimService *zim.Handler
	AcceptGzip bool
}
