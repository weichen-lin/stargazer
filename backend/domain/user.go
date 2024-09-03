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

func (u *User) ToUserEntity() *UserEntity {
	return &UserEntity{
		Name:              u.name,
		Email:             u.email,
		Image:             u.image,
		AccessToken:       u.accessToken,
		Provider:          u.provider,
		ProviderAccountId: u.providerAccountId,
		Scope:             u.scope,
		AuthType:          u.authType,
		TokenType:         u.tokenType,
	}
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Image() string {
	return u.image
}

func (u *User) AccessToken() string {
	return u.accessToken
}

func (u *User) Provider() string {
	return u.provider
}

func (u *User) ProviderAccountId() string {
	return u.providerAccountId
}

func (u *User) Scope() string {
	return u.scope
}

func (u *User) AuthType() string {
	return u.authType
}

func (u *User) TokenType() string {
	return u.tokenType
}
