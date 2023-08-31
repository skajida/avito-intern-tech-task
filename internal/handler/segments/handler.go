package segments

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	hrf "service-segs/internal/model/http-response-forms"
)

type SegmentsHandler struct {
	service segmentor
}

func NewHandler(svc segmentor) *SegmentsHandler {
	return &SegmentsHandler{service: svc}
}

type segment struct {
	SegId   string `json:"seg_id"`
	Percent int    `json:"percent,omitempty"`
}

func createHandle(service creator, w http.ResponseWriter, r *http.Request) {
	var requestSegment segment
	body, _ := io.ReadAll(r.Body) // TODO whats the danger
	if err := json.Unmarshal(body, &requestSegment); err != nil || requestSegment.SegId == "" {
		hrf.NewErrorResponse(r, "Input parameters specified incorrectly").
			Send(w, http.StatusUnprocessableEntity)
		return
	}
	if err := service.AddSegment(r.Context(), requestSegment.SegId); err != nil {
		hrf.NewErrorResponse(r, "Specified segment already exists").Send(w, http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func deleteHandle(service deletor, w http.ResponseWriter, r *http.Request) {
	segId := r.FormValue("seg_id")
	if segId == "" {
		hrf.NewErrorResponse(r, "Segment specified incorrectly").
			Send(w, http.StatusUnprocessableEntity)
		return
	}
	if err := service.DeleteSegment(r.Context(), segId); err != nil {
		hrf.NewErrorResponse(r, "Specified segment doesn't exist").Send(w, http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (this *SegmentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createHandle(this.service, w, r)
	case http.MethodDelete:
		deleteHandle(this.service, w, r)
	default:
		w.Header().Set("Allow", fmt.Sprint(http.MethodPost, ", ", http.MethodDelete))
		hrf.NewErrorResponse(r, "API doesn't support the method").
			Send(w, http.StatusMethodNotAllowed)
	}
}
