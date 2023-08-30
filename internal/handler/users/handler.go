package users

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	hrf "service-segs/internal/model/http-response-forms"
	"strconv"
	"strings"
)

type UsersHandler struct {
	business useror
}

func NewHandler(service useror) *UsersHandler {
	return &UsersHandler{business: service}
}

func userIdParse(w http.ResponseWriter, r *http.Request, message string) (uint, error) {
	args := strings.Split(r.URL.Path, "/")[2:]
	user_id, err := strconv.Atoi(strings.Split(args[0], "?")[0])
	if len(args) > 1 || err != nil || user_id < 0 {
		hrf.NewErrorResponse(r, message).Send(w, http.StatusUnprocessableEntity)
		return 0, fmt.Errorf("Wrong user_id format")
	}
	return uint(user_id), nil
}

type activeSegments struct {
	ActiveUserSegments []string `json:"segs_active"`
}

func readHandle(service getter, w http.ResponseWriter, r *http.Request) {
	user_id, uierr := userIdParse(w, r, "User specified incorrectly")
	if uierr != nil {
		return
	}
	if segments, err := service.GetUserSegments(r.Context(), user_id); err != nil {
		hrf.NewErrorResponse(r, "User doesn't exist").Send(w, http.StatusNotFound)
	} else {
		json_repr, _ := json.Marshal(activeSegments{ActiveUserSegments: segments})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(json_repr)
	}
}

type requestSegments struct {
	WantedSegments   []string `json:"segs_to_add"`
	UnwantedSegments []string `json:"segs_to_remove"`
}

func updateHandle(service updater, w http.ResponseWriter, r *http.Request) {
	user_id, uierr := userIdParse(w, r, "Input parameters specified incorrectly")
	if uierr != nil {
		return
	}
	var requestSegments requestSegments
	body, _ := io.ReadAll(r.Body) // TODO whats the danger?
	if err := json.Unmarshal(body, &requestSegments); err != nil {
		hrf.NewErrorResponse(r, "Input parameters specified incorrectly").
			Send(w, http.StatusUnprocessableEntity)
		return
	}
	if err := service.UpdateUserSegments(r.Context(), user_id, requestSegments.WantedSegments, requestSegments.UnwantedSegments); err != nil {
		switch err {
		// case: User doesn't exist
		default: // TODO different errors
			hrf.NewErrorResponse(r, "Specified segment doesn't exist").
				Send(w, http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (this *UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		readHandle(this.business, w, r)
	case http.MethodPut:
		updateHandle(this.business, w, r)
	default:
		w.Header().Set("Allow", fmt.Sprint(http.MethodPost, ", ", http.MethodDelete))
		hrf.NewErrorResponse(r, "API doesn't support the method").
			Send(w, http.StatusMethodNotAllowed)
	}
}
