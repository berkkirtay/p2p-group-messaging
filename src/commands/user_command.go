// Copyright (c) 2024 Berk Kirtay

package commands

import (
	"encoding/json"
	"fmt"
	"main/infra/cryptography"
	"main/infra/http"
	"main/services/auth"
	"main/services/user"
)

var sessionAuth auth.AuthenticationModel
var CurrentUser user.User

func HandleRegister(command []string) {
	if len(command) != 2 {
		fmt.Printf("Wrong usage.\n")
		return
	}

	name := command[1]
	userCrypto := cryptography.CreateCommonCrypto(
		name)
	dumpToFile(userCrypto.PrivateKey, "PRIVATE_KEY")
	dumpToFile(userCrypto.PublicKey, "PUBLIC_KEY")
	dumpToFile(userCrypto.Sign, "SIGN")
	userCrypto.PrivateKey = ""
	var user user.User = user.CreateUser(
		user.WithName(name),
		user.WithIsPeer(true),
		user.WithCryptography(userCrypto))

	body, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	res := http.POST(assignedPeer.Address+"/users", string(body), &user)
	if res.StatusCode == http.CREATED {
		fmt.Printf("Successfully registered as a new peer user.\n")
		fmt.Printf("Crypto files are exported.\n")
	}
}

func HandleLogin(command []string) {
	if len(command) > 3 || len(command) < 2 {
		fmt.Printf("Wrong usage.\n")
		return
	}
	userLogin := command[1]
	respBody := make([]user.User, 1)
	res := http.GET(assignedPeer.Address+"/users", &respBody, "id", userLogin, "name", userLogin)
	if res.StatusCode != http.OK {
		fmt.Printf("User %s does not exist.\n", userLogin)
		return
	}
	fmt.Printf("Hello %s! Your authentication process is ongoing..\n", respBody[0].Name)
	var authentication auth.AuthenticationModel = auth.CreateDefaultAuthenticationModel()
	if len(command) == 2 {
		publicKey := readFromFile("PUBLIC_KEY")
		privateKey := readFromFile("PRIVATE_KEY")
		sign := readFromFile("SIGN")
		userCrypto := cryptography.CreateCryptography(
			cryptography.WithSign(sign),
			cryptography.WithPublicKey(publicKey))
		authentication = auth.CreateAuthenticationModel(
			auth.WithId(userLogin),
			auth.WithName(userLogin),
			auth.WithCryptography(userCrypto))
		login(&authentication)
		if authentication.Token != "" {
			authentication.Token = cryptography.DecryptRSA(authentication.Token, privateKey)
		}
	} else if len(command) == 3 {
		password := command[2]
		authentication = auth.CreateAuthenticationModel(
			auth.WithId(userLogin),
			auth.WithName(userLogin),
			auth.WithPassword(password))
		login(&authentication)
	}
	if authentication.Token != "" {
		sessionAuth = authentication
		CurrentUser = respBody[0]
		http.InitializeService(sessionAuth.Cookies, sessionAuth.Id, sessionAuth.Token)
		fmt.Printf("You are authorized with the user-id: %s\n", authentication.Id)
	} else {
		fmt.Printf("Authentication process failed.. somehow.\n")
	}
}

func login(authenticationModel *auth.AuthenticationModel) {
	body, err := json.Marshal(authenticationModel)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	res := http.POST(assignedPeer.Address+"/auth", string(body), &authenticationModel)
	if res.StatusCode == http.ACCEPTED {
		authenticationModel.Cookies = res.Cookies()
	}
}
