package routes

import (
	"contact-store/controllers"

	"github.com/gorilla/mux"
)

func AuthRoutes(router *mux.Router) {
	router.HandleFunc("/user/register", controllers.Register).Methods("POST")
	router.HandleFunc("/user/login", controllers.Login).Methods("POST")
}
