package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/RedBrand88/breadmachine/utility"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	opt := option.WithCredentialsFile("/etc/breadmachine/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("failed to create firebase app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	defer client.Close()

	log.Println("Firebase initialized successfully")

	iter := client.Collection("Recipes").Documents(ctx)
	defer iter.Stop()

	updated := 0
	skipped := 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("failed to iterate: %v", err)
		}

		data := doc.Data()
		ingredients, ok := data["ingredients"].([]interface{})
		if !ok {
			skipped++
			continue
		}

		updatedIngredients := make([]interface{}, len(ingredients))
		for i, ing := range ingredients {
			ingMap, ok := ing.(map[string]interface{})
			if !ok {
				updatedIngredients[i] = ing
				continue
			}

			// Only backfill if missing or zero
			existing, _ := ingMap["densityGPerMl"].(float64)
			if existing == 0 {
				name, _ := ingMap["ingredientName"].(string)
				density := utility.LookupDensity(name)
				ingMap["densityGPerMl"] = density
				fmt.Printf("  %s → %.3f\n", name, density)
			}

			updatedIngredients[i] = ingMap
		}

		_, err = doc.Ref.Update(ctx, []firestore.Update{
			{Path: "ingredients", Value: updatedIngredients},
		})
		if err != nil {
			log.Fatalf("failed to update doc %s: %v", doc.Ref.ID, err)
		}

		updated++
		fmt.Printf("updated recipe: %s\n", doc.Ref.ID)
	}

	fmt.Printf("\ndone — updated %d recipes, skipped %d\n", updated, skipped)
}
