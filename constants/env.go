package constants

import "os"

func GetRegion() string {
	return os.Getenv("REGION")
}

func GetCognitoUserPoolId() string {
	return os.Getenv("COGNITO_USER_POOL_ID")
}
