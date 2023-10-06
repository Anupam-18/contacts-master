package controllers

import (
	"contact-store/models"
	"contact-store/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func CreateContact(w http.ResponseWriter, r *http.Request) {
	contact := &models.Contact{}
	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		utils.Respond(w, 400, utils.Message("400", "Invalid request body"))
		return
	}
	if resp, ok := contact.ValidateReqBody(); !ok {
		utils.Respond(w, 400, resp)
		return
	}
	tempContact := &models.Contact{}
	err = models.GetDB().Raw("select * from contacts where phone_number=?", contact.PhoneNumber).First(tempContact).Error

	if tempContact.PhoneNumber == contact.PhoneNumber {
		utils.Respond(w, 201, utils.Message("201", "contact already exists"))
		return
	}
	models.GetDB().Create(contact)
	if contact.ID <= 0 {
		utils.Respond(w, 500, utils.Message("500", "Failed to create contact"))
		return
	}
	utils.Respond(w, 201, map[string]interface{}{
		"inserted_id": contact.ID,
	})
	return
}

func GetAllContacts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("user_id")
	contacts := make([]*models.Contact, 0)
	err := models.GetDB().Table("contacts").Where("user_id=?", id).Find(&contacts).Error
	if err != nil {
		fmt.Println(err)
		utils.Respond(w, 500, utils.Message("500", "Failed to fetch contacts"))
		return
	}
	utils.Respond(w, 200, map[string]interface{}{
		"data": map[string]interface{}{
			"status":   "200",
			"contacts": contacts,
		},
	})
}

func DeleteContact(w http.ResponseWriter, r *http.Request) {

	contact := &models.Contact{}
	err := json.NewDecoder(r.Body).Decode(contact)
	fmt.Println(contact.PhoneNumber)
	if err != nil {
		utils.Respond(w, 400, utils.Message("400", "Invalid request body"))
		return
	}
	if !utils.IsPhoneValid(contact.PhoneNumber) {
		utils.Respond(w, 400, utils.Message("400", "Invalid phone number"))
		return
	}
	reqContact := &models.Contact{}
	err = models.GetDB().Raw("select * from contacts where phone_number=?", contact.PhoneNumber).First(reqContact).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Respond(w, 404, utils.Message("404", "record does not exist, create it first"))
			return
		}
		utils.Respond(w, 500, utils.Message("500", "Connection error, Please retry"))
		return
	}
	if err := models.GetDB().Delete(&reqContact).Error; err != nil {
		utils.Respond(w, 500, utils.Message("500", "Internal server error"))
		return
	}
	utils.Respond(w, 204, utils.Message("204", "contact deleted"))
	return
}

func UpdateContact(w http.ResponseWriter, r *http.Request) {

	contact := &models.Contact{}
	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		utils.Respond(w, 400, utils.Message("400", "Invalid request body"))
		return
	}
	if len(strings.TrimSpace(contact.PhoneNumber)) == 0 {
		utils.Respond(w, 400, utils.Message("400", "Phone number required"))
		return
	}
	if !utils.IsPhoneValid(contact.PhoneNumber) {
		utils.Respond(w, 400, utils.Message("400", "Invalid phone number"))
		return
	}
	reqContact := &models.Contact{}
	err = models.GetDB().Raw("select * from contacts where phone_number=?", contact.PhoneNumber).First(reqContact).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Respond(w, 404, utils.Message("404", "record does not exist, create it first"))
			return
		}
		utils.Respond(w, 500, utils.Message("500", "Internal server error, Please retry"))
		return
	}
	if len(strings.TrimSpace(contact.Name)) != 0 {
		reqContact.Name = contact.Name
	}
	reqContact.Context = contact.Context
	if err := models.GetDB().Save(reqContact).Error; err != nil {
		utils.Respond(w, 500, utils.Message("500", "Internal server error, Please retry"))
		return
	}
	utils.Respond(w, 201, map[string]interface{}{
		"updated_data": reqContact,
	})
	return
}
