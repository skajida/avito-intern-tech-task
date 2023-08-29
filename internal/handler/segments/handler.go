package segments

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"service-segs/internal/model/exchange/request"
	hrf "service-segs/internal/model/http-response-forms"
)

type SegmentsHandler struct {
	business segmentor
}

func NewHandler(service segmentor) *SegmentsHandler {
	return &SegmentsHandler{business: service}
}

func createHandle(service creator, w http.ResponseWriter, r *http.Request) {
	var requestSegment request.Segment
	body, _ := io.ReadAll(r.Body) // TODO whats the danger?
	if err := json.Unmarshal(body, &requestSegment); err != nil || requestSegment.SegId == "" {
		hrf.NewErrorResponse(r, "Segment specified incorrectly").
			Send(w, http.StatusUnprocessableEntity)
		return
	}
	if err := service.AddSegment(r.Context(), requestSegment); err != nil {
		hrf.NewErrorResponse(r, "Specified segment already exists").Send(w, http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func deleteHandle(service deletor, w http.ResponseWriter, r *http.Request) {
	seg_id := r.FormValue("seg_id") // TODO post/url?
	if seg_id == "" {
		hrf.NewErrorResponse(r, "Segment specified incorrectly").
			Send(w, http.StatusUnprocessableEntity)
		return
	}
	if err := service.RemoveSegment(r.Context(), request.Segment{SegId: seg_id}); err != nil {
		hrf.NewErrorResponse(r, "Specified segment doesn't exist").Send(w, http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (this *SegmentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createHandle(this.business, w, r)
	case http.MethodDelete:
		deleteHandle(this.business, w, r)
	default:
		w.Header().Set("Allow", fmt.Sprint(http.MethodPost, ", ", http.MethodDelete))
		hrf.NewErrorResponse(r, "API doesn't support the method").
			Send(w, http.StatusMethodNotAllowed)
	}
}
