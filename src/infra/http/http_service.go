// Copyright (c) 2024 Berk Kirtay

package http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

/*
 * Service to make HTTP requests. Handles session cookies upon user authentication.
 */

var sessionHeader HeaderModel
var client *http.Client

func InitializeService(cookies []*http.Cookie, headerData ...string) {
	if len(headerData) > 1 {
		sessionHeader = CreateHeaderModel(
			WithContentType("application/json"),
			WithCookie(cookies),
			WithSession(headerData[0]),
			WithAuthorization(headerData[1]))
	}
	if client == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			panic(err)
		}
		client = &http.Client{
			Jar: jar,
		}
	}

}

func GET(path string, respType interface{}, params ...string) *http.Response {
	if client == nil {
		InitializeService(nil)
	}
	path = handleQueryParams(path, params)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		panic(err)
	}
	prepareHeadersForRequest(req)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if body == nil {
		panic(err)
	}
	if len(body) != 0 {
		err = json.Unmarshal(body, &respType)
		if err != nil {
			panic(err)
		}
	}
	return res
}

func POST(path string, payload string, respType interface{}, params ...string) *http.Response {
	if client == nil {
		InitializeService(nil)
	}
	path = handleQueryParams(path, params)
	req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(payload))
	if err != nil {
		panic(err)
	}
	prepareHeadersForRequest(req)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if body == nil {
		panic(err)
	}
	if len(body) != 0 {
		err = json.Unmarshal(body, &respType)
		if err != nil {
			panic(err)
		}
	}
	return res
}

func PUT() {

}

func DELETE(path string, respType interface{}, params ...string) *http.Response {
	if client == nil {
		InitializeService(nil)
	}
	path = handleQueryParams(path, params)
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		panic(err)
	}
	prepareHeadersForRequest(req)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	return res
}

func prepareHeadersForRequest(req *http.Request) {
	var headerMap map[string]string
	str, _ := json.Marshal(sessionHeader)
	json.Unmarshal(str, &headerMap)
	for key, value := range headerMap {
		req.Header.Add(key, value)
	}
	urlObj, _ := url.Parse("/")
	if len(sessionHeader.Cookie) > 0 {
		client.Jar.SetCookies(urlObj, sessionHeader.Cookie)
	}
}

func handleQueryParams(url string, params []string) string {
	for i, param := range params {
		if i == 0 {
			url += "?"
		} else if i%2 == 0 {
			url += "&"
		} else {
			url += "="
		}
		url += param
	}
	return url
}
