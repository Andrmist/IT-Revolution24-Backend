package auth

import (
	"fmt"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	gomail "gopkg.in/mail.v2"
)

type registerRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	server := r.Context().Value("server").(types.ServerContext)
	var req registerRequest
	err := render.DecodeJSON(r.Body, &req)
	smtpPort, _ := strconv.Atoi(server.Config.SMTPPort)
	smtpD := gomail.NewDialer(server.Config.SMTPHost, smtpPort, server.Config.SMTPUser, server.Config.SMTPPass)

	if err != nil {
		domain.HTTPError(w, r, http.StatusBadRequest, errors.New("failed to decode body"))
		return
	}

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		domain.HTTPError(w, r, http.StatusBadRequest, domain.ValidationError(validateErr))
		return
	}

	var existingUser domain.User
	server.DB.First(&existingUser, "name = ?", req.Username)

	if existingUser.Name != "" {
		domain.HTTPError(w, r, http.StatusBadRequest, errors.New("user with this username already exists"))
		return
	}
	existingUser = domain.User{}

	var user domain.User
	user.Name = req.Username
	user.Email = req.Email
	user.Role = req.Role
	user.Password = req.Password
	if user.Role == "child" {
		user.Balance = types.STANDARD_BALANCE
	}

	//password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	//if err != nil {
	//	domain.HTTPInternalServerError(w, r, err)
	//	return
	//}
	//user.Password = string(password)

	for {
		user.AuthCode = 10000 + rand.Intn(99999-10000)
		var existingUserCode domain.User
		server.DB.First(existingUserCode, "auth_code = ?", user.AuthCode)
		if existingUserCode.AuthCode != user.AuthCode {
			break
		}
	}

	if user.Role == "parent" {
		server.DB.First(&existingUser, "email = ? AND role = ?", req.Email, "parent")
		if existingUser.Email != "" {
			domain.HTTPError(w, r, http.StatusBadRequest, errors.New("parent with this email already exists"))
			return
		}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", server.Config.SMTPUser)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", fmt.Sprintf("Account %s verification", user.Name))
	if req.Role == "parent" {
		m.SetBody("text/plain", fmt.Sprintf("Hi! Please, confirm %s account creation with this code: %d", user.Role, user.AuthCode))
	} else {
		m.SetBody("text/plain", fmt.Sprintf("Hi! Please, confirm %s account creation with this link: https://hackaton.dev.m0e.space/code/%d", user.Role, user.AuthCode))
	}
	if err := smtpD.DialAndSend(m); err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}

	if err := server.DB.Save(&user).Error; err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}

	if user.Role == "child" {
		if err := server.DB.Create(&domain.Pet{
			Type:      types.TYPE_FISH,
			Sex:       types.SEX_MALE,
			Satiety:   100,
			LoveMeter: 0,
			Cost:      types.FISH_COST,
			UserID:    user.ID,
		}).Error; err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}

		if err := server.DB.Create(&domain.Pet{
			Type:      types.TYPE_FISH,
			Sex:       types.SEX_FEMALE,
			Satiety:   100,
			LoveMeter: 0,
			Cost:      types.FISH_COST,
			UserID:    user.ID,
		}).Error; err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}

		if err := server.DB.Create(&domain.Pet{
			Type:      types.TYPE_SNAIL,
			Sex:       types.SEX_MALE,
			Satiety:   100,
			LoveMeter: 0,
			Cost:      types.SNAIL_COST,
			UserID:    user.ID,
		}).Error; err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}

		if err := server.DB.Create(&domain.Pet{
			Type:      types.TYPE_SNAIL,
			Sex:       types.SEX_FEMALE,
			Satiety:   100,
			LoveMeter: 0,
			Cost:      types.SNAIL_COST,
			UserID:    user.ID,
		}).Error; err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}

		if err := server.DB.Create(&domain.Pet{
			Type:      types.TYPE_SHRIMP,
			Sex:       types.SEX_MALE,
			Satiety:   100,
			LoveMeter: 0,
			Cost:      types.SHRIMP_COST,
			UserID:    user.ID,
		}).Error; err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}

		if err := server.DB.Create(&domain.Pet{
			Type:      types.TYPE_SHRIMP,
			Sex:       types.SEX_FEMALE,
			Satiety:   100,
			LoveMeter: 0,
			Cost:      types.SHRIMP_COST,
			UserID:    user.ID,
		}).Error; err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}
	}

	responseTokens, err := generateTokens(user, server.Config)

	if err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}

	render.JSON(w, r, RefreshTokenResponse{
		Response: types.Response{},
		Tokens:   responseTokens,
	})

}
