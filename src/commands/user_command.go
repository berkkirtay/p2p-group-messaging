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
	userCrypto := cryptography.CreateCommonCrypto(name)
	dumpToFile(userCrypto.PublicKey, fmt.Sprintf("./keys/PUBLIC_KEY_%s", name))
	dumpToFile(userCrypto.PrivateKey, fmt.Sprintf("./keys/PRIVATE_KEY_%s", name))
	dumpToFile(userCrypto.Sign, fmt.Sprintf("./keys/SIGN_%s", name))
	userCrypto.PrivateKey = ""
	userCrypto.Elliptic.PrivateKey = nil
	var user user.User = user.CreateUser(
		user.WithName(name),
		user.WithIsPeer(true),
		user.WithCryptography(userCrypto))

	body, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	res := http.POST(assignedPeer, assignedPeer.Address+"/users", string(body), &user)
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
	res := http.GET(assignedPeer, assignedPeer.Address+"/users", &respBody, "id", userLogin, "name", userLogin)
	if res.StatusCode != http.OK {
		fmt.Printf("User %s does not exist.\n", userLogin)
		return
	}
	fmt.Printf("Hello %s! Your authentication process is ongoing..\n", respBody[0].Name)
	var authentication auth.AuthenticationModel = auth.CreateDefaultAuthenticationModel()
	if len(command) == 2 {
		authentication = loginWithECDH(userLogin, respBody[0].Id)
	} else if len(command) == 3 {
		password := command[2]
		authentication = loginWithPassword(userLogin, password)
	}
	if authentication.Token != "" {
		sessionAuth = authentication
		CurrentUser = respBody[0]
		http.InitializeService(&sessionAuth)
		fmt.Printf("You are authorized %s!\n", authentication.Name)
	} else {
		fmt.Printf("Authentication process failed.. somehow.\n")
	}
}

func loginWithECDH(userLogin string, userId string) auth.AuthenticationModel {
	publicKey := readFromFile(fmt.Sprintf("./keys/PUBLIC_KEY_%s", userLogin))
	privateKey := readFromFile(fmt.Sprintf("./keys/PRIVATE_KEY_%s", userLogin))
	sign := readFromFile(fmt.Sprintf("./keys/SIGN_%s", userLogin))

	ellipticPrivate, ellipticPublic := cryptography.GenerateEllipticCurveKeys()
	userCrypto := cryptography.CreateCryptography(
		cryptography.WithSign(sign),
		cryptography.WithPublicKey(publicKey),
		cryptography.WithPrivateKey(privateKey),
		cryptography.WithElliptic(
			cryptography.CreateElliptic(
				cryptography.WithEllipticPublicKey(ellipticPublic),
				cryptography.WithEllipticPrivateKey(ellipticPrivate))))

	authentication := postLogin(auth.CreateAuthenticationModel(
		auth.WithId(userId),
		auth.WithName(userLogin),
		auth.WithCryptography(userCrypto)))
	if authentication.Id != "" {
		key := cryptography.DiffieHellman(
			ellipticPrivate,
			authentication.Cryptography.Elliptic.PublicKey)
		authentication.Token = cryptography.DecryptAES(authentication.Token, key)
		authentication.Cryptography.PrivateKey = privateKey
	}
	return authentication
}

// func loginWithPKCS(userLogin string, userId string) auth.AuthenticationModel {
// 	publicKey := readFromFile(fmt.Sprintf("./keys/PUBLIC_KEY_%s", userLogin))
// 	privateKey := readFromFile(fmt.Sprintf("./keys/PRIVATE_KEY_%s", userLogin))
// 	sign := readFromFile(fmt.Sprintf("./keys/SIGN_%s", userLogin))

// 	ellipticPrivate, ellipticPublic := cryptography.GenerateEllipticCurveKeys()
// 	userCrypto := cryptography.CreateCryptography(
// 		cryptography.WithSign(sign),
// 		cryptography.WithPublicKey(publicKey),
// 		cryptography.WithElliptic(
// 			cryptography.CreateElliptic(
// 				cryptography.WithEllipticPublicKey(ellipticPublic))))

// 	authentication := postLogin(auth.CreateAuthenticationModel(
// 		auth.WithId(userId),
// 		auth.WithName(userLogin),
// 		auth.WithCryptography(userCrypto)))

// 	if authentication.Id != "" {
// 		key := cryptography.DiffieHellman(
// 			ellipticPrivate,
// 			assignedPeer.Cryptography.Elliptic.PublicKey)
// 		authentication.Token = cryptography.DecryptAES(authentication.Token, key)
// 		authentication.Cryptography.PrivateKey = privateKey
// 	}
// 	return authentication
// }

func loginWithPassword(userLogin string, password string) auth.AuthenticationModel {
	return postLogin(auth.CreateAuthenticationModel(
		auth.WithId(userLogin),
		auth.WithName(userLogin),
		auth.WithPassword(password)))
}

func postLogin(authenticationModel auth.AuthenticationModel) auth.AuthenticationModel {
	body, err := json.Marshal(authenticationModel)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return auth.CreateDefaultAuthenticationModel()
	}
	res := http.POST(
		assignedPeer,
		assignedPeer.Address+"/auth",
		string(body),
		&authenticationModel)
	if res.StatusCode == http.ACCEPTED {
		authenticationModel.Cookies = res.Cookies()
		return authenticationModel
	} else {
		return auth.CreateDefaultAuthenticationModel()
	}
}
