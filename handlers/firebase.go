package handlers

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/RedBrand88/breadmachine/models"
)

var client *firestore.Client

// InitFirebase initializes the Firebase Firestore client
func InitFirebase() error {
	opt := option.WithCredentialsFile("/etc/breadmachine/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	client, err = app.Firestore(context.Background())
	if err != nil {
		return err
	}

	log.Println("Firebase initialized successfully")
	return nil
}

// SaveRecipe saves a recipe to Firebase
func SaveRecipe(recipe models.Recipe) (string, error) {
	docRef, _, err := client.Collection("Recipes").Add(context.Background(), recipe)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

// FetchAllRecipesFromFirebase retrieves all recipes from Firestore
func FetchAllRecipesFromFirebase() ([]models.Recipe, error) {
	var recipes []models.Recipe
	iter := client.Collection("Recipes").Documents(context.Background())

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var recipe models.Recipe
		doc.DataTo(&recipe)
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}
