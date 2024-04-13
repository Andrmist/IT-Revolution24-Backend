package auth

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/mail.v2"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"math/rand"
	"net/http"
	"strconv"
)

type registerRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required"`
}

const emailFormat = `From: IT-Revolution24-Backend <%s>
To: %s
Subject: Account %s verification

Hi! Please, confirm %s account creation with this Link: %d`

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
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}
	user.Password = string(password)

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
	m.SetBody("text/plain", fmt.Sprintf("Hi! Please, confirm %s account creation with this code: %d", user.Role, user.AuthCode))
	if err := smtpD.DialAndSend(m); err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}

	if err := server.DB.Save(&user).Error; err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}

	responseTokens, err := generateTokens(user, server.Config)

	if err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}

	render.JSON(w, r, RefreshTokenResponse{
		Response: types.Response{},
		Tokens:   responseTokens,
		User:     user,
	})

}
