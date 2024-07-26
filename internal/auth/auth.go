package auth

import (
	"errors"
	"net/http"
	"strings"
)

// find the api key in a HTTP request header and return it
// Auth header looks like: Authorization: apikey {key}
func GetAPIKey(headers http.Header) (string, error) {
	value := headers.Get("Authorization")

	if value == "" {
		return "", errors.New("no mapi key found")
	}

	values := strings.Split(value, " ")

	return values[1], nil

}
