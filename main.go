package main

import (
	"contact-store/models"
	"contact-store/routes"

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
	routes.AuthRoutes(router)
	subRouter := router.PathPrefix("/auth").Subrouter()
	routes.ContactRoutes(subRouter)

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
