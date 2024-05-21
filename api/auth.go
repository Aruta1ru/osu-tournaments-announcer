package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type AccessData struct {
	ClientID     int    `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Scope        string `json:"scope"`
}

type TokenData struct {
	AccessToken  string `json:"access_token"`
	TokenExpires int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

const apiURL = "https://osu.ppy.sh/api/v2/"

var CLIENT_ID int
var CLIENT_SECRET string

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	clientId, exists := os.LookupEnv("OSU_CLIENT_ID")
	if exists {
		CLIENT_ID, _ = strconv.Atoi(clientId)
	}
	clientSecret, exists := os.LookupEnv("OSU_CLIENT_SECRET")
	if exists {
		CLIENT_SECRET = clientSecret
	}
}

func GetAuthTokenFromFile(client *http.Client) (token string, err error) {

	tokenData := TokenData{}

	tokenFile, err := os.Open("tokenData.txt")
	if err != nil {
		return token, err
	}
	defer tokenFile.Close()

	var tokenStrings []string

	scanner := bufio.NewScanner(tokenFile)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "\n" {
			break
		}
		tokenStrings = strings.Split(line, "\t")
	}

	if len(tokenStrings) != 3 {
		return token, err
	}

	if tokenStrings[0] == "" || tokenStrings[1] == "" || tokenStrings[2] == "" {
		return token, err
	}

	expiredInt, err := strconv.Atoi(tokenStrings[2])
	if err != nil {
		return token, err
	}

	if int64(expiredInt) < time.Now().Unix() {
		return token, err
	}

	tokenData = TokenData{
		TokenType:    tokenStrings[0],
		AccessToken:  tokenStrings[1],
		TokenExpires: int64(expiredInt),
	}
	token = tokenData.TokenType + " " + tokenData.AccessToken
	return token, nil
}

func GetNewAuthToken(client *http.Client) (token string, err error) {

	tokenData := TokenData{}

	authURL := "https://osu.ppy.sh/oauth/token"

	accessData := AccessData{
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		GrantType:    "client_credentials",
		Scope:        "public",
	}

	requestBody := fmt.Sprintf(
		"client_id=%d&client_secret=%s&grant_type=%s&scope=%s",
		accessData.ClientID,
		accessData.ClientSecret,
		accessData.GrantType,
		accessData.Scope)

	req, _ := http.NewRequest("POST", authURL, bytes.NewReader([]byte(requestBody)))

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(body, &tokenData)
	if err != nil {
		return token, err
	}

	tokenFile, err := os.OpenFile("tokenData.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return token, err
	}
	defer tokenFile.Close()

	tokenData.TokenExpires = time.Now().Unix() + tokenData.TokenExpires

	fmt.Fprintf(tokenFile, "%s\t%s\t%d\n",
		tokenData.TokenType, tokenData.AccessToken, tokenData.TokenExpires)
	token = tokenData.TokenType + " " + tokenData.AccessToken
	return token, nil
}
