package db

import (
	"database/sql"
	"discord-go/api"
	"fmt"
)

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func checkExistingRecordById(tablename string, id int) (bool, error) {
	db := ConnectDB()
	if db == nil {
		err := fmt.Errorf("%s", "Cannot connect to database")
		db.Close()
		return false, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf(`
		SELECT EXISTS(SELECT 1 FROM %s WHERE id=%d)`, tablename, id))
	if err != nil {
		return false, err
	}

	var exists bool

	for rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return false, err
		}
	}

	return exists, nil
}

func InsertUser(userToAdd *api.User) bool {
	db := ConnectDB()
	if db == nil {
		fmt.Println("Cannot connect to database!")
		db.Close()
		return false
	}
	defer db.Close()

	userExists, err := checkExistingRecordById("users", userToAdd.ID)
	if err != nil {
		fmt.Println("Error find user:", err)
		return false
	}
	if userExists {
		return true
	}

	stmt, err := db.Prepare(`
		INSERT INTO users (id, username, avatar, country_code) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		fmt.Println("Error insert prepared:", err)
		return false
	}
	_, err = stmt.Exec(userToAdd.ID, userToAdd.Username, userToAdd.AvatarURL, userToAdd.CountryCode)
	if err != nil {
		fmt.Println("Error insert exec:", err)
		return false
	}

	return true
}

func SetForumpostsNotified(ids []int) error {
	db := ConnectDB()
	if db == nil {
		err := fmt.Errorf("%s", "Cannot connect to database")
		db.Close()
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE forumposts SET is_notified = true WHERE id = $1")
	if err != nil {
		return err
	}
	for _, id := range ids {
		_, err = stmt.Exec(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertForumpost(forumpostToAdd *api.ForumPostDb) bool {
	db := ConnectDB()
	if db == nil {
		fmt.Println("Cannot connect to database!")
		db.Close()
		return false
	}
	defer db.Close()

	forumpostExists, err := checkExistingRecordById("forumposts", forumpostToAdd.ID)
	if err != nil {
		fmt.Println("Error find user:", err)
		return false
	}
	if !forumpostExists {
		stmt, err := db.Prepare(`INSERT INTO forumposts 
		(id, title, user_id, created_at, edited_at, picture_preview, is_valid) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`)
		if err != nil {
			fmt.Println("Error:", err)
			return false
		}

		_, err = stmt.Exec(
			forumpostToAdd.ID,
			forumpostToAdd.Title,
			forumpostToAdd.UserID,
			forumpostToAdd.CreatedAt,
			NewNullString(forumpostToAdd.EditedAt),
			forumpostToAdd.PicPreview,
			forumpostToAdd.IsValid)
		if err != nil {
			fmt.Println("Error:", err)
			return false
		}
	}

	stmt, err := db.Prepare(`
		INSERT INTO forumpost_links (forumpost_id, name, url) VALUES ($1, $2, $3)`)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	for _, link := range forumpostToAdd.Links {
		_, err = stmt.Exec(link.ForumpostID, link.Name, link.URL)
		if err != nil {
			fmt.Println("Error:", err)
			return false
		}
	}

	return true
}
