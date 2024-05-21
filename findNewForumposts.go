package main

import (
	"discord-go/api"
	"discord-go/db"
	"discord-go/scraper"
	"fmt"
	"net/http"
	"time"
)

func main() {

	latestID, err := db.GetLatestForumpostID()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	newForumpostIDs, err := scraper.ScrapeNewForumposts(latestID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(newForumpostIDs) == 0 {
		fmt.Println("No new forumposts found at the moment")
		return
	}

	client := &http.Client{}

	token, err := api.GetAuthTokenFromFile(client)
	if err != nil {
		fmt.Println("Oops! there was an error:", err)
		return
	}

	if token == "" {
		token, err = api.GetNewAuthToken(client)
		if err != nil {
			fmt.Println("Oops! there was an error:", err)
			return
		}
	}

	for _, id := range newForumpostIDs {
		forumpost, err := api.GetForumpostData(client, id, token)
		if err != nil {
			fmt.Printf("Error get forumpost data (ID = %d): %s\n", id, err)
			return
		}

		user, err := api.GetUserData(client, forumpost.UserID, token)
		if err != nil {
			fmt.Printf("Error get user data (ID = %d): %s\n", forumpost.UserID, err)
			return
		}

		time.Sleep(time.Millisecond * 5)

		if !db.InsertUser(user) {
			fmt.Println("Cannot insert user data ID =", forumpost.UserID)
			return
		}

		if !db.InsertForumpost(forumpost) {
			fmt.Println("Cannot insert forumpost data ID =", id)
			return
		}
	}

}
