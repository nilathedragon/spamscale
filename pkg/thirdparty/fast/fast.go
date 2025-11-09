package fast

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"

	"github.com/go-mojito/mojito"
	"github.com/nilathedragon/spamscale/config"
)

const FastBlocklistCacheKey = "blocklist:fast"

func IsBlocked(userId int64) (bool, error) {
	blocklist, err := GetBlocklist()
	if err != nil {
		return false, err
	}
	return slices.Contains(blocklist, strconv.FormatInt(userId, 10)), nil
}

func GetBlocklist() ([]string, error) {
	if exists, err := mojito.DefaultCache().Contains(FastBlocklistCacheKey); err != nil {
		return nil, err
	} else if exists {
		var blocklist []string
		if err := mojito.DefaultCache().Get(FastBlocklistCacheKey, &blocklist); err != nil {
			return nil, err
		}
		return blocklist, nil
	}
	blocklist, err := fetchBlocklist()
	if err != nil {
		return nil, err
	}
	mojito.DefaultCache().Set(FastBlocklistCacheKey, blocklist)
	mojito.DefaultCache().ExpireAfter(FastBlocklistCacheKey, config.GetCacheDuration())
	return blocklist, nil
}

func fetchBlocklist() ([]string, error) {
	response, err := http.Get("https://countersign.chat/api/scammer_ids.json")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var blocklist []string
	if err := json.NewDecoder(response.Body).Decode(&blocklist); err != nil {
		return nil, err
	}

	return blocklist, nil
}
