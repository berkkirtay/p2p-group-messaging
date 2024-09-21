// Copyright (c) 2024 Berk Kirtay

package http

import (
	"encoding/json"
	"io"
	"main/infra/cryptography"
	"main/services/auth"
	"main/services/peer"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

/*
 * Service to make HTTP requests. Handles session cookies upon user authentication.
 */

const (
	OK                    = 200
	CREATED               = 201
	ACCEPTED              = 202
	NOT_FOUND             = 404
	INTERNAL_SERVER_ERROR = 500
)

var sessionHeader HeaderModel
var client *http.Client

func InitializeService(sessionAuth *auth.AuthenticationModel) {
	if sessionAuth != nil && sessionAuth.Token != "" {
		sessionHeader = CreateHeaderModel(
			WithContentType("application/json"),
			WithCookie(sessionAuth.Cookies),
			WithSession(sessionAuth.Id),
			WithAuthorization(sessionAuth.Token),
			WithPublicKey(sessionAuth.Cryptography.PublicKey))

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

func generateANewWithDiffieHellmanAuthKey(peer peer.Peer) {
	privateKey, publicKey := cryptography.GenerateEllipticCurveKeys()
	var token string = cryptography.DiffieHellman(privateKey, peer.Cryptography.Elliptic.PublicKey)
	sessionHeader = CreateHeaderModel(
		WithContentType("application/json"),
		WithCookie(sessionHeader.Cookie),
		WithSession(sessionHeader.Session),
		WithAuthorization(token),
		WithPublicKey(publicKey))
}

func GET(peer peer.Peer, path string, respType interface{}, params ...string) *http.Response {
	InitializeService(nil)
	generateANewWithDiffieHellmanAuthKey(peer)
	path = handleQueryParams(path, params)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		panic(err)
	}
	prepareHeadersForRequest(req)
	res, err := client.Do(req)
	if err != nil {
		return res
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if body == nil {
		return res
	}
	if len(body) != 0 {
		json.Unmarshal(body, &respType)
	}
	return res
}

func POST(peer peer.Peer, path string, payload string, respType interface{}, params ...string) *http.Response {
	InitializeService(nil)
	generateANewWithDiffieHellmanAuthKey(peer)
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
		return res
	}
	if len(body) != 0 {
		json.Unmarshal(body, &respType)
	}
	return res
}

func PUT() {

}

func DELETE(peer peer.Peer, path string, respType interface{}, params ...string) *http.Response {
	InitializeService(nil)
	generateANewWithDiffieHellmanAuthKey(peer)
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
	// urlObj, _ := url.Parse("/")
	// if len(sessionHeader.Cookie) > 0 {
	// 	client.Jar.SetCookies(urlObj, sessionHeader.Cookie)
	// }
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
