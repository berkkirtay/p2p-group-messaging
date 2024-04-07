package auth

type AuthenticationModel struct {
	Id       string `json:"id,omitempty" bson:"id,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Token    string `json:"token,omitempty" bson:"token,omitempty"`
}

type AuthenticationModelOption func(AuthenticationModel) AuthenticationModel

func WithId(id string) AuthenticationModelOption {
	return func(authenticationModel AuthenticationModel) AuthenticationModel {
		authenticationModel.Id = id
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
