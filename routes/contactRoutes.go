package routes

import (
	"contact-store/controllers"
	"contact-store/middleware"

	"github.com/gorilla/mux"
)

func ContactRoutes(router *mux.Router) {
	router.Use(middleware.HandleJwtAuth)
	router.HandleFunc("/create-contact", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/contacts", controllers.GetAllContacts).Methods("GET")
	router.HandleFunc("/contacts", controllers.DeleteContact).Methods("DELETE")
	router.HandleFunc("/contacts", controllers.UpdateContact).Methods("PATCH")
}
