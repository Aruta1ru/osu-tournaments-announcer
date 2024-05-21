package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type ForumTopic struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	UserID    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type ForumPost struct {
	EditedAt string        `json:"edited_at"`
	Body     ForumPostBody `json:"body"`
}

type ForumPostBody struct {
	HTMLdata string `json:"html"`
}

type TopicResponse struct {
	Posts []ForumPost `json:"posts"`
	Topic ForumTopic  `json:"topic"`
}

type ForumPostDb struct {
	ID         int
	Title      string
	UserID     int
	CreatedAt  string
	EditedAt   string
	PicPreview string
	IsValid    bool
	IsNotified bool
	Links      []ForumPostLink
}

type ForumPostLink struct {
	ForumpostID int
	Name        string
	URL         string
}

var BLACKLISTED_LINKS = []string{
	"https://osu.ppy.sh/",
	"http://osu.ppy.sh/",
	"http://puu.sh/",
	"https://pif.ephemeral.ink",
	"#",
	"https://es.wikipedia.org",
	"https://www.geogebra.org",
	"http://www.timeanddate.com",
	"https://www.timeanddate.com",
	"https://pickem.hwc.hr/",
	"https://en.wikipedia.org",
	"mailto:",
	"https://soontm",
	".bandcamp.com",
	"twitter.com",
}

func getLinks(forumpostID int, body io.Reader) *[]ForumPostLink {
	z := html.NewTokenizer(body)
	links := make(map[string]string)
	var linksOrder []string
	for {
		var url, name string
		blacklisted := false
		tokenType := z.Next()
		switch tokenType {
		case html.StartTagToken:
			token := z.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						url = attr.Val
						for _, blacklistedLink := range BLACKLISTED_LINKS {
							if strings.HasPrefix(url, blacklistedLink) {
								blacklisted = true
								break
							}
						}
						nameFound := false
						for {
							tokenType = z.Next()
							switch tokenType {
							case html.TextToken:
								name = strings.TrimSpace(z.Token().Data)
								if len(name) > 0 {
									nameFound = true
								}
							case html.EndTagToken:
								token := z.Token()
								if token.Data == "a" {
									nameFound = true
								}
							}
							if nameFound {
								break
							}
						}

						if !blacklisted {
							foundDuplicate := false
							for _, link := range linksOrder {
								if url == link {
									foundDuplicate = true
									break
								}
							}
							if !foundDuplicate {
								links[url] = name
								linksOrder = append(linksOrder, url)
							}
						}
					}
				}
			}
		}
		if tokenType == html.ErrorToken {
			break
		}
	}
	var forumpostLinks []ForumPostLink
	for _, linkURL := range linksOrder {
		tempForumpostLink := ForumPostLink{
			ForumpostID: forumpostID,
			URL:         linkURL,
			Name:        links[linkURL],
		}
		forumpostLinks = append(forumpostLinks, tempForumpostLink)
	}
	return &forumpostLinks
}

func getPicPreview(body io.Reader) string {
	z := html.NewTokenizer(body)
	var pictureURL string
	for {
		tokenType := z.Next()
		switch tokenType {
		case html.StartTagToken, html.SelfClosingTagToken:
			token := z.Token()
			if token.Data == "img" {
				for _, attr := range token.Attr {
					if attr.Key == "src" {
						pictureURL = attr.Val
					}
				}
			}
		}
		if tokenType == html.ErrorToken || pictureURL != "" {
			break
		}
	}
	return pictureURL
}

func GetForumpostData(client *http.Client, forumpostID int, token string) (*ForumPostDb, error) {
	endpoint := fmt.Sprintf("%sforums/topics/%d", apiURL, forumpostID)

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	q := req.URL.Query()
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 {
		err = fmt.Errorf("%s", "You are not authorized!")
		return nil, err
	}

	if resp.StatusCode == 404 {
		err = fmt.Errorf("%s: %d", "This post does not exist", forumpostID)
		return nil, err
	}

	topicResponse := TopicResponse{}
	err = json.Unmarshal(body, &topicResponse)
	if err != nil {
		return nil, err
	}

	topicLinks := getLinks(topicResponse.Topic.ID,
		strings.NewReader(topicResponse.Posts[0].Body.HTMLdata))
	picPreview := getPicPreview(strings.NewReader(topicResponse.Posts[0].Body.HTMLdata))

	forumpostToSave := ForumPostDb{
		ID:         topicResponse.Topic.ID,
		Title:      topicResponse.Topic.Title,
		UserID:     topicResponse.Topic.UserID,
		CreatedAt:  topicResponse.Topic.CreatedAt,
		EditedAt:   topicResponse.Posts[0].EditedAt,
		PicPreview: picPreview,
		IsValid:    true,
		IsNotified: false,
	}

	forumpostToSave.Links = *topicLinks

	if len(*topicLinks) == 0 {
		forumpostToSave.IsValid = false
	}

	return &forumpostToSave, nil
}
