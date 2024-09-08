// Copyright (c) 2024 Berk Kirtay

package auth

import (
	"main/infra/cryptography"
	"net/http"
)

type AuthenticationModel struct {
	Id           string                     `json:"id,omitempty" bson:"id,omitempty"`
	Name         string                     `json:"name,omitempty" bson:"name,omitempty"`
	Password     string                     `json:"password,omitempty" bson:"password,omitempty"`
	Token        string                     `json:"token,omitempty" bson:"token,omitempty"`
	Cryptography *cryptography.Cryptography `json:"cryptography,omitempty" bson:"cryptography,omitempty"`
	Cookies      []*http.Cookie             `json:"cookies,omitempty" bson:"cookies,omitempty"`
}

type AuthenticationModelOption func(AuthenticationModel) AuthenticationModel

func WithId(id string) AuthenticationModelOption {
	return func(authenticationModel AuthenticationModel) AuthenticationModel {
		authenticationModel.Id = id
		return authenticationModel
	}
}

func WithName(name string) AuthenticationModelOption {
	return func(authenticationModel AuthenticationModel) AuthenticationModel {
		authenticationModel.Name = name
		return authenticationModel
	}
}

func WithPassword(password string) AuthenticationModelOption {
	return func(authenticationModel AuthenticationModel) AuthenticationModel {
		authenticationModel.Password = password
		return authenticationModel
	}
}

func WithToken(token string) AuthenticationModelOption {
	return func(authenticationModel AuthenticationModel) AuthenticationModel {
		authenticationModel.Token = token
		return authenticationModel
	}
}

func WithCryptography(cryptography *cryptography.Cryptography) AuthenticationModelOption {
	return func(authenticationModel AuthenticationModel) AuthenticationModel {
		authenticationModel.Cryptography = cryptography
		return authenticationModel
	}
}

func WithCookies(cookies []*http.Cookie) AuthenticationModelOption {
	return func(authenticationModel AuthenticationModel) AuthenticationModel {
		authenticationModel.Cookies = cookies
		return authenticationModel
	}
}

func CreateDefaultAuthenticationModel() AuthenticationModel {
	return AuthenticationModel{}
}

func CreateAuthenticationModel(options ...AuthenticationModelOption) AuthenticationModel {
	authenticationModel := CreateDefaultAuthenticationModel()

	for _, option := range options {
		authenticationModel = option(authenticationModel)
	}

	return authenticationModel
}
