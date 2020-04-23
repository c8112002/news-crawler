package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	qiitaBaseURL = "https://qiita.com/api/v2/"
	qiitaAction = "items"
	qiitaBaseParams = "?page=1&per_page=100"
	qiitaTimeFormat = "2006-01-02"
)

type ArticleCrawler struct {
	Token string
	Tags []string
	From time.Time
	To time.Time
}

func (a *ArticleCrawler) Run() ([]qiitaResult, error) {
	fromDay := a.From.Format(qiitaTimeFormat)
	toDay := a.To.Format(qiitaTimeFormat)

	var results []qiitaResult

	for _, tag := range a.Tags {

		query := fmt.Sprintf("&query=tag:%s+created:>=%s+created:<=%s", tag, fromDay, toDay)

		var items []qiitaItem

		endpoint, err := url.Parse(qiitaBaseURL + qiitaAction + qiitaBaseParams + query)
		if err != nil {
			return results, err
		}

		var header http.Header

		if len(a.Token) > 0 {
			header = http.Header{
				"Content-Type": {"application/json"},
				"Authorization": {"Bearer " + a.Token},
			}
		} else {
			header = http.Header{
				"Content-Type": {"application/json"},
			}
		}

		resp, err := http.DefaultClient.Do(&http.Request{
			Method:           http.MethodGet,
			URL:              endpoint,
			Header: header,
		})

		results, err = func() ([]qiitaResult, error) {
			defer func() {
				if err := resp.Body.Close(); err != nil {
					panic(err)
				}
			}()

			if err != nil {
				return results, err
			}

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return results, err
			}

			if err := json.Unmarshal(b, &items); err != nil {
				return results, err
			}

			var popularItems []qiitaItem
			for _, item := range extractPopularItems(items) {
				crawlThumbnail(tag, &item)
				popularItems = append(popularItems, item)
			}
			results = append(results, qiitaResult{Tag: tag, Items: popularItems})
			return results, nil
		}()
	}

	return results, nil
}

func extractPopularItems(source []qiitaItem) []qiitaItem {
	var articles []qiitaItem
	for _, a := range source {
		if a.Likes >= 10 {
			articles = append(articles, a)
		}
	}

	return articles
}

func crawlThumbnail(tag string, item *qiitaItem)  {
	res, err := http.Get(item.URL)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			panic(err.Error())
		}
	}()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	thumbnail, _ := doc.Find("meta[property='og:image']").First().Attr("content")
	fmt.Println(thumbnail)
	item.Thumbnail = thumbnail
}

type qiitaItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Likes int    `json:"likes_count"`
	Thumbnail string
}


func (a qiitaItem) String() string {
	return fmt.Sprintf("id: %v, title: %v, url: %v, likes: %v, thumbnail: %v", a.ID, a.Title, a.URL, a.Likes, a.Thumbnail)
}

type qiitaResult struct {
	Tag string
	Items []qiitaItem
}