package handlers

import "git.jojii.de/jojii/zimserver/zim"

// HandlerData data for handler funcs
type HandlerData struct {
	ZimService *zim.Handler
	AcceptGzip bool
}
