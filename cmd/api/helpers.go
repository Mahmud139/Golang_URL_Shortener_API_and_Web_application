package main

import (
	"os"
	"strings"
)

func (app *application) checkCustomShort(url string) (bool, error) {
	exist, err := app.rdb.Exists(app.config.db.ctx, url).Result()
	if err != nil {
		return false, err
	}

	if exist == 1 {
		return true, nil
	}

	return false, nil
}

func DomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}

	newUrl := strings.Replace(url, "http://", "", 1)
	newUrl = strings.Replace(newUrl, "https://", "", 1)
	newUrl = strings.Replace(newUrl, "www.", "", 1)
	newUrl = strings.Split(newUrl, "/")[0]

	return newUrl != os.Getenv("DOMAIN")
}

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}

	return url
}
