package utils

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func ExtractAPIKeyDetails(
	req *http.Request,
) (string, string, error) {

	strategy := "basic_auth"
	var username, password string

	auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		strategy = "url_params"
	}

	switch strategy {
	case "basic_auth":
		payload, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			return "", "", fmt.Errorf("%v", "Unauthorized. Basic Authentication required.")
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			return "", "", fmt.Errorf("%v", "Unauthorized. Basic Authentication required, given auth not correct.")
		}

		username = pair[0]
		password = pair[1]
	case "url_params":

		username = req.URL.Query().Get("username")
		password = req.URL.Query().Get("password")

		if username == "" || password == "" {
			return "", "", fmt.Errorf("%v", "Unauthorized. Authentication required.")
		}

	default:
		return "", "", fmt.Errorf("%v", "Unauthorized. Authentication required.")
	}

	return username, password, nil
}
