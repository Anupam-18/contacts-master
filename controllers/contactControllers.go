package controllers

import (
	"contact-store/models"
	"contact-store/utils"
	"encoding/json"
	"fmt"
	"net/http"
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
	// id := r.Context().Value("user_id").(uint)
	query := r.URL.Query()
	id := query.Get("user_id")
	// fmt.Println(id)
	contacts := make([]*models.Contact, 0)
	// err = models.GetDB().Raw("select * from users where email=?", user.Email).First(tempUser).Error
	err := models.GetDB().Table("contacts").Where("user_id=?", id).Find(&contacts).Error
	if err != nil {
		fmt.Println(err)
		utils.Respond(w, 500, utils.Message("500", "Failed to fetch contacts"))
	}
	utils.Respond(w, 200, map[string]interface{}{
		"data": map[string]interface{}{
			"status":   "200",
			"contacts": contacts,
		},
	})
}
