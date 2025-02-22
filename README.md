# BreadMachine API

A simple Go API that manages bread recipes and calculates baker's percentages.

## Features
- Create, read, update, and delete recipes
- Convert recipes to baker's math for scaling
- Store recipes in Firebase

## API Endpoints

| Method | Endpoint | Description |
|--------|---------|-------------|
| GET | `/recipes` | Fetch all recipes |
| GET | `/recipes/{id}` | Get a specific recipe |
| POST | `/recipes` | Create a new recipe |
| DELETE | `/recipes/{id}` | Delete a recipe |

### Prerequisites
- Go 1.21+
- Firebase setup

