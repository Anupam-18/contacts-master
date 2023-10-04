package main

import (
	"contact-store/controllers"
	"contact-store/models"

	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func init() {
	models.InitDB()
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	// router.Use(middleware.HandleJwtAuth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println(port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print("Error creating server", err)
		return
	}
}

// 	router.NotFoundHandler = app.NotFoundHandler
