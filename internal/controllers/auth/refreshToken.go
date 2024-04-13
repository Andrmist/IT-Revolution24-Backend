package auth

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	types.Response
	Tokens tokens      `json:"tokens"`
	User   domain.User `json:"user"`
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	server := r.Context().Value("server").(types.ServerContext)
	authorization := r.Header.Get("Authorization")
	parsedAuthorization := strings.Split(authorization, " ")
	if len(parsedAuthorization) == 2 {
		refreshToken, err := jwt.Parse(parsedAuthorization[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(server.Config.JWTSecret), nil
		})

		if err != nil {
			domain.HTTPError(w, r, http.StatusForbidden, errors.New("access token is invalid"))
			return
		}
		server.Log.Info(refreshToken.Valid)

		var user domain.User

		sub, err := refreshToken.Claims.GetSubject()
		if err != nil {
			domain.HTTPError(w, r, http.StatusUnauthorized, nil)
			return
		}

		server.DB.First(&user, "ID = ?", sub)
		exp, err := refreshToken.Claims.GetExpirationTime()

		if err != nil || exp.Sub(time.Now()).Seconds() <= 0 || !refreshToken.Valid || user.ID == 0 {
			domain.HTTPError(w, r, http.StatusUnauthorized, nil)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		tokens, err := generateTokens(user, server.Config)
		if err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}
		res := RefreshTokenResponse{
			User:   user,
			Tokens: tokens,
		}
		render.JSON(w, r, res)
	} else {
		domain.HTTPError(w, r, http.StatusBadRequest, nil)
	}
}

func GetUserFromAccessToken(server types.ServerContext, token string) (domain.User, error) {
	accessToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(server.Config.JWTSecret), nil
	})

	if err != nil {
		return domain.User{}, err
	}

	var user domain.User

	sub, err := accessToken.Claims.GetSubject()
	if err != nil {
		return domain.User{}, err
	}

	server.DB.First(&user, "ID = ?", sub)

	exp, err := accessToken.Claims.GetExpirationTime()

	if err != nil {
		return domain.User{}, err
	}

	if exp.Sub(time.Now()).Seconds() <= 0 || !accessToken.Valid || user.ID == 0 {
		return domain.User{}, errors.New("invalid or expired token")
	}

	return user, nil
}

func generateTokens(user domain.User, config types.Config) (tokens, error) {
	fmt.Println(user.ID)
	id := strconv.Itoa(int(user.ID))
	accessTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(24 * 14 * time.Hour).Unix(),
	})
	refreshTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(24 * 30 * time.Hour).Unix(),
	})
	accessToken, err := accessTokenJWT.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return tokens{}, err
	}
	refreshToken, err := refreshTokenJWT.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return tokens{}, err
	}
	return tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
