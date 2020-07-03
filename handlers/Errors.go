package handlers

import "errors"

// ...
var (
	ErrNotFound          = errors.New("Not found")
	ErrNamespaceNotFound = errors.New("Namespace not found")
)
