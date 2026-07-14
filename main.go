package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BreadBrand/breadmachine/handlers"
)

func main() {
	fmt.Println("starting Breadmachine API...")

	//initialize Firebase
	err := handlers.InitFirebase()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	//define API routes
	http.HandleFunc("GET /api/recipes", handlers.GetAllRecipes)
	http.HandleFunc("POST /api/recipes", handlers.CreateRecipe)
	http.HandleFunc("POST /api/recipes/parse", handlers.ParseHandler)
	http.HandleFunc("GET /api/recipes/{id}", handlers.GetRecipe)
	http.HandleFunc("DELETE /api/recipes/{id}", handlers.DeleteRecipe)

	//start the server
	log.Println("server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
