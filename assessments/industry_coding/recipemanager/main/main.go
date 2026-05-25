package main

import (
	"fmt"
	"strconv"

	"github.com/Bruwald/CodeSignalPractice/assessments/industry_coding/recipemanager"
)

func main() {
	recipeManager := recipemanager.NewRecipeManager()

	queries := [][]string{
		{"CREATE_USER", "1", "alice"},
		{"CREATE_USER", "2", "bob"},
		{"CREATE_USER", "3", "charlie"},
		{"CREATE_USER", "4", "diana"},
		{"CREATE_USER", "5", "eve"},
		{"ADD_RECIPE", "6", "alice", "pasta", "flour,eggs,salt", "Mix flour and eggs, cook in boiling water"},
		{"ADD_RECIPE", "7", "alice", "pizza", "flour,tomato,cheese", "Mix dough, add toppings, bake at 450F"},
		{"ADD_RECIPE", "8", "bob", "salad", "lettuce,tomato,olive_oil", "Chop vegetables, mix with olive oil"},
		{"ADD_RECIPE", "9", "bob", "soup", "chicken,broth,vegetables", "Simmer chicken in broth with vegetables"},
		{"ADD_RECIPE", "10", "charlie", "bread", "flour,water,salt,yeast", "Knead dough, let rise, bake"},
		{"ADD_RECIPE", "11", "charlie", "cake", "flour,sugar,eggs,butter", "Mix ingredients, bake at 350F"},
		{"ADD_RECIPE", "12", "diana", "cookies", "flour,sugar,butter,chocolate", "Cream butter and sugar, add flour and chocolate"},
		{"ADD_RECIPE", "13", "eve", "ice_cream", "cream,milk,sugar,vanilla", "Mix ingredients, churn until frozen"},
		{"RATE_RECIPE", "14", "bob", "pasta", "5"},
		{"RATE_RECIPE", "15", "charlie", "pasta", "4"},
		{"RATE_RECIPE", "16", "diana", "pasta", "5"},
		{"RATE_RECIPE", "17", "eve", "pasta", "3"},
		{"RATE_RECIPE", "18", "alice", "salad", "3"},
		{"RATE_RECIPE", "19", "charlie", "salad", "4"},
		{"RATE_RECIPE", "20", "diana", "salad", "5"},
		{"RATE_RECIPE", "21", "alice", "soup", "2"},
		{"RATE_RECIPE", "22", "eve", "soup", "4"},
		{"RATE_RECIPE", "23", "bob", "bread", "4"},
		{"RATE_RECIPE", "24", "diana", "bread", "5"},
		{"RATE_RECIPE", "25", "eve", "bread", "3"},
		{"RATE_RECIPE", "26", "alice", "cake", "5"},
		{"RATE_RECIPE", "27", "bob", "cake", "4"},
		{"RATE_RECIPE", "28", "diana", "cake", "5"},
		{"RATE_RECIPE", "29", "charlie", "cookies", "4"},
		{"RATE_RECIPE", "30", "alice", "cookies", "3"},
		{"RATE_RECIPE", "31", "eve", "cookies", "4"},
		{"RATE_RECIPE", "32", "bob", "ice_cream", "5"},
		{"RATE_RECIPE", "33", "charlie", "ice_cream", "5"},
		{"RATE_RECIPE", "34", "diana", "ice_cream", "4"},
		{"FAVORITE_RECIPE", "35", "alice", "salad"},
		{"FAVORITE_RECIPE", "36", "alice", "bread"},
		{"FAVORITE_RECIPE", "37", "bob", "pasta"},
		{"FAVORITE_RECIPE", "38", "bob", "pizza"},
		{"FAVORITE_RECIPE", "39", "charlie", "pasta"},
		{"FAVORITE_RECIPE", "40", "charlie", "bread"},
		{"FAVORITE_RECIPE", "41", "diana", "pasta"},
		{"FAVORITE_RECIPE", "42", "diana", "pizza"},
		{"FAVORITE_RECIPE", "43", "diana", "salad"},
		{"FAVORITE_RECIPE", "44", "eve", "ice_cream"},
		{"FAVORITE_RECIPE", "45", "eve", "cookies"},
		{"GET_TOP_RECIPES", "46", "5"},
		{"SEARCH_BY_INGREDIENT", "47", "flour"},
		{"GET_USER_FAVORITES", "48", "diana"},
		{"GET_RECOMMENDATIONS", "49", "alice", "3"},
		{"GET_LEADERBOARD", "50", "5"},
	}

	processQueries(queries, recipeManager)
	printRecipeManager(recipeManager)
}

func processQueries(queries [][]string, recipeManager recipemanager.RecipeManager) {
	for _, query := range queries {
		cmd := query[0]

		switch cmd {
		case "CREATE_USER":
			timestamp, _ := strconv.Atoi(query[1])
			userID := recipemanager.UserID(query[2])
			_ = recipeManager.CreateUser(int64(timestamp), userID)

		case "ADD_RECIPE":
			timestamp, _ := strconv.Atoi(query[1])
			userID := recipemanager.UserID(query[2])
			recipeID := recipemanager.RecipeID(query[3])
			ingredients := parseIngredients(query[4])
			instructions := query[5]
			_ = recipeManager.AddRecipe(int64(timestamp), userID, recipeID, ingredients, instructions)

		case "RATE_RECIPE":
			timestamp, _ := strconv.Atoi(query[1])
			userID := recipemanager.UserID(query[2])
			recipeID := recipemanager.RecipeID(query[3])
			score, _ := strconv.Atoi(query[4])
			_ = recipeManager.RateRecipe(int64(timestamp), userID, recipeID, int64(score))

		case "FAVORITE_RECIPE":
			timestamp, _ := strconv.Atoi(query[1])
			userID := recipemanager.UserID(query[2])
			recipeID := recipemanager.RecipeID(query[3])
			_ = recipeManager.FavoriteRecipe(int64(timestamp), userID, recipeID)

		case "GET_TOP_RECIPES":
			timestamp, _ := strconv.Atoi(query[1])
			n, _ := strconv.Atoi(query[2])
			_ = recipeManager.GetTopRecipes(int64(timestamp), n)

		case "SEARCH_BY_INGREDIENT":
			timestamp, _ := strconv.Atoi(query[1])
			ingredient := query[2]
			_ = recipeManager.SearchByIngredient(int64(timestamp), ingredient)

		case "GET_USER_FAVORITES":
			timestamp, _ := strconv.Atoi(query[1])
			userID := recipemanager.UserID(query[2])
			_ = recipeManager.GetUserFavorites(int64(timestamp), userID)

		case "GET_RECOMMENDATIONS":
			timestamp, _ := strconv.Atoi(query[1])
			userID := recipemanager.UserID(query[2])
			n, _ := strconv.Atoi(query[3])
			_ = recipeManager.GetRecommendations(int64(timestamp), userID, n)

		case "GET_LEADERBOARD":
			timestamp, _ := strconv.Atoi(query[1])
			n, _ := strconv.Atoi(query[2])
			_ = recipeManager.GetLeaderboard(int64(timestamp), n)
		}
	}
}

func printRecipeManager(rm recipemanager.RecipeManager) {
	fmt.Println("RECIPE MANAGER STATE")

	impl, ok := rm.(*recipemanager.RecipeManagerImpl)
	if !ok {
		fmt.Println("Cannot access internal state for printing")
		return
	}

	fmt.Println("USERS:")
	for userID, user := range impl.Users {
		fmt.Printf("  User: %s (Created at: %d)\n", userID, user.Timestamp)
		if len(user.FavoriteRecipes) > 0 {
			fmt.Printf("    Favorites: ")
			for i, recipe := range user.FavoriteRecipes {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Print(recipe.ID)
			}
			fmt.Println()
		}
	}

	fmt.Println("\nRECIPES:")
	for recipeID, recipe := range impl.Recipes {
		if recipe != nil {
			fmt.Printf("  Recipe: %s\n", recipeID)
			fmt.Printf("    Created by: %s (at timestamp: %d)\n", recipe.UserID, recipe.Timestamp)
			fmt.Printf("    Ingredients: %v\n", recipe.Ingredients)
			fmt.Printf("    Instructions: %s\n", recipe.Instructions)
			fmt.Printf("    Average Rating: %.2f (%d ratings)\n", recipe.AverageRating, len(recipe.Ratings))
			if len(recipe.Ratings) > 0 {
				fmt.Printf("    Ratings: ")
				for i, rating := range recipe.Ratings {
					if i > 0 {
						fmt.Print(", ")
					}
					fmt.Printf("%s:%d", rating.UserID, rating.Rating)
				}
				fmt.Println()
			}
			fmt.Printf("    Favorites: %d users\n", len(recipe.Favorites))
		}
	}

	fmt.Println("\nSTATISTICS:")
	fmt.Printf("  Total Users: %d\n", len(impl.Users))
	fmt.Printf("  Total Recipes: %d\n", len(impl.Recipes))

	totalRatings := 0
	totalFavorites := 0
	for _, recipe := range impl.Recipes {
		if recipe != nil {
			totalRatings += len(recipe.Ratings)
			totalFavorites += len(recipe.Favorites)
		}
	}
	fmt.Printf("  Total Ratings: %d\n", totalRatings)
	fmt.Printf("  Total Favorites: %d\n", totalFavorites)
}

func parseIngredients(ingredientStr string) []string {
	ingredients := []string{}
	var current string
	for _, ch := range ingredientStr {
		if ch == ',' {
			if current != "" {
				ingredients = append(ingredients, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		ingredients = append(ingredients, current)
	}
	return ingredients
}
