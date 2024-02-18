package user

import (
	"main/services/audit"
	"main/services/cryptography"
)

type User struct {
	Id        string                  `json:"id,omitempty" bson:"id,omitempty"`
	Name      string                  `json:"name,omitempty" bson:"name,omitempty"`
	Password  string                  `json:"password,omitempty" bson:"password,omitempty"`
	Role      string                  `json:"role,omitempty" bson:"role,omitempty"`
	Signature *cryptography.Signature `json:"signature,omitempty" bson:"signature,omitempty"`
	Actions   []*Action               `json:"actions,omitempty" bson:"actions,omitempty"`
	Features  []*Feature              `json:"features,omitempty" bson:"features,omitempty"`
	Audit     *audit.Audit            `json:"audit,omitempty" bson:"audit,omitempty"`
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

func WithSignature(signature *cryptography.Signature) UserOption {
	return func(user User) User {
		user.Signature = signature
		return user
	}
}

func WithActions(actions []*Action) UserOption {
	return func(user User) User {
		user.Actions = actions
		return user
	}
}

func WithFeatures(features []*Feature) UserOption {
	return func(user User) User {
		user.Features = features
		return user
	}
}

func WithAudit(audit *audit.Audit) UserOption {
	return func(user User) User {
		user.Audit = audit
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
