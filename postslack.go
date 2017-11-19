package clitoolgoaws

import (
    "encoding/json"
    "io/ioutil"
	"net/http"
    "net/url"
	"fmt"
)

const (
	USERNAME = "GoBot"
	CHANNEL  = "iot_platform_alert"
	URL      = "https://hooks.slack.com/services/XXXXXX"
)

var (
	WebhookUrl string = URL
)

type Slack struct {
	Data        string `json:"text"`
	Username    string `json:"username"`
	Icon_emoji  string `json:"icon_emoji"`
	Icon_url    string `json:"icon_url"`
	Channel     string `json:"channel"`
}

func PostSlack(billing float64) {
	fmt.Println(billing)
	params, _ := json.Marshal(Slack{
		"AWS(bct-Prd)現在の料金: $"+fmt.Sprint(billing) ,
		USERNAME,
		"",
		"http://www.techscore.com/blog/wp/wp-content/uploads/2016/12/gopher_ueda.png",
		CHANNEL})

	res, _ := http.PostForm(
		WebhookUrl,
		url.Values{"payload": {string(params)}},
	)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	println(string(body))
}