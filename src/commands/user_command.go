// Copyright (c) 2024 Berk Kirtay

package commands

import (
	"encoding/json"
	"fmt"
	"main/infra/http"
	"main/services/auth"
	"main/services/user"
)

var sessionAuth auth.AuthenticationModel
var CurrentUser user.User

func HandleRegister(command []string) {

}

func HandleLogin(command []string) {
	if len(command) != 3 {
		fmt.Printf("Wrong usage.\n")
		return
	}

	userLogin := command[1]
	password := command[2]

	respBody := make([]user.User, 1)
	res := http.GET(assignedPeer.Address+"/users", &respBody, "id", userLogin, "name", userLogin)

	if res.StatusCode != http.OK {
		fmt.Printf("User %s does not exist.\n", userLogin)
		return
	}

	fmt.Printf("Hello %s! Your authentication process is ongoing..\n", respBody[0].Name)

	auth := auth.CreateAuthenticationModel(
		auth.WithId(userLogin),
		auth.WithName(userLogin),
		auth.WithPassword(password))
	login(&auth)

	if auth.Token != "" {
		sessionAuth = auth
		CurrentUser = respBody[0]
		http.InitializeService(sessionAuth.Cookies, sessionAuth.Id, sessionAuth.Token)
		fmt.Printf("You are authorized with the user-id: %s\n", auth.Id)
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
