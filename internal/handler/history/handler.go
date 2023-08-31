package history

import (
	"encoding/json"
	"fmt"
	"net/http"
	hrf "service-segs/internal/model/http-response-forms"
	"strconv"
	"strings"
	"time"
)

type HistoryHandler struct {
	service requester
}

func NewHandler(svc requester) *HistoryHandler {
	return &HistoryHandler{service: svc}
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

type historyUrl struct {
	Url string `json:"url"`
}

func requestHandler(service requester, w http.ResponseWriter, r *http.Request) {
	userId, uierr := userIdParse(w, r, "User specified incorrectly")
	if uierr != nil {
		return
	}
	year, month := r.FormValue("year"), r.FormValue("month")
	if year == "" || month == "" {
		hrf.NewErrorResponse(r, "Input parameters specified incorrectly").
			Send(w, http.StatusUnprocessableEntity)
		return
	}
	iYear, iyerr := strconv.Atoi(year)
	iMonth, imerr := strconv.Atoi(month)
	if iyerr != nil || imerr != nil || iMonth < 1 || iMonth > 12 {
		hrf.NewErrorResponse(r, "Input parameters specified incorrectly").
			Send(w, http.StatusUnprocessableEntity)
		return
	}
	url, rerr := service.RequestHistory(r.Context(), userId, time.Date(iYear, time.Month(iMonth), 1, 0, 0, 0, 0, time.Local))
	if rerr != nil {
		hrf.NewErrorResponse(r, "User doesn't exist").
			Send(w, http.StatusNotFound)
		return
	}
	jsonRepr, _ := json.Marshal(historyUrl{Url: url})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRepr)
}

func (this *HistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		requestHandler(this.service, w, r)
	default:
		w.Header().Set("Allow", fmt.Sprint(http.MethodPost, ", ", http.MethodDelete))
		hrf.NewErrorResponse(r, "API doesn't support the method").
			Send(w, http.StatusMethodNotAllowed)
	}
}
