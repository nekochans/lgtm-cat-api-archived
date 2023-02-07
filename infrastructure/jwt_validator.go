package infrastructure

import (
	"context"
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/nekochans/lgtm-cat-api/constants"
	"github.com/nekochans/lgtm-cat-api/derrors"
)

type JwtValidator struct {
	issuer string
	jwkSet jwk.Set
}

func NewJwtValidator() (validator JwtValidator, err error) {
	defer derrors.Wrap(&err, "NewJwtValidator()")

	jwksUrl := constants.JwksUri()
	iss := fmt.Sprintf(
		"https://cognito-idp.%v.amazonaws.com/%v",
		constants.GetRegion(),
		constants.GetCognitoUserPoolId(),
	)

	c := jwk.NewCache(context.Background())

	if err := c.Register(""); err != nil {
		return JwtValidator{}, err
	}

	return JwtValidator{
		issuer: iss,
		jwkSet: jwk.NewCachedSet(c, jwksUrl),
	}, nil
}

func (v JwtValidator) Validate(accessToken string) (err error) {
	defer derrors.Wrap(&err, "JwtValidator.Validate(%s)", accessToken)

	token, err := jwt.ParseString(
		accessToken,
		jwt.WithKeySet(v.jwkSet),
	)
	if err != nil {
		return err
	}

	err = jwt.Validate(
		token,
		jwt.WithIssuer(v.issuer),
		jwt.WithClaimValue("token_use", "access"),
	)
	if err != nil {
		return err
	}

	return nil
}
