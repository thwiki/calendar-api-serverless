package utils

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	client "github.com/bozd4g/go-http-client"
)

var (
	Api        = os.Getenv("API_SOURCE")
	httpClient = client.New(Api)
)

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

func ParseText(text string) (result string, err error) {
	request, err := httpClient.GetWith("/api.php", ParseParameter{
		Action:             "parse",
		Format:             "json",
		Text:               text,
		Prop:               "text",
		ContentModel:       "wikitext",
		Wrapoutputclass:    "",
		Disablelimitreport: 1,
		Disableeditsection: 1,
		Disabletoc:         1,
	})

	if err != nil {
		return
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return
	}

	var parseResult ParseResponse
	response.Get().To(&parseResult)

	if err != nil {
		return
	}

	resultHtml := parseResult.Parse.Text.Content
	reader := strings.NewReader(resultHtml)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return
	}

	result = strings.Join(doc.Find("body>p").Map(func(i int, s *goquery.Selection) string {
		if title, err := s.Html(); err == nil {
			return title
		}
		return ""
	}), "<br />")

	return
}

func GetEvents(start string, end string) (bytes []byte, err error) {
	request, err := httpClient.GetWith("/api.php", EventsParameter{
		Action:        "ask",
		Format:        "json",
		Formatversion: 2,
		Query:         "[[事件开始::>" + start + "]][[事件开始::<" + end + "]]|?事件类型=type|?事件颜色=color|?事件页面=name|?事件开始=start|?事件结束=end|?事件描述=desc|?事件图标=icon|sort=事件开始|order=asc",
	})

	if err != nil {
		return
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return
	}

	var originalResult SMWResponse
	response.Get().To(&originalResult)

	convertedResult, err := ConvertApiResult(&originalResult)
	if err != nil {
		return
	}

	return json.Marshal(convertedResult)
}

func ConvertApiResult(response *SMWResponse) (result ApiResult, err error) {
	result.Version = ""
	result.Meta = response.Query.Meta

	result.Results = make([]ApiResultEntry, len(response.Query.Results.Entries))
	for i, entry := range response.Query.Results.Entries {
		resultEntry := ApiResultEntry{
			Id: entry.Fulltext,
		}

		resultEntry.Start, err = strconv.Atoi(entry.Printouts.Start[0].Timestamp)
		if err != nil {
			return
		}
		resultEntry.StartStr, err = SanitizeSMWDate(entry.Printouts.Start[0].Raw)
		if err != nil {
			return
		}
		resultEntry.End, err = strconv.Atoi(entry.Printouts.End[0].Timestamp)
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
				html, err := ParseText(resultEntry.Desc)
				if err != nil {
					resultEntry.Desc = SanitizeWikiText(resultEntry.Desc)
				}
				resultEntry.Desc = html
			}
		}

		resultEntry.Url = entry.Fullurl

		resultEntry.Type = entry.Printouts.Type
		if len(entry.Printouts.Icon) > 0 {
			resultEntry.Icon = entry.Printouts.Icon[0].Fullurl
		}
		if len(entry.Printouts.Color) > 0 {
			resultEntry.Color = entry.Printouts.Color[0]
		}

		result.Results[i] = resultEntry
	}
	return
}
