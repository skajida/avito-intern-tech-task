package download

import (
	"fmt"
	"net/http"
	mod "service-segs/internal/model"
	c "service-segs/internal/model/constants"
	hrf "service-segs/internal/model/http-response-forms"
	"strings"
)

type DownloadHandler struct {
	service downloader
}

func NewHandler(svc downloader) *DownloadHandler {
	return &DownloadHandler{service: svc}
}

func filenameParse(w http.ResponseWriter, r *http.Request, message string) (string, error) {
	args := strings.Split(r.URL.Path, "/")[2:]
	if len(args) > 1 {
		hrf.NewErrorResponse(r, message).Send(w, http.StatusUnprocessableEntity)
		return "", c.WrongRequest
	}
	return strings.Split(args[0], "?")[0], nil
}

func downloadHandle(service downloader, w http.ResponseWriter, r *http.Request) {
	filename, err := filenameParse(w, r, "Filename specified incorrectly")
	if err != nil {
		return
	}
	if file, err := service.DownloadFile(r.Context(), mod.Filename(filename)); err != nil {
		hrf.NewErrorResponse(r, "File not found").Send(w, http.StatusNotFound)
		return
	} else {
		w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		w.WriteHeader(http.StatusOK)
		w.Write(file)
	}
}

func (this *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		downloadHandle(this.service, w, r)
	default:
		w.Header().Set("Allow", fmt.Sprint(http.MethodPost, ", ", http.MethodDelete))
		hrf.NewErrorResponse(r, "API doesn't support the method").
			Send(w, http.StatusMethodNotAllowed)
	}
}
