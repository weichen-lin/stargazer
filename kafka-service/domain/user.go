package domain

type User struct {
	name  string
	email string
	image string

	accessToken       string
	provider          string
	providerAccountId string
	scope             string
	authType          string
	tokenType         string
}

type UserEntity struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Image             string `json:"image"`
	AccessToken       string `json:"access_token"`
	Provider          string `json:"provider"`
	ProviderAccountId string `json:"providerAccountId"`
	Scope             string `json:"scope"`
	AuthType          string `json:"type"`
	TokenType         string `json:"token_type"`
}

func FromUserEntity(entity *UserEntity) *User {
	user := &User{
		name:              entity.Name,
		email:             entity.Email,
		image:             entity.Image,
		accessToken:       entity.AccessToken,
		provider:          entity.Provider,
		providerAccountId: entity.ProviderAccountId,
		scope:             entity.Scope,
		authType:          entity.AuthType,
		tokenType:         entity.TokenType,
	}
	return user
}