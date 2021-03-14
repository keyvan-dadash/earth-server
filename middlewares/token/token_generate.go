package token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type TokenDetails struct {
	AccessToken        string
	RefreshToken       string
	AccessTokenUUID    string
	RefreshTokenUUID   string
	AccessTokenExpire  int64
	RefreshTokenExpire int64
}

//CreateToken is function that generate token base on given username and return generated token
func CreateToken(username string) (*TokenDetails, error) {

	accessSecretKey := os.Getenv("ACCESS_SECRET_KEY") //you can pass keys and gain better performance bcs every time you are reading from env
	if len(accessSecretKey) == 0 {
		return nil, errors.New("cannot get access secret key from env")
	}

	refreshSecretKey := os.Getenv("REFRESH_SECRET_KEY")
	if len(refreshSecretKey) == 0 {
		return nil, errors.New("cannot get refresh secret key from env")
	}

	//TODO: its better to determine and set access and refresh expire time from env variable

	//creating access
	accessClaims := jwt.MapClaims{}
	accessClaims["username"] = username

	accessUUID := uuid.NewV4().String()
	accessClaims["access_uuid"] = accessUUID

	accessExpire := time.Now().Add(2 * time.Hour).Unix()
	accessClaims["expire_time"] = time.Now().Add(2 * time.Hour).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	//creating refresh
	refreshClaims := jwt.MapClaims{}
	refreshClaims["username"] = username

	refreshUUID := uuid.NewV4().String()
	refreshClaims["refresh_uuid"] = refreshUUID

	refreshExpire := time.Now().Add(2 * time.Hour).Unix()
	refreshClaims["expire_time"] = time.Now().Add(7 * 24 * time.Hour).Unix()

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedAccess, _ := at.SignedString([]byte(accessSecretKey))
	signedRefresh, _ := rt.SignedString([]byte(refreshSecretKey))

	td := &TokenDetails{
		AccessToken:        signedAccess,
		RefreshToken:       signedRefresh,
		AccessTokenUUID:    accessUUID,
		RefreshTokenUUID:   refreshUUID,
		AccessTokenExpire:  accessExpire,
		RefreshTokenExpire: refreshExpire,
	}

	return td, nil
}

func CreateTokenBasedOnRefreshToken(refreshToken string) (*TokenDetails, error) {

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method. err: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("REFRESH_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token is expired")
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, fmt.Errorf("token validation faild")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid { //TODO: extracting refresh and access token login must be sperate
		refreshUUID, ok := claims["refresh_uuid"].(string)

		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve refresh uuid")
		}

		username, ok := claims["username"].(string)

		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve username")
		}

		refreshExpireTime, ok := claims["expire_time"].(int64)

		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve expire time")
		}

		if _, err := DeleteRefreshToken(refreshUUID); err != nil { //TODO: we have some passing db connection we must resolve this
			return nil, err
		}

		rt := time.Unix(refreshExpireTime, 0)
		now := time.Now()

		newToken, err := CreateToken(username)

		if err != nil {
			return nil, err
		}

		newToken.RefreshTokenExpire = int64(rt.Sub(now))

		return newToken, nil
	}
}
