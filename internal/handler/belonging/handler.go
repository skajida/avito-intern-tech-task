package belonging

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	hrf "service-segs/internal/model/http-response-forms"
	"strconv"
	"strings"
	"time"
)

type BelongingHandler struct {
	service belonger
}

func NewHandler(svc belonger) *BelongingHandler {
	return &BelongingHandler{service: svc}
}

func userIdParse(w http.ResponseWriter, r *http.Request, message string) (int, error) {
	args := strings.Split(r.URL.Path, "/")[2:]
	userId, err := strconv.Atoi(strings.Split(args[0], "?")[0])
	if len(args) > 1 || err != nil {
		hrf.NewErrorResponse(r, message).Send(w, http.StatusUnprocessableEntity)
		return 0, fmt.Errorf("Wrong user_id format")
	}
	return userId, nil
}

type activeSegments struct {
	ActiveUserSegments []string `json:"segs_active"`
}

func readHandle(service reader, w http.ResponseWriter, r *http.Request) {
	userId, uierr := userIdParse(w, r, "User specified incorrectly")
	if uierr != nil {
		return
	}
	if segments, err := service.ReadBelonging(r.Context(), userId); err != nil {
		hrf.NewErrorResponse(r, "User doesn't exist").Send(w, http.StatusNotFound)
	} else {
		jsonRepr, _ := json.Marshal(activeSegments{ActiveUserSegments: segments})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonRepr)
	}
}

type requestSegments struct {
	WantedSegments   []string `json:"segs_to_add"`
	UnwantedSegments []string `json:"segs_to_remove"`
}

func updateHandle(service modifier, w http.ResponseWriter, r *http.Request) {
	userId, uierr := userIdParse(w, r, "Input parameters specified incorrectly")
	if uierr != nil {
		return
	}
	var requestSegments requestSegments
	body, _ := io.ReadAll(r.Body) // TODO whats the danger
	if err := json.Unmarshal(body, &requestSegments); err != nil {
		hrf.NewErrorResponse(r, "Input parameters specified incorrectly").
			Send(w, http.StatusUnprocessableEntity)
		return
	}
	var err error
	if before, err := time.Parse(time.RFC3339, r.FormValue("before")); err != nil {
		err = service.ModifyBelongingTimer(
			r.Context(),
			userId,
			requestSegments.WantedSegments,
			requestSegments.UnwantedSegments,
			before,
		)
	} else {
		err = service.ModifyBelonging(
			r.Context(),
			userId,
			requestSegments.WantedSegments,
			requestSegments.UnwantedSegments,
		)
	}
	if err != nil {
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

func (this *BelongingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		readHandle(this.service, w, r)
	case http.MethodPut:
		updateHandle(this.service, w, r)
	default:
		w.Header().Set("Allow", fmt.Sprint(http.MethodPost, ", ", http.MethodDelete))
		hrf.NewErrorResponse(r, "API doesn't support the method").
			Send(w, http.StatusMethodNotAllowed)
	}
}
