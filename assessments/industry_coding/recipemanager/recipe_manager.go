package recipemanager

import (
	"math"
	"slices"
	"sort"
	"strconv"
)

const (
	LeaderBoardRecipesMultiplicationFactor   int = 10
	LeaderBoardRatingsMultiplicationFactor   int = 3
	LeaderBoardFavoritesMultiplicationFactor int = 5
)

type (
	RecipeID string
	UserID   string

	RecipeManager interface {
		// Creates a new user.
		// Returns `True` if the user was successfully created, `False` if the user already exists.
		CreateUser(timestamp int64, userID UserID) bool
		// Adds a new recipe.
		// Returns `True` if successful, `False` if the user doesn't exist or the recipe already exists.
		AddRecipe(timestamp int64, userID UserID, recipeID RecipeID, ingredients []string, instructions string) bool
		// Returns recipe information.
		GetRecipe(timestamp int64, recipeID RecipeID) *Recipe
		// Allows a user to rate a recipe (score must be between 1 and 5 inclusive).
		// Returns `False` if user doesn't exist, recipe doesn't exist, user is trying to rate their own recipe, or they already rated it.
		RateRecipe(timestamp int64, userID UserID, recipeID RecipeID, score int64) bool
		// Returns the top `n` recipes by average rating (descending).
		// In case of ties, sort by `recipe_id` alphabetically ascending.
		// Format: `["recipe_id1(4.5)", "recipe_id2(4.0)", ...]`
		GetTopRecipes(timestamp int64, n int) []string
		// Returns list of `recipe_id`s that contain the given ingredient (case-sensitive), sorted alphabetically.
		SearchByIngredient(timestamp int64, ingredient string) []string
		// Adds a recipe to the user's favorites. Returns `True` if successful.
		FavoriteRecipe(timestamp int64, userID UserID, recipeID RecipeID) bool
		// Returns list of recipe ids favorited by the user, sorted alphabetically.
		GetUserFavorites(timestamp int64, userID UserID) []string
		// Returns up to `n` recommended recipes for the user.
		// Highest-rated recipes that the user has **not rated** and **not favorited**, preferably from users
		// with similar taste (shared at least one recipe with rating difference ≤ 1).
		GetRecommendations(timestamp int64, userID UserID, n int) []string
		// Returns top `n` users by engagement score:
		// `score = (recipes_created * 10) + (total_ratings_received * 3) + (total_favorites_received * 5)`
		// Format: `["user1(245)", "user2(180)", ...]`
		GetLeaderboard(timestamp int64, n int) []string
	}

	User struct {
		ID              UserID
		Timestamp       int64
		FavoriteRecipes []*Recipe
	}

	Rating struct {
		UserID UserID
		Rating int64
	}

	Recipe struct {
		ID            RecipeID
		UserID        UserID
		Ingredients   []string
		Instructions  string
		AverageRating float64
		Timestamp     int64
		Ratings       []Rating
		Favorites     []*User
	}

	LeaderBoardStatus struct {
		RecipesCreated         int
		TotalRatingsReceived   int
		TotalFavoritesReceived int
	}

	RecipeManagerImpl struct {
		Users   map[UserID]*User
		Recipes map[RecipeID]*Recipe
	}
)

func NewRecipeManager() RecipeManager {
	return &RecipeManagerImpl{
		Users:   map[UserID]*User{},
		Recipes: map[RecipeID]*Recipe{},
	}
}

func (r *RecipeManagerImpl) CreateUser(timestamp int64, userID UserID) bool {
	if _, userExists := r.Users[userID]; userExists {
		return false
	}
	r.Users[userID] = &User{
		ID:        userID,
		Timestamp: timestamp,
	}
	return true
}

func (r *RecipeManagerImpl) AddRecipe(timestamp int64, userID UserID, recipeID RecipeID, ingredients []string, instructions string) bool {
	user, userExists := r.Users[userID]
	_, recipeExists := r.Recipes[recipeID]
	if !userExists || recipeExists {
		return false
	}
	recipe := Recipe{
		ID:            recipeID,
		UserID:        user.ID,
		Ingredients:   ingredients,
		Instructions:  instructions,
		AverageRating: float64(0),
		Timestamp:     timestamp,
	}
	r.Recipes[recipeID] = &recipe
	return true
}

func (r *RecipeManagerImpl) GetRecipe(timestamp int64, recipeID RecipeID) *Recipe {
	return r.Recipes[recipeID]
}

func (r *RecipeManagerImpl) RateRecipe(timestamp int64, userID UserID, recipeID RecipeID, score int64) bool {
	user, userExists := r.Users[userID]
	recipe, recipeExists := r.Recipes[recipeID]
	if !userExists || !recipeExists || score < 1 || score > 5 || recipe.UserID == userID {
		return false
	}

	for _, rating := range recipe.Ratings {
		if rating.UserID == userID {
			return false
		}
	}

	recipe.Ratings = append(recipe.Ratings, Rating{
		UserID: user.ID,
		Rating: score,
	})

	sumRatings := float64(0)
	for _, rating := range recipe.Ratings {
		sumRatings += float64(rating.Rating)
	}
	recipe.AverageRating = math.Round(sumRatings/float64(len(recipe.Ratings))*100) / 100

	return true
}

func (r *RecipeManagerImpl) GetTopRecipes(timestamp int64, n int) []string {
	if n <= 0 {
		return []string{}
	}

	allRecipes := []Recipe{}
	for _, recipe := range r.Recipes {
		if recipe != nil {
			allRecipes = append(allRecipes, *recipe)
		}
	}

	sort.Slice(allRecipes, func(i, j int) bool {
		if allRecipes[i].AverageRating != allRecipes[j].AverageRating {
			return allRecipes[i].AverageRating > allRecipes[j].AverageRating
		}
		return allRecipes[i].ID < allRecipes[j].ID
	})

	if len(allRecipes) > n {
		allRecipes = allRecipes[:n]
	}

	topRecipes := []string{}
	for _, recipe := range allRecipes {
		topRecipes = append(topRecipes, string(recipe.ID)+"("+strconv.FormatFloat(recipe.AverageRating, 'f', 2, 64)+")")
	}

	return topRecipes
}

func (r *RecipeManagerImpl) SearchByIngredient(timestamp int64, ingredient string) []string {
	recipesWithIngredient := []string{}
	for recipeID, recipe := range r.Recipes {
		if recipe != nil && slices.Contains(recipe.Ingredients, ingredient) {
			recipesWithIngredient = append(recipesWithIngredient, string(recipeID))
		}
	}
	sort.Strings(recipesWithIngredient)
	return recipesWithIngredient
}

func (r *RecipeManagerImpl) FavoriteRecipe(timestamp int64, userID UserID, recipeID RecipeID) bool {
	user, userExists := r.Users[userID]
	recipe, recipeExists := r.Recipes[recipeID]
	if !userExists || !recipeExists {
		return false
	}

	user.FavoriteRecipes = append(user.FavoriteRecipes, recipe)
	recipe.Favorites = append(recipe.Favorites, user)

	return true
}

func (r *RecipeManagerImpl) GetUserFavorites(timestamp int64, userID UserID) []string {
	user, userExists := r.Users[userID]
	if !userExists {
		return []string{}
	}

	userFavorites := make([]string, len(user.FavoriteRecipes))
	for _, recipe := range user.FavoriteRecipes {
		userFavorites = append(userFavorites, string(recipe.ID))
	}

	sort.Slice(userFavorites, func(i, j int) bool {
		return userFavorites[i] < userFavorites[j]
	})

	return userFavorites
}

func (r *RecipeManagerImpl) GetRecommendations(timestamp int64, userID UserID, n int) []string {
	_, userExists := r.Users[userID]
	if !userExists || n <= 0 {
		return []string{}
	}

	recipesToRecommend := []Recipe{}
	for _, recipe := range r.Recipes {
		if recipe != nil && (recipe.UserID == userID || len(recipe.Ratings) == 0) {
			continue
		}
		for _, rating := range recipe.Ratings {
			if rating.UserID == userID {
				continue
			}
		}
		recipesToRecommend = append(recipesToRecommend, *recipe)
	}

	sort.Slice(recipesToRecommend, func(i, j int) bool {
		return recipesToRecommend[i].AverageRating > recipesToRecommend[j].AverageRating
	})

	if len(recipesToRecommend) > n {
		recipesToRecommend = recipesToRecommend[:n]
	}

	recipeIDsToRecommend := make([]string, len(recipesToRecommend))
	for _, recipeToRecommend := range recipesToRecommend {
		recipeIDsToRecommend = append(recipeIDsToRecommend, string(recipeToRecommend.ID))
	}

	return recipeIDsToRecommend
}

func (r *RecipeManagerImpl) GetLeaderboard(timestamp int64, n int) []string {
	if n <= 0 {
		return []string{}
	}

	leaderBoardStatusesByUser := map[UserID]*LeaderBoardStatus{}
	for _, recipe := range r.Recipes {
		if recipe != nil {
			if leaderBoardStatus, exists := leaderBoardStatusesByUser[recipe.UserID]; exists {
				leaderBoardStatus.RecipesCreated++
				leaderBoardStatus.TotalRatingsReceived += len(recipe.Ratings)
				leaderBoardStatus.TotalFavoritesReceived += len(recipe.Favorites)
				continue
			}

			leaderBoardStatusesByUser[recipe.UserID] = &LeaderBoardStatus{
				RecipesCreated:         1,
				TotalRatingsReceived:   len(recipe.Ratings),
				TotalFavoritesReceived: len(recipe.Favorites),
			}
		}
	}

	type scoreByUser struct {
		UserID UserID
		Score  int
	}

	scoresByUser := []scoreByUser{}
	for userID, status := range leaderBoardStatusesByUser {
		scoresByUser = append(scoresByUser, scoreByUser{
			UserID: userID,
			Score: status.RecipesCreated*LeaderBoardRecipesMultiplicationFactor +
				status.TotalRatingsReceived*LeaderBoardRatingsMultiplicationFactor +
				status.TotalFavoritesReceived*LeaderBoardFavoritesMultiplicationFactor,
		})
	}

	sort.Slice(scoresByUser, func(i, j int) bool {
		if scoresByUser[i].Score != scoresByUser[j].Score {
			return scoresByUser[i].Score > scoresByUser[j].Score
		}
		return scoresByUser[i].UserID < scoresByUser[j].UserID
	})

	if len(scoresByUser) > n {
		scoresByUser = scoresByUser[:n]
	}

	results := []string{}
	for _, scoreByUser := range scoresByUser {
		results = append(results, string(scoreByUser.UserID)+"("+strconv.Itoa(scoreByUser.Score)+")")
	}

	return results
}
