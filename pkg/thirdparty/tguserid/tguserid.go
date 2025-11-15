package tguserid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/go-mojito/mojito"
	"github.com/nilathedragon/spamscale/config"
)

const TgUserIdCacheKey = "tguserid:tguserid:%s"

type TgUserData struct {
	ID         int64  `json:"id,string"`
	Type       string `json:"type"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Username   string `json:"username"`
	IsBot      bool   `json:"isBot"`
	IsVerified bool   `json:"isVerified"`
}

func GetUserID(username string) (int64, error) {
	user, err := GetUser(username)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func GetUser(username string) (*TgUserData, error) {
	cacheKey := fmt.Sprintf(TgUserIdCacheKey, username)
	if exists, err := mojito.DefaultCache().Contains(cacheKey); err != nil {
		return nil, err
	} else if exists {
		var user TgUserData
		if err := mojito.DefaultCache().Get(cacheKey, &user); err != nil {
			return nil, err
		}
		return &user, nil
	}

	user, err := fetchUser(username)
	if err != nil {
		return nil, err
	}
	mojito.DefaultCache().Set(cacheKey, *user)
	mojito.DefaultCache().ExpireAfter(cacheKey, config.GetCacheDuration())
	return user, nil
}

func fetchUser(username string) (*TgUserData, error) {
	csrfToken, err := fetchCsrfToken(username)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(
		"POST",
		"https://tg-user.id/api/get-userid",
		bytes.NewBufferString(fmt.Sprintf(`{"username": "%s"}`, username)),
	)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", "https://tg-user.id")
	request.Header.Set("Referer", "https://tg-user.id/from/username/"+username)
	request.Header.Set("X-CSRF-Token", csrfToken)

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user TgUserData
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return &user, nil
}

func fetchCsrfToken(username string) (string, error) {
	resp, err := http.Get("https://tg-user.id/from/username/" + username)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	regex, err := regexp.Compile(`<meta name="csrf-token" content="(.*)">`)
	if err != nil {
		return "", err
	}

	matches := regex.FindStringSubmatch(string(body))
	if len(matches) != 2 {
		return "", fmt.Errorf("failed to parse csrf token")
	}

	return matches[1], nil
}
