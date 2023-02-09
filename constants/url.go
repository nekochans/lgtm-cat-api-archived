package constants

import "fmt"

func JwksUri() string {
	return fmt.Sprintf(
		"https://cognito-idp.%v.amazonaws.com/%v/.well-known/jwks.json",
		GetRegion(),
		GetCognitoUserPoolId(),
	)
}
