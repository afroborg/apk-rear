package models

import "net/http"

type Handler struct {
	Path    string
	Handler http.HandlerFunc
	Methods []string
}
