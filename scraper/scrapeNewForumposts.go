package scraper

import (
	"log"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

func ScrapeNewForumposts(latestID int) ([]int, error) {
	forumpostIDs := make([]int, 0)

	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	c.OnHTML("li.forum-topic-entry", func(e *colly.HTMLElement) {
		id, err := strconv.Atoi(e.Attr("data-topic-id"))
		if err != nil {
			log.Println("Something went wrong: ", err)
		}
		if id > latestID {
			forumpostIDs = append(forumpostIDs, id)
		}
	})

	url := "https://osu.ppy.sh/community/forums/55?sort=created&page=1"
	err := c.Visit(url)
	if err != nil {
		return forumpostIDs, err
	}

	log.Println("Found:", len(forumpostIDs), "forumposts")

	return forumpostIDs, nil
}
