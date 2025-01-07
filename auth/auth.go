package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKeyFromHeader(header http.Header) (string, error) {
	key := header.Get("Authorization")
	if key == "" {
		return "", errors.New("Api key is not found!")
	}

	splits := strings.Split(key, " ")

	if len(splits) != 2 || len(splits[1]) != 64 {
		return "", errors.New("Invalid api key!")
	}
	if splits[0] != "ApiKey" {
		return "", errors.New("Invalid api key!")
	}
	return splits[1], nil
}
