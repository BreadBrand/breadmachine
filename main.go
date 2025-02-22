package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RedBrand88/breadmachine/handlers"
)

func main() {
	fmt.Println("starting Breadmachine API...")

	//initialize Firebase
	err := handlers.InitFirebase()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	//define API routes
	http.HandleFunc("/recipes", handlers.GetAllRecipes)
	http.HandleFunc("/recipes/", handlers.RecipeHandler)
	http.HandleFunc("/recipes/create", handlers.CreateRecipe)

	//start the server
	log.Println("server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
