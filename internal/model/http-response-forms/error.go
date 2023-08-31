package httpresponseforms

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ErrorResponse struct {
	Title        string `json:"title"`
	Request      string `json:"request"`
	Time         string `json:"time"`
	ErrorTraceId string `json:"errorTraceId"`
}

func NewErrorResponse(r *http.Request, title string) *ErrorResponse {
	return &ErrorResponse{
		Title:        title,
		Request:      fmt.Sprint(r.Method, " ", r.URL),
		Time:         time.Now().Format(time.RFC3339),
		ErrorTraceId: uuid.New().String(),
	}
}

func (this *ErrorResponse) Send(w http.ResponseWriter, status int) {
	defer log.Println("Responded with errorTraceId", this.ErrorTraceId)

	jsonRepr, _ := json.Marshal(this)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonRepr)
}
