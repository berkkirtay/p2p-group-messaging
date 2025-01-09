// Copyright (c) 2024 Berk Kirtay

package http

import (
	"net/http"
)

type HeaderModel struct {
	Connection    string         `json:"Connection,omitempty" bson:"Connection,omitempty"`
	Authorization string         `json:"Authorization,omitempty" bson:"Authorization,omitempty"`
	Session       string         `json:"Session,omitempty" bson:"Session,omitempty"`
	PublicKey     string         `json:"PublicKey,omitempty" bson:"PublicKey,omitempty"`
	Cookie        []*http.Cookie `json:"Cookie,omitempty" bson:"Cookie,omitempty"`
	ContentType   string         `json:"Content-Type,omitempty" bson:"Content-Type,omitempty"`
}

type HeaderModelOption func(HeaderModel) HeaderModel

func WithConnection(connection string) HeaderModelOption {
	return func(headerModel HeaderModel) HeaderModel {
		headerModel.Connection = connection
		return headerModel
	}
}

func WithAuthorization(authorization string) HeaderModelOption {
	return func(headerModel HeaderModel) HeaderModel {
		headerModel.Authorization = authorization
		return headerModel
	}
}

func WithSession(session string) HeaderModelOption {
	return func(headerModel HeaderModel) HeaderModel {
		headerModel.Session = session
		return headerModel
	}
}

func WithPublicKey(publicKey string) HeaderModelOption {
	return func(headerModel HeaderModel) HeaderModel {
		headerModel.PublicKey = publicKey
		return headerModel
	}
}

func WithCookie(cookie []*http.Cookie) HeaderModelOption {
	return func(headerModel HeaderModel) HeaderModel {
		headerModel.Cookie = cookie
		return headerModel
	}
}

func WithContentType(contentType string) HeaderModelOption {
	return func(headerModel HeaderModel) HeaderModel {
		headerModel.ContentType = contentType
		return headerModel
	}
}

func WithHeader(headerModelCopy HeaderModel) HeaderModelOption {
	return func(headerModel HeaderModel) HeaderModel {
		headerModel.Connection = headerModelCopy.Connection
		headerModel.Authorization = headerModelCopy.Authorization
		headerModel.Session = headerModelCopy.Session
		headerModel.PublicKey = headerModelCopy.PublicKey
		headerModel.Cookie = headerModelCopy.Cookie
		headerModel.ContentType = headerModelCopy.ContentType
		return headerModel
	}
}

func CreateDefaultHeaderModel() HeaderModel {
	return HeaderModel{}
}

func CreateHeaderModel(options ...HeaderModelOption) HeaderModel {
	headerModel := CreateDefaultHeaderModel()

	for _, option := range options {
		headerModel = option(headerModel)
	}

	return headerModel
}
