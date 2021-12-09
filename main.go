package main

import (
	"encoding/json"
	"fmt"
	"os"
)

//goland:noinspection SpellCheckingInspection,GoUnusedConst
const (
	EnvSlackWebhook   = "SLACK_WEBHOOK"
	EnvSlackIcon      = "SLACK_ICON"
	EnvSlackIconEmoji = "SLACK_ICON_EMOJI"
	EnvSlackChannel   = "SLACK_CHANNEL"
	EnvSlackTitle     = "SLACK_TITLE"
	EnvSlackMessage   = "SLACK_MESSAGE"
	EnvSlackColor     = "SLACK_COLOR"
	EnvSlackUserName  = "SLACK_USERNAME"
	EnvGithubActor    = "GITHUB_ACTOR"
	EnvSiteName       = "SITE_NAME"
	EnvHostName       = "HOST_NAME"
	EnvDeployPath     = "DEPLOY_PATH"
	EnvMinimal        = "MSG_MINIMAL"
	EnvPSEUrl         = "PSE_URL"
	EnvPSEIP          = "PSE_IP"
	EnvPullRequestURL = "PULL_REQUEST_URL"
	EnvPSEVersion     = "PSE_VERSION"
	EnvUuid           = "UUID"
	EnvBiLink         = "BI_LINK"
	EnvBqLink         = "BQ_LINK"

	BlockSectionTypeHeader  = "header"
	BlockSectionTypeSection = "section"

	TextTypePlainText     = "plain_text"
	TextTypePlainMarkdown = "mrkdwn"

	BlockImageAccessory = "image"
)

type BlockText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type BlockAccessory struct {
	Type     string `json:"type"`
	ImageUrl string `json:"image_url"`
	AltText  string `json:"alt_text"`
}

type Block struct {
	Type      string          `json:"type"`
	Text      BlockText       `json:"text"`
	Accessory *BlockAccessory `json:"accessory,omitempty"`
}

type Webhook struct {
	Text        string       `json:"text,omitempty"`
	UserName    string       `json:"username,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Blocks      []Block      `json:"blocks,omitempty"`
}

type Attachment struct {
	Fallback   string  `json:"fallback"`
	Pretext    string  `json:"pretext,omitempty"`
	Color      string  `json:"color,omitempty"`
	AuthorName string  `json:"author_name,omitempty"`
	AuthorLink string  `json:"author_link,omitempty"`
	AuthorIcon string  `json:"author_icon,omitempty"`
	Footer     string  `json:"footer,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
}

type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

func main() {
	endpoint := os.Getenv(EnvSlackWebhook)
	if endpoint == "" {
		_, _ = fmt.Fprintln(os.Stderr, "URL is required")
		os.Exit(1)
	}

	actionUrl := fmt.Sprintf("%s/%s/actions/runs/%s/attempts/%s",
		os.Getenv("GITHUB_SERVER_URL"),
		os.Getenv("GITHUB_REPOSITORY"),
		os.Getenv("GITHUB_RUN_ID"),
		os.Getenv("GITHUB_RUN_ATTEMPT"))
	minimal := os.Getenv(EnvMinimal)
	var fields []Field
	if minimal == "true" {
		mainFields := []Field{
			{
				Title: os.Getenv(EnvSlackTitle),
				Value: os.Getenv(EnvSlackMessage),
				Short: false,
			},
		}
		fields = append(mainFields, fields...)
	} else {

		if os.Getenv(EnvPullRequestURL) != "" {
			fields = append([]Field{
				{
					Title: "Pull Request URL",
					Value: os.Getenv(EnvPullRequestURL),
					Short: false,
				},
			}, fields...)
		}
		if os.Getenv(EnvPSEUrl) != "" {
			fields = append([]Field{
				{
					Title: "PSE URL",
					Value: os.Getenv(EnvPSEUrl),
					Short: false,
				},
			}, fields...)
		}
		if os.Getenv(EnvPSEIP) != "" {
			fields = append([]Field{
				{
					Title: "PSE IP",
					Value: os.Getenv(EnvPSEIP),
					Short: false,
				},
			}, fields...)
		}
	}

	hostName := os.Getenv(EnvHostName)
	if hostName != "" {
		newFields := []Field{
			{
				Title: os.Getenv("SITE_TITLE"),
				Value: os.Getenv(EnvSiteName),
				Short: true,
			},
			{
				Title: os.Getenv("HOST_TITLE"),
				Value: os.Getenv(EnvHostName),
				Short: true,
			},
		}
		fields = append(newFields, fields...)
	}

	//goland:noinspection ALL
	githubActor := "http://github.com/" + os.Getenv(EnvGithubActor)
	var blocks = []Block{
		{
			Type: BlockSectionTypeHeader,
			Text: BlockText{
				Type: TextTypePlainText,
				Text: os.Getenv(EnvSlackTitle),
			},
		},
		{
			Type: BlockSectionTypeSection,
			Text: BlockText{
				Type: TextTypePlainMarkdown,
				Text: os.Getenv(EnvSlackMessage),
			},
		},
		{
			Type: BlockSectionTypeSection,
			Text: BlockText{
				Type: TextTypePlainMarkdown,
				Text: "*Actions URL:*\n" + actionUrl,
			},
		},
		{
			Type: BlockSectionTypeSection,
			Text: BlockText{
				Type: TextTypePlainMarkdown,
				Text: "*Run UUID:*\n" + os.Getenv(EnvUuid),
			},
		},
		{
			Type: BlockSectionTypeSection,
			Text: BlockText{
				Type: TextTypePlainMarkdown,
				Text: "*PSE Version:*\n" + os.Getenv(EnvPSEVersion),
			},
		},
		{
			Type: BlockSectionTypeSection,
			Text: BlockText{
				Type: TextTypePlainMarkdown,
				Text: "*BI (Metabase):*\n" + os.Getenv(EnvBiLink),
			},
		},
		{
			Type: BlockSectionTypeSection,
			Text: BlockText{
				Type: TextTypePlainMarkdown,
				Text: "*BigQuery:*\n" + os.Getenv(EnvBqLink),
			},
		},
		{
			Type: BlockSectionTypeSection,
			Text: BlockText{
				Type: TextTypePlainMarkdown,
				Text: githubActor,
			},
			//Accessory: BlockAccessory{
			//	Type: BlockImageAccessory,
			//	//ImageUrl: githubActor + ".png?size=32",
			//	ImageUrl: "https://api.slack.com/img/blocks/bkb_template_images/tripAgent_1.png",
			//	AltText:  envOr(EnvGithubActor, ""),
			//},
		},
	}

	//goland:noinspection HttpUrlsUsage
	msg := Webhook{
		UserName:  os.Getenv(EnvSlackUserName),
		IconURL:   os.Getenv(EnvSlackIcon),
		IconEmoji: os.Getenv(EnvSlackIconEmoji),
		Channel:   os.Getenv(EnvSlackChannel),
		Blocks:    blocks,
		Attachments: []Attachment{
			{
				Fallback:   envOr(EnvSlackMessage, "GITHUB_ACTION="+os.Getenv("GITHUB_ACTION")+" \n GITHUB_ACTOR="+os.Getenv("GITHUB_ACTOR")+" \n GITHUB_EVENT_NAME="+os.Getenv("GITHUB_EVENT_NAME")+" \n GITHUB_REF="+os.Getenv("GITHUB_REF")+" \n GITHUB_REPOSITORY="+os.Getenv("GITHUB_REPOSITORY")+" \n GITHUB_WORKFLOW="+os.Getenv("GITHUB_WORKFLOW")),
				Color:      envOr(EnvSlackColor, "good"),
				AuthorName: envOr(EnvGithubActor, ""),
				AuthorLink: "http://github.com/" + os.Getenv(EnvGithubActor),
				AuthorIcon: "http://github.com/" + os.Getenv(EnvGithubActor) + ".png?size=32",
				Footer:     "<https://github.com/rtCamp/github-actions-library|Powered By rtCamp's GitHub Actions Library>",
				Fields:     fields,
			},
		},
	}

	if err := send(endpoint, msg); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error sending message: %s\n", err)
		os.Exit(2)
	}
}

func envOr(name, def string) string {
	if d, ok := os.LookupEnv(name); ok {
		return d
	}
	return def
}

func send(endpoint string, msg Webhook) error {
	enc, _ := json.Marshal(msg)
	return fmt.Errorf("encoding %s", enc)
	//if err != nil {
	//	return err
	//}
	//b := bytes.NewBuffer(enc)
	//res, err := http.Post(endpoint, "application/json", b)
	//if err != nil {
	//	return err
	//}
	//
	//if res.StatusCode >= 299 {
	//	return fmt.Errorf("Error on message: %s\n", res.Status)
	//}
	//fmt.Println(res.Status)
	//return nil
}
