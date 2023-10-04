package models

import (
	"contact-store/utils"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}
type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

func (user *User) ValidateReqBody(flag string) (map[string]interface{}, bool) {

	if len(strings.TrimSpace(user.Email)) == 0 {
		return utils.Message("400", "Email requied"), false
	} else if !utils.IsEmailValid(user.Email) {
		return utils.Message("400", "Email invalid"), false
	} else if len(strings.TrimSpace(user.Password)) == 0 {
		return utils.Message("400", "Password required"), false
	} else if len(strings.TrimSpace(user.Password)) < 6 {
		return utils.Message("400", "Password should have at least 6 chars"), false
	}

	tempUser := &User{}
	err := GetDB().Raw("select * from users where email=?", user.Email).First(tempUser).Error
	// fmt.Println(gorm.ErrRecordNotFound, err != gorm.ErrRecordNotFound)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		return utils.Message("400", "Connection Error. Please retry"), false
	}
	if flag == "login" {
		return utils.Message("200", "Requirement passed"), true
	}
	if tempUser.Email != "" {
		return utils.Message("400", "Email already in use"), false
	}
	return utils.Message("200", "Requirement passed"), true
}
