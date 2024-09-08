package utils

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	client "github.com/bozd4g/go-http-client"
)

var (
	Api        = os.Getenv("API_SOURCE")
	ApiToken   = os.Getenv("API_TOKEN")
	Url        = os.Getenv("URL_SOURCE")
	ICSName    = os.Getenv("ICS_NAME")
	ICSTZ      = os.Getenv("ICS_TZ")
	httpClient = client.New(Api)
)

type UpstreamError struct{}

func (e *UpstreamError) Error() string {
	return "upstream error"
}

type EventsParameter struct {
	Action        string `url:"action"`
	Format        string `url:"format"`
	Formatversion int    `url:"formatversion"`
	Query         string `url:"query"`
}

type ParseParameter struct {
	Action             string `url:"action"`
	Format             string `url:"format"`
	Text               string `url:"text"`
	Prop               string `url:"prop"`
	ContentModel       string `url:"contentmodel"`
	Wrapoutputclass    string `url:"wrapoutputclass"`
	Disablelimitreport int    `url:"disablelimitreport"`
	Disableeditsection int    `url:"disableeditsection"`
	Disabletoc         int    `url:"disabletoc"`
}

func GetEvents(start string, end string) (events ApiResult, err error) {
	request, err := httpClient.GetWith("/api.php", EventsParameter{
		Action:        "ask",
		Format:        "json",
		Formatversion: 2,
		Query:         "[[事件开始::>" + start + "]][[事件开始::<" + end + "]]|?事件类型=type|?事件颜色=color|?事件页面=name|?事件开始=start|?事件结束=end|?事件描述=desc|?事件图标=icon|limit=200",
	})

	if err != nil {
		return
	}

	request.Header.Set("X-Api-Token", ApiToken)

	response, err := httpClient.Do(request)
	if err != nil {
		return
	}

	var originalResult SMWResponse
	response.Get().To(&originalResult)

	events, err = ConvertApiResult(&originalResult)
	return
}

func ConvertApiResult(response *SMWResponse) (result ApiResult, err error) {
	result.Version = ""
	result.Meta = response.Query.Meta

	if result.Meta.Hash == "" {
		err = &UpstreamError{}
		return
	}

	result.Results = make([]ApiResultEntry, len(response.Query.Results.Entries))
	for i, entry := range response.Query.Results.Entries {
		resultEntry := ApiResultEntry{
			Id: entry.Fulltext,
		}

		resultEntry.Start, err = strconv.ParseInt(entry.Printouts.Start[0].Timestamp, 10, 0)
		if err != nil {
			return
		}
		resultEntry.StartStr, err = SanitizeSMWDate(entry.Printouts.Start[0].Raw)
		if err != nil {
			return
		}
		resultEntry.End, err = strconv.ParseInt(entry.Printouts.End[0].Timestamp, 10, 0)
		if err != nil {
			return
		}
		resultEntry.EndStr, err = SanitizeSMWDate(entry.Printouts.End[0].Raw)
		if err != nil {
			return
		}

		resultEntry.Title = strings.TrimSpace(strings.Join(entry.Printouts.Name, " "))
		if resultEntry.Title == "" {
			if entry.Displaytitle != "" {
				resultEntry.Title = entry.Displaytitle
			} else {
				resultEntry.Title = entry.Fulltext[0 : len(entry.Fulltext)-strings.LastIndex(entry.Fulltext, "#")+1]
			}
		}

		resultEntry.Desc = strings.TrimSpace(strings.Join(entry.Printouts.Desc, " "))
		if resultEntry.Desc != "" {
			if strings.ContainsAny(resultEntry.Desc, "[{'}]") {
				resultEntry.Desc = SanitizeWikiText(resultEntry.Desc)
			}
		}

		resultEntry.Url = strings.Replace(entry.Fullurl, Api, Url, 1)

		resultEntry.Type = entry.Printouts.Type
		if len(entry.Printouts.Icon) > 0 {
			resultEntry.Icon = strings.Replace(entry.Printouts.Icon[0].Fullurl, Api, Url, 1)
		}
		if len(entry.Printouts.Color) > 0 {
			resultEntry.Color = entry.Printouts.Color[0]
		}

		result.Results[i] = resultEntry
	}

	sort.Slice(result.Results, func(i, j int) bool {
		a := result.Results[i]
		b := result.Results[j]
		if a.Start == b.Start {
			return strings.Compare(a.Title, b.Title) < 0
		}
		return a.Start < b.Start
	})

	return
}

func GetICS(now time.Time) (calendar *ics.Calendar, err error) {
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 6, 0)

	events, err := GetEvents(FormatDate(startDate), FormatDate(endDate))

	if err != nil {
		return
	}

	calendar = ics.NewCalendar()
	calendar.SetName(ICSName)
	calendar.SetTzid("Asia/Shanghai")
	calendar.SetProductId("-//THBWiki//Touhou Related Events Calendar")
	calendar.SetRefreshInterval("P30D")

	length := len(events.Results)
	for i := 0; i < length; i++ {
		resultEntry := events.Results[i]

		event := calendar.AddEvent(resultEntry.Id)
		event.SetSummary(resultEntry.Title)
		event.SetDescription(resultEntry.Desc)

		startDate := time.Unix(resultEntry.Start, 0)
		endDate := time.Unix(resultEntry.End, 0).AddDate(0, 0, 1)

		event.SetProperty(ics.ComponentPropertyDtStart, startDate.Format("20060102")+"T000000")
		event.SetProperty(ics.ComponentPropertyDtEnd, endDate.Format("20060102")+"T000000")
		event.SetProperty(ics.ComponentProperty(ics.PropertyTzid), ICSTZ)
		event.SetDtStampTime(now)
		event.SetURL(resultEntry.Url)
	}

	return
}
