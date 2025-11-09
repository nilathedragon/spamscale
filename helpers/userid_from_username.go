package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type TgUserData struct {
	ID         int64  `json:"id,string"`
	Type       string `json:"type"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Username   string `json:"username"`
	IsBot      bool   `json:"isBot"`
	IsVerified bool   `json:"isVerified"`
}

func GetUserIDFromUsername(username string) (int64, error) {
	resp, err := http.Get("https://tg-user.id/from/username/" + username)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)

	regex, err := regexp.Compile(`<meta name="csrf-token" content="(.*)">`)
	if err != nil {
		return 0, err
	}

	matches := regex.FindStringSubmatch(string(body))
	if len(matches) != 2 {
		fmt.Println(string(body))
		return 0, fmt.Errorf("failed to parse csrf token")
	}

	csrfToken := matches[1]

	request, err := http.NewRequest(
		"POST",
		"https://tg-user.id/api/get-userid",
		bytes.NewBufferString(fmt.Sprintf(`{"username": "%s"}`, username)),
	)

	if err != nil {
		return 0, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", "https://tg-user.id")
	request.Header.Set("Referer", "https://tg-user.id/from/username/"+username)
	request.Header.Set("X-CSRF-Token", csrfToken)

	resp, err = http.DefaultClient.Do(request)

	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err = io.ReadAll(resp.Body)

	var user TgUserData
	err = json.Unmarshal(body, &user)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
