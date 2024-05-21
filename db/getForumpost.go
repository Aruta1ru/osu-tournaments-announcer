package db

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func processLink(url string) string {
	if strings.HasPrefix(url, "https://forms.gle") || strings.HasPrefix(url, "https://docs.google.com/forms") {
		return "Google Form"
	}
	if strings.HasPrefix(url, "https://discord.gg") {
		return "Discord Server"
	}
	if strings.HasPrefix(url, "https://challonge.com") {
		return "Challonge Bracket"
	}
	if strings.HasPrefix(url, "https://www.twitch.tv") {
		return "Twitch Channel"
	}
	if strings.HasPrefix(url, "https://docs.google.com/spreadsheets") {
		return "Google Spreadsheet"
	}
	if strings.HasPrefix(url, "https://docs.google.com/document") {
		return "Google Document"
	}
	return fmt.Sprintf("Unrecognized: %s", strings.ReplaceAll(url, "https://", ""))
}

func GetNonNotifiedForumpostsID() ([]int, error) {
	var ids []int
	db := ConnectDB()
	if db == nil {
		err := fmt.Errorf("%s", "Cannot connect to database")
		db.Close()
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT id 
		  FROM forumposts 
		 WHERE is_valid = true 
		   AND is_notified = false 
		 ORDER BY id`)
	if err != nil {
		return nil, err
	}

	var id int

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func GetLatestForumpostID() (int, error) {
	db := ConnectDB()
	if db == nil {
		err := fmt.Errorf("%s", "Cannot connect to database")
		db.Close()
		return -1, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT MAX(id) AS latest_id 
		  FROM forumposts`)
	if err != nil {
		return -1, err
	}

	var latestID int

	for rows.Next() {
		err = rows.Scan(&latestID)
		if err != nil {
			return -1, err
		}
	}

	return latestID, nil
}

func GetForumpostDataById(id int) (*discordgo.MessageEmbed, error) {
	db := ConnectDB()
	if db == nil {
		err := fmt.Errorf("%s", "Cannot connect to database")
		db.Close()
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf(`
		SELECT forumposts.id, title, user_id, created_at, picture_preview, username, avatar, is_valid 
		  FROM forumposts 
		 	   INNER JOIN users ON forumposts.user_id = users.id 
		 WHERE forumposts.id=%d`, id))
	if err != nil {
		return nil, err
	}

	var (
		forumpostID int
		title       string
		userID      int
		created     string
		pictureURL  string
		username    string
		avatar      string
		valid       bool
	)

	exists := false

	for rows.Next() {
		exists = true
		err = rows.Scan(&forumpostID, &title, &userID, &created, &pictureURL, &username, &avatar, &valid)
		if err != nil {
			return nil, err
		}
	}

	if !exists {
		return nil, nil
	}

	embedForumpost := discordgo.MessageEmbed{}

	embedFooter := discordgo.MessageEmbedFooter{
		Text:    "Posted at ",
		IconURL: avatar,
	}

	embedAuthor := discordgo.MessageEmbedAuthor{
		URL:     fmt.Sprintf("%s%d", USER_URL, userID),
		Name:    "Host: " + username,
		IconURL: avatar,
	}

	embedForumpost.Title = ":trophy: " + title
	embedForumpost.Color = 0x00ff00
	embedForumpost.URL = fmt.Sprintf("%s%d", FORUMPOST_URL, forumpostID)
	embedForumpost.Timestamp = created
	embedForumpost.Footer = &embedFooter
	embedForumpost.Image = &discordgo.MessageEmbedImage{URL: pictureURL}
	embedForumpost.Author = &embedAuthor

	if valid {
		embedForumpost.Description = "TOURNEY LINKS:\n"

		rows, err = db.Query(fmt.Sprintf(`
			SELECT name, url 
			  FROM forumpost_links 
			 WHERE forumpost_id=%d`, forumpostID))
		if err != nil {
			return nil, err
		}

		var (
			linkName string
			linkURL  string
		)

		for rows.Next() {
			err = rows.Scan(&linkName, &linkURL)
			if err != nil {
				return nil, err
			}
			if linkName == "" || strings.HasPrefix(linkName, "http") {
				linkName = processLink(linkURL)
			}
			embedForumpost.Description += fmt.Sprintf(":link: [%s](%s)\n", linkName, linkURL)
		}
	}

	return &embedForumpost, nil
}
