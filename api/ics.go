package handler

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/thwiki/calendar-api-serverless/utils"
)

var (
	icsMaxAge  = os.Getenv("ICS_MAX_AGE")
	icsSMaxAge = os.Getenv("ICS_S_MAX_AGE")
)

func ICS(w http.ResponseWriter, r *http.Request) {
	isDownload := strings.HasSuffix(r.URL.Path, ".ics")

	calendar, err := utils.GetICS(time.Now())

	header := w.Header()
	header.Set("Cache-Control", "max-age="+icsMaxAge+", s-maxage="+icsSMaxAge+", public")

	if err != nil {
		header.Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write(utils.PrintError(err))
		return
	}

	if isDownload {
		header.Set("Content-Type", "text/calendar; charset=utf-8")
		header.Set("Content-Disposition", "attachment")
	} else {
		header.Set("Content-Type", "text/plain; charset=utf-8")
	}

	w.WriteHeader(http.StatusOK)
	calendar.SerializeTo(w)
}
