package handler

import (
	"fmt"
	"net/http"
	create_seg "service-segs/internal/handler/create-seg"
	delete_seg "service-segs/internal/handler/delete-seg"
)

type SegmentsHandler struct {
}

func (uh *SegmentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		create_seg.Handler(w, r)
	case http.MethodDelete:
		delete_seg.Handler(w, r)
	default:
		fmt.Fprintln(w, "unknown segments method")
	}
}
