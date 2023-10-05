package models

import (
	"contact-store/utils"
	"strings"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	UserId      uint   `json:"user_id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_num"`
	Context     string `json:"context"`
}

func (contact *Contact) ValidateReqBody() (map[string]interface{}, bool) {
	if contact.UserId <= 0 {
		return utils.Message("400", "User id invalid"), false
	} else if len(strings.TrimSpace(contact.Name)) == 0 {
		return utils.Message("400", "Name requied"), false
	} else if len(strings.TrimSpace(contact.PhoneNumber)) == 0 {
		return utils.Message("400", "Phone num required"), false
	} else if !utils.IsPhoneValid(contact.PhoneNumber) {
		return utils.Message("400", "Invalid phone number"), false
	}
	return utils.Message("200", "Requirement passed"), true
}
