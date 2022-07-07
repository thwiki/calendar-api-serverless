package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/thwiki/calendar-api-serverless/utils"
)

var (
	responseMaxAge  = os.Getenv("RESPONSE_MAX_AGE")
	responseSMaxAge = os.Getenv("RESPONSE_S_MAX_AGE")
)

func Handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	startStr := query.Get("start")
	endStr := query.Get("end")

	pathParts := strings.Split(strings.TrimRight(r.URL.Path, "/"), "/")
	if len(pathParts) >= 4 {
		startStr = pathParts[len(pathParts)-2]
		endStr = pathParts[len(pathParts)-1]
	}

	header := w.Header()
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Content-Type", "application/json; charset=utf-8")

	start, err := utils.SanitizeDate(startStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(printError(err))
		return
	}
	end, err := utils.SanitizeDate(endStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(printError(err))
		return
	}

	events, err := utils.GetEvents(start, end)

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write(printError(err))
		return
	}

	header.Set("Cache-Control", "max-age="+responseMaxAge+", s-maxage="+responseSMaxAge+", public")
	w.WriteHeader(http.StatusOK)
	w.Write(events)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func printError(err error) []byte {
	response := ErrorResponse{
		Error: err.Error(),
	}
	errJson, _ := json.Marshal(response)
	return errJson
}
