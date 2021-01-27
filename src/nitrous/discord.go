package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	api = "http://discord.com/api/v8/"
)

func check(code string) (bool, error) {
	path := fmt.Sprintf("entitlements/gift-codes/%v?with_application=false&with_subscription_plan=false", code)

	transport, err := getTransport()
	if err != nil {
		return false, err
	}

	client := http.Client{
		Timeout:   time.Second * time.Duration(timeout),
		Transport: transport,
	}

	request, err := http.NewRequest("GET", api+path, nil)
	if err != nil {
		return false, err
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/0.0.309 Chrome/83.0.4103.122 Electron/9.3.5 Safari/537.36")

	response, err := client.Do(request)
	if err != nil {
		return false, err
	}

	defer response.Body.Close()

	switch response.StatusCode {
	case 200:
		return true, nil
	case 404:
		return false, nil
	case 429:
		return false, errors.New("rate limited")
	}

	return false, errors.New("unkown error")
}
