package domain

type JwtValidator interface {
	Validate(accessToken string) (err error)
}
