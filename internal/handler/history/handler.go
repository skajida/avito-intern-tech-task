package history

import (
	"encoding/json"
	"fmt"
	"net/http"
	c "service-segs/internal/model/constants"
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
		return 0, c.WrongUser
	}
	return userId, nil
}

type historyUrl struct {
	Url string `json:"url"`
}

func requestHandle(service requester, w http.ResponseWriter, r *http.Request) {
	userId, err := userIdParse(w, r, "User specified incorrectly")
	if err != nil {
		return
	}
	year, month := r.FormValue("year"), r.FormValue("month")
	if year == "" || month == "" {
		hrf.NewErrorResponse(r, "Input parameters specified incorrectly").
			Send(w, http.StatusUnprocessableEntity)
		return
	}
	iYear, yerr := strconv.Atoi(year)
	iMonth, merr := strconv.Atoi(month)
	if yerr != nil || merr != nil || iMonth < 1 || iMonth > 12 {
		hrf.NewErrorResponse(r, "Input parameters specified incorrectly").
			Send(w, http.StatusUnprocessableEntity)
		return
	}
	url, err := service.RequestHistory(
		r.Context(),
		userId,
		time.Date(iYear, time.Month(iMonth), 1, 0, 0, 0, 0, time.Local),
	)
	if err != nil {
		hrf.NewErrorResponse(r, "User doesn't exist").
			Send(w, http.StatusNotFound)
		return
	}
	jsonRepr, _ := json.Marshal(historyUrl{Url: "/download/" + string(url)})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRepr)
}

func (this *HistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		requestHandle(this.service, w, r)
	default:
		w.Header().Set("Allow", fmt.Sprint(http.MethodPost, ", ", http.MethodDelete))
		hrf.NewErrorResponse(r, "API doesn't support the method").
			Send(w, http.StatusMethodNotAllowed)
	}
}
