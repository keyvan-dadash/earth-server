package authentication

import (
	"errors"
	"os"

	"github.com/sod-lol/earth-server/internal/app/authentication/models"
)

//CreateToken is function that generate token base on given User and return generated token
func CreateToken(user models.User) (string, error) {

	accessSecretKey := os.Getenv("ACCESS_SECRET_KEY") //you can pass keys and gain better performance
	if len(accessSecretKey) == 0 {
		return "", errors.New("cannot get access secret key from env")
	}

	refreshSecretKey := os.Getenv("REFRESH_SECRET_KEY")
	if len(refreshSecretKey) == 0 {
		return "", errors.New("cannot get refresh secret key from env")
	}

	return "", nil
}
