// Copyright (c) 2024 Berk Kirtay

package user

import (
	"main/infra/cryptography"
	"main/services/audit"
)

type User struct {
	Id           string                     `json:"id,omitempty" bson:"id,omitempty"`
	Name         string                     `json:"name,omitempty" bson:"name,omitempty"`
	Password     string                     `json:"password,omitempty" bson:"password,omitempty"`
	Role         string                     `json:"role,omitempty" bson:"role,omitempty"`
	Cryptography *cryptography.Cryptography `json:"cryptography,omitempty" bson:"cryptography,omitempty"`
	Actions      []*Action                  `json:"actions,omitempty" bson:"actions,omitempty"`
	Audit        *audit.Audit               `json:"audit,omitempty" bson:"audit,omitempty"`
	IsPeer       bool                       `json:"isPeer,omitempty" bson:"isPeer,omitempty"`
}

type UserOption func(User) User

func WithId(id string) UserOption {
	return func(user User) User {
		user.Id = id
		return user
	}
}

func WithName(name string) UserOption {
	return func(user User) User {
		user.Name = name
		return user
	}
}

func WithPassword(password string) UserOption {
	return func(user User) User {
		user.Password = password
		return user
	}
}

func WithRole(role string) UserOption {
	return func(user User) User {
		user.Role = role
		return user
	}
}

func WithCryptography(cryptography *cryptography.Cryptography) UserOption {
	return func(user User) User {
		user.Cryptography = cryptography
		return user
	}
}

func WithActions(actions []*Action) UserOption {
	return func(user User) User {
		user.Actions = actions
		return user
	}
}

func WithAudit(audit *audit.Audit) UserOption {
	return func(user User) User {
		user.Audit = audit
		return user
	}
}

func WithIsPeer(isPeer bool) UserOption {
	return func(user User) User {
		user.IsPeer = isPeer
		return user
	}
}

func WithUser(newUser User) UserOption {
	return func(user User) User {
		user.Id = newUser.Id
		user.Name = newUser.Name
		user.Password = newUser.Password
		user.Role = newUser.Role
		user.Cryptography = newUser.Cryptography
		user.Actions = newUser.Actions
		user.Audit = newUser.Audit
		user.IsPeer = newUser.IsPeer
		return user
	}
}

func CreateDefaultUser() User {
	return User{}
}

func CreateUser(options ...UserOption) User {
	user := CreateDefaultUser()

	for _, option := range options {
		user = option(user)
	}

	return user
}
