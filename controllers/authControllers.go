package controllers

import (
	"contact-store/models"
	"contact-store/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //handle the error inline
	if err != nil {
		utils.Respond(w, 400, utils.Message("400", "Invalid request body"))
		return
	}
	if resp, ok := user.ValidateReqBody(""); !ok {
		fmt.Println(resp)
		utils.Respond(w, 400, resp)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Respond(w, 500, utils.Message("500", "Error hashing password"))
		return
	}
	user.Password = string(hashedPassword)
	models.GetDB().Create(user)
	if user.ID <= 0 {
		utils.Respond(w, 500, utils.Message("500", "Failed to create account"))
		return
	}
	utils.Respond(w, 201, utils.Message("201", "account successfully created"))
	return
}

func Login(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	response := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.Respond(w, 400, utils.Message("400", "Invalid request body"))
		return
	}
	if resp, ok := user.ValidateReqBody("login"); !ok {
		utils.Respond(w, 400, resp)
		return
	}
	tempUser := &models.User{}
	err = models.GetDB().Raw("select * from users where email=?", user.Email).First(tempUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Respond(w, 404, utils.Message("404", "Email address not found"))
		}
		utils.Respond(w, 400, utils.Message("400", "Connection error, Please retry"))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(tempUser.Password), []byte(user.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		utils.Respond(w, 400, utils.Message("400", "Invalid credentials, Please retry"))
		return
	}
	tokenClaims := &models.Token{
		UserId: tempUser.ID,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(1 * time.Hour).Unix(),
		}}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenClaims)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	response["token"] = tokenString
	response["id"] = tempUser.ID
	utils.Respond(w, 200, response)
	return
}
