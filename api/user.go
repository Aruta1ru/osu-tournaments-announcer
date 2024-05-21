package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	AvatarURL   string `json:"avatar_url"`
	CountryCode string `json:"country_code"`
}

func GetUserData(client *http.Client, userId int, token string) (*User, error) {
	endpoint := fmt.Sprintf("%susers/%d", apiURL, userId)

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	user := User{}

	if resp.StatusCode == 401 {
		return nil, err
	}

	if resp.StatusCode == 404 {
		user = User{
			ID:          userId,
			Username:    "Not found or restricted",
			AvatarURL:   "",
			CountryCode: "UG",
		}
		return &user, nil
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
