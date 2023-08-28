package handler

import (
	"fmt"
	"net/http"
	get_user_segs "service-segs/internal/handler/get-user-segs"
	update_user_segs "service-segs/internal/handler/update-user-segs"
)

type UserHandler struct {
}

func (uh *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		update_user_segs.Handler(w, r)
	case http.MethodGet:
		get_user_segs.Handler(w, r)
	default:
		fmt.Fprintln(w, "unknown user method")
	}
}
