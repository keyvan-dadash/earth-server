package token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type AccessTokenDetail struct {
	AccessToken       string
	AccessTokenUUID   string
	AccessTokenExpire int64

	//below filed's is extra filed that we extract from token
	Username string
}

type RefreshTokenDetail struct {
	RefreshToken       string
	RefreshTokenUUID   string
	RefreshTokenExpire int64

	//below filed's is extra filed that we extract from token
	Username string
}

type TokenDetails struct {
	at *AccessTokenDetail
	rt *RefreshTokenDetail
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

	signedAccess, _ := at.SignedString([]byte(accessSecretKey))

	accessTokenDetail := &AccessTokenDetail{
		AccessToken:       signedAccess,
		AccessTokenUUID:   accessUUID,
		AccessTokenExpire: accessExpire,
	}

	//creating refresh
	refreshClaims := jwt.MapClaims{}
	refreshClaims["username"] = username

	refreshUUID := uuid.NewV4().String()
	refreshClaims["refresh_uuid"] = refreshUUID

	refreshExpire := time.Now().Add(2 * time.Hour).Unix()
	refreshClaims["expire_time"] = time.Now().Add(7 * 24 * time.Hour).Unix()

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedRefresh, _ := rt.SignedString([]byte(refreshSecretKey))

	refreshTokenDetail := &RefreshTokenDetail{
		RefreshToken:       signedRefresh,
		RefreshTokenUUID:   refreshUUID,
		RefreshTokenExpire: refreshExpire,
	}

	td := &TokenDetails{
		at: accessTokenDetail,
		rt: refreshTokenDetail,
	}

	return td, nil
}

func (t *TokenDetails) GetAccessToken() string {
	return t.at.AccessToken
}

func (t *TokenDetails) GetRefreshToken() string {
	return t.rt.RefreshToken
}

//CreateTokenBasedOnRefreshToken will renew access and refresh token but expire time of refresh token
//remain same because we want after speific period user login in to his account
//Note: this function will not delete previous refresh token uuid from database
//therefore you should delete previous refresh token uuid from database by yourself
func CreateTokenBasedOnRefreshToken(rt *RefreshTokenDetail) (*TokenDetails, error) {

	newToken, err := CreateToken(rt.Username)

	if err != nil {
		return nil, err
	}

	remainingRefreshTokenExpireTime := time.Unix(rt.RefreshTokenExpire, 0)
	now := time.Now()

	newToken.rt.RefreshTokenExpire = int64(remainingRefreshTokenExpireTime.Sub(now))

	return newToken, nil

}

//ExtractRefreshTokenFrom given refreshToken stirng
func ExtractRefreshTokenFrom(refreshToken string) (*RefreshTokenDetail, error) {

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

	if ok && token.Valid {
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

		return &RefreshTokenDetail{
			RefreshToken:       refreshToken,
			RefreshTokenUUID:   refreshUUID,
			RefreshTokenExpire: refreshExpireTime,
			Username:           username,
		}, nil
	}

	return nil, fmt.Errorf("token is invalid")

}

////ExtractAccessTokenFrom given accessToken stirng
func ExtractAccessTokenFrom(accessToken string) (*AccessTokenDetail, error) {

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method. err: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token is expired")
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, fmt.Errorf("token validation faild")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)

		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve refresh uuid")
		}

		username, ok := claims["username"].(string)

		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve username")
		}

		accessExpireTime, ok := claims["expire_time"].(int64)

		if !ok {
			return nil, fmt.Errorf("token is invalid because cannot retrieve expire time")
		}

		return &AccessTokenDetail{
			AccessToken:       accessToken,
			AccessTokenUUID:   accessUUID,
			AccessTokenExpire: accessExpireTime,
			Username:          username,
		}, nil
	}

	return nil, fmt.Errorf("token is invalid")

}
