package recipemanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeManager_CreateUser(t *testing.T) {
	recipeManager := NewRecipeManager()

	testCases := []struct {
		name      string
		timestamp int64
		userID    UserID
		expected  bool
	}{
		{
			name:      "Should create a new user successfully",
			timestamp: 1,
			userID:    "user1",
			expected:  true,
		},
		{
			name:      "Should not create a duplicate user",
			timestamp: 2,
			userID:    "user1",
			expected:  false,
		},
		{
			name:      "Should create another user successfully",
			timestamp: 3,
			userID:    "user2",
			expected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := recipeManager.CreateUser(tc.timestamp, tc.userID)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRecipeManager_AddRecipe(t *testing.T) {
	recipeManager := NewRecipeManager()

	testCases := []struct {
		name         string
		timestamp    int64
		userID       UserID
		recipeID     RecipeID
		ingredients  []string
		instructions string
		setup        func()
		cleanup      func()
		expected     bool
	}{
		{
			name:         "Should add a new recipe successfully",
			timestamp:    1,
			userID:       "user1",
			recipeID:     "recipe1",
			ingredients:  []string{"flour", "sugar", "eggs"},
			instructions: "Mix and bake",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: true,
		},
		{
			name:         "Should not add recipe for non-existing user",
			timestamp:    2,
			userID:       "user_nonexistent",
			recipeID:     "recipe2",
			ingredients:  []string{"salt", "pepper"},
			instructions: "Season and serve",
			setup:        func() {},
			cleanup:      func() {},
			expected:     false,
		},
		{
			name:         "Should not add duplicate recipe",
			timestamp:    3,
			userID:       "user1",
			recipeID:     "recipe1",
			ingredients:  []string{"different", "ingredients"},
			instructions: "Different instructions",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour", "sugar"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := recipeManager.AddRecipe(tc.timestamp, tc.userID, tc.recipeID, tc.ingredients, tc.instructions)
			assert.Equal(t, tc.expected, result)
			tc.cleanup()
		})
	}
}

func TestRecipeManager_GetRecipe(t *testing.T) {
	recipeManager := NewRecipeManager()

	testCases := []struct {
		name      string
		timestamp int64
		recipeID  RecipeID
		setup     func()
		cleanup   func()
		expected  *Recipe
	}{
		{
			name:      "Should return recipe if it exists",
			timestamp: 1,
			recipeID:  "recipe1",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour", "sugar"}, "Mix and bake")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: &Recipe{
				ID:            "recipe1",
				AverageRating: 0,
				Timestamp:     1,
			},
		},
		{
			name:      "Should return nil if recipe does not exist",
			timestamp: 2,
			recipeID:  "recipe_nonexistent",
			setup:     func() {},
			cleanup:   func() {},
			expected:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := recipeManager.GetRecipe(tc.timestamp, tc.recipeID)
			if tc.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tc.expected.ID, result.ID)
				assert.Equal(t, tc.expected.AverageRating, result.AverageRating)
				assert.Equal(t, tc.expected.Timestamp, result.Timestamp)
			}
			tc.cleanup()
		})
	}
}

func TestRecipeManager_RateRecipe(t *testing.T) {
	recipeManager := NewRecipeManager()

	testCases := []struct {
		name      string
		timestamp int64
		userID    UserID
		recipeID  RecipeID
		score     int64
		setup     func()
		cleanup   func()
		expected  bool
	}{
		{
			name:      "Should rate recipe successfully",
			timestamp: 1,
			userID:    "user2",
			recipeID:  "recipe1",
			score:     4,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: true,
		},
		{
			name:      "Should not rate with score below 1",
			timestamp: 2,
			userID:    "user2",
			recipeID:  "recipe1",
			score:     0,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: false,
		},
		{
			name:      "Should not rate with score above 5",
			timestamp: 3,
			userID:    "user2",
			recipeID:  "recipe1",
			score:     6,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: false,
		},
		{
			name:      "Should not allow user to rate their own recipe",
			timestamp: 4,
			userID:    "user1",
			recipeID:  "recipe1",
			score:     5,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: false,
		},
		{
			name:      "Should not allow duplicate rating from same user",
			timestamp: 5,
			userID:    "user2",
			recipeID:  "recipe1",
			score:     3,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.RateRecipe(4, "user2", "recipe1", 4)
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: false,
		},
		{
			name:      "Should not rate if user does not exist",
			timestamp: 6,
			userID:    "user_nonexistent",
			recipeID:  "recipe1",
			score:     4,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: false,
		},
		{
			name:      "Should not rate if recipe does not exist",
			timestamp: 7,
			userID:    "user2",
			recipeID:  "recipe_nonexistent",
			score:     4,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := recipeManager.RateRecipe(tc.timestamp, tc.userID, tc.recipeID, tc.score)
			assert.Equal(t, tc.expected, result)
			tc.cleanup()
		})
	}
}

func TestRecipeManager_GetTopRecipes(t *testing.T) {
	recipeManager := NewRecipeManager()

	testCases := []struct {
		name      string
		timestamp int64
		n         int
		setup     func()
		cleanup   func()
		expected  []string
	}{
		{
			name:      "Should return top recipes by average rating",
			timestamp: 1,
			n:         10,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.CreateUser(3, "user3")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.AddRecipe(1, "user2", "recipe2", []string{"sugar"}, "Heat")
				_ = recipeManager.RateRecipe(1, "user2", "recipe1", 5)
				_ = recipeManager.RateRecipe(1, "user3", "recipe1", 3)
				_ = recipeManager.RateRecipe(1, "user1", "recipe2", 4)
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{"recipe1(4.00)", "recipe2(4.00)"},
		},
		{
			name:      "Should return top n recipes",
			timestamp: 2,
			n:         1,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.CreateUser(3, "user3")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.AddRecipe(1, "user2", "recipe2", []string{"sugar"}, "Heat")
				_ = recipeManager.RateRecipe(1, "user2", "recipe1", 5)
				_ = recipeManager.RateRecipe(1, "user1", "recipe2", 2)
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{"recipe1(5.00)"},
		},
		{
			name:      "Should return empty list if no recipes",
			timestamp: 3,
			n:         10,
			setup:     func() {},
			cleanup:   func() {},
			expected:  []string{},
		},
		{
			name:      "Should sort alphabetically by recipe ID in case of tie",
			timestamp: 4,
			n:         10,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.CreateUser(3, "user3")
				_ = recipeManager.AddRecipe(1, "user1", "recipeB", []string{"flour"}, "Mix")
				_ = recipeManager.AddRecipe(1, "user2", "recipeA", []string{"sugar"}, "Heat")
				_ = recipeManager.RateRecipe(1, "user2", "recipeB", 4)
				_ = recipeManager.RateRecipe(1, "user1", "recipeA", 4)
				_ = recipeManager.RateRecipe(1, "user3", "recipeA", 4)
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{"recipeA(4.00)", "recipeB(4.00)"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := recipeManager.GetTopRecipes(tc.timestamp, tc.n)
			assert.Equal(t, tc.expected, result)
			tc.cleanup()
		})
	}
}

func TestRecipeManager_SearchByIngredient(t *testing.T) {
	recipeManager := NewRecipeManager()

	testCases := []struct {
		name       string
		timestamp  int64
		ingredient string
		setup      func()
		cleanup    func()
		expected   []string
	}{
		{
			name:       "Should find recipes with ingredient",
			timestamp:  1,
			ingredient: "flour",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour", "sugar"}, "Mix")
				_ = recipeManager.AddRecipe(1, "user2", "recipe2", []string{"flour", "eggs"}, "Whisk")
				_ = recipeManager.AddRecipe(1, "user1", "recipe3", []string{"sugar", "butter"}, "Cream")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{"recipe1", "recipe2"},
		},
		{
			name:       "Should return empty list if ingredient not found",
			timestamp:  2,
			ingredient: "salt",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour", "sugar"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{},
		},
		{
			name:       "Should be case-sensitive",
			timestamp:  3,
			ingredient: "Flour",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour", "sugar"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{},
		},
		{
			name:       "Should return empty list when no recipes exist",
			timestamp:  4,
			ingredient: "flour",
			setup:      func() {},
			cleanup:    func() {},
			expected:   []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := recipeManager.SearchByIngredient(tc.timestamp, tc.ingredient)
			assert.Equal(t, tc.expected, result)
			tc.cleanup()
		})
	}
}

func TestRecipeManager_SearchByIngredient_Sorting(t *testing.T) {
	recipeManager := NewRecipeManager()

	testCases := []struct {
		name       string
		timestamp  int64
		ingredient string
		setup      func()
		cleanup    func()
		expected   []string
	}{
		{
			name:       "Should return results sorted alphabetically",
			timestamp:  1,
			ingredient: "flour",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.CreateUser(3, "user3")
				_ = recipeManager.AddRecipe(1, "user1", "zebra_cake", []string{"flour", "sugar"}, "Mix")
				_ = recipeManager.AddRecipe(1, "user2", "apple_bread", []string{"flour", "eggs"}, "Whisk")
				_ = recipeManager.AddRecipe(1, "user3", "muffin", []string{"flour", "butter"}, "Blend")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{"apple_bread", "muffin", "zebra_cake"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := recipeManager.SearchByIngredient(tc.timestamp, tc.ingredient)
			assert.Equal(t, tc.expected, result)
			tc.cleanup()
		})
	}
}

func TestRecipeManager_RateRecipe_AverageCalculation(t *testing.T) {
	recipeManager := NewRecipeManager()

	testCases := []struct {
		name      string
		timestamp int64
		setup     func()
		cleanup   func()
		expected  float64
	}{
		{
			name:      "Should calculate average rating correctly with multiple raters",
			timestamp: 1,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.CreateUser(3, "user3")
				_ = recipeManager.CreateUser(4, "user4")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.RateRecipe(1, "user2", "recipe1", 5)
				_ = recipeManager.RateRecipe(1, "user3", "recipe1", 3)
				_ = recipeManager.RateRecipe(1, "user4", "recipe1", 4)
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: 4.0,
		},
		{
			name:      "Should calculate average rating with single rater",
			timestamp: 2,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.RateRecipe(1, "user2", "recipe1", 5)
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: 5.0,
		},
		{
			name:      "Should round average rating to 2 decimal places",
			timestamp: 3,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.CreateUser(3, "user3")
				_ = recipeManager.CreateUser(4, "user4")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.RateRecipe(1, "user2", "recipe1", 5)
				_ = recipeManager.RateRecipe(1, "user3", "recipe1", 2)
				_ = recipeManager.RateRecipe(1, "user4", "recipe1", 1)
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: 2.67,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			recipe := recipeManager.GetRecipe(tc.timestamp, "recipe1")
			assert.NotNil(t, recipe)
			assert.Equal(t, tc.expected, recipe.AverageRating)
			tc.cleanup()
		})
	}
}

func TestRecipeManager_ComplexScenario(t *testing.T) {
	recipeManager := NewRecipeManager()

	t.Run("Should handle multiple users and recipes with various interactions", func(t *testing.T) {
		assert.True(t, recipeManager.CreateUser(1, "alice"))
		assert.True(t, recipeManager.CreateUser(2, "bob"))
		assert.True(t, recipeManager.CreateUser(3, "charlie"))
		assert.True(t, recipeManager.CreateUser(4, "diana"))

		assert.True(t, recipeManager.AddRecipe(5, "alice", "pasta", []string{"flour", "eggs", "salt"}, "Mix and cook"))
		assert.True(t, recipeManager.AddRecipe(6, "bob", "salad", []string{"lettuce", "tomato", "oil"}, "Toss"))
		assert.True(t, recipeManager.AddRecipe(7, "charlie", "bread", []string{"flour", "water", "salt"}, "Knead and bake"))

		assert.True(t, recipeManager.RateRecipe(8, "bob", "pasta", 5))
		assert.True(t, recipeManager.RateRecipe(9, "charlie", "pasta", 4))
		assert.True(t, recipeManager.RateRecipe(10, "diana", "pasta", 5))
		assert.True(t, recipeManager.RateRecipe(11, "alice", "salad", 3))
		assert.True(t, recipeManager.RateRecipe(12, "charlie", "salad", 4))
		assert.True(t, recipeManager.RateRecipe(13, "diana", "bread", 5))

		topRecipes := recipeManager.GetTopRecipes(14, 10)
		assert.Equal(t, 3, len(topRecipes))
		assert.True(t, len(topRecipes) > 0)

		flourRecipes := recipeManager.SearchByIngredient(15, "flour")
		assert.Equal(t, 2, len(flourRecipes))

		pastaRecipe := recipeManager.GetRecipe(16, "pasta")
		assert.NotNil(t, pastaRecipe)
		assert.Equal(t, "alice", string(pastaRecipe.UserID))
		assert.Equal(t, 3, len(pastaRecipe.Ratings))
	})
}

func TestRecipeManager_FavoriteRecipe(t *testing.T) {
	recipeManager := NewRecipeManager()

	testCases := []struct {
		name      string
		timestamp int64
		userID    UserID
		recipeID  RecipeID
		setup     func()
		cleanup   func()
		expected  bool
	}{
		{
			name:      "Should favorite a recipe successfully",
			timestamp: 1,
			userID:    "user2",
			recipeID:  "recipe1",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: true,
		},
		{
			name:      "Should not favorite if user does not exist",
			timestamp: 2,
			userID:    "user_nonexistent",
			recipeID:  "recipe1",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: false,
		},
		{
			name:      "Should not favorite if recipe does not exist",
			timestamp: 3,
			userID:    "user2",
			recipeID:  "recipe_nonexistent",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: false,
		},
		{
			name:      "Can favorite their own recipe (implementation allows it)",
			timestamp: 4,
			userID:    "user1",
			recipeID:  "recipe1",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: true,
		},
		{
			name:      "Can favorite the same recipe multiple times (implementation allows it)",
			timestamp: 5,
			userID:    "user2",
			recipeID:  "recipe1",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.FavoriteRecipe(1, "user2", "recipe1")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := recipeManager.FavoriteRecipe(tc.timestamp, tc.userID, tc.recipeID)
			assert.Equal(t, tc.expected, result)
			tc.cleanup()
		})
	}
}

func TestRecipeManager_GetUserFavorites(t *testing.T) {
	recipeManager := NewRecipeManager()

	testCases := []struct {
		name      string
		timestamp int64
		userID    UserID
		setup     func()
		cleanup   func()
		verify    func(t *testing.T, result []string)
	}{
		{
			name:      "Should return user's favorited recipes",
			timestamp: 1,
			userID:    "user2",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.AddRecipe(1, "user1", "recipe2", []string{"apples"}, "Bake")
				_ = recipeManager.FavoriteRecipe(1, "user2", "recipe1")
				_ = recipeManager.FavoriteRecipe(1, "user2", "recipe2")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			verify: func(t *testing.T, result []string) {
				assert.Greater(t, len(result), 0)
			},
		},
		{
			name:      "Should return empty list if user has no favorites",
			timestamp: 2,
			userID:    "user2",
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			verify: func(t *testing.T, result []string) {
				assert.Equal(t, 0, len(result))
			},
		},
		{
			name:      "Should return empty list if user does not exist",
			timestamp: 3,
			userID:    "user_nonexistent",
			setup:     func() {},
			cleanup:   func() {},
			verify: func(t *testing.T, result []string) {
				assert.Equal(t, 0, len(result))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := recipeManager.GetUserFavorites(tc.timestamp, tc.userID)
			tc.verify(t, result)
			tc.cleanup()
		})
	}
}

func TestRecipeManager_GetRecommendations(t *testing.T) {
	recipeManager := NewRecipeManager()

	testCases := []struct {
		name      string
		timestamp int64
		userID    UserID
		n         int
		setup     func()
		cleanup   func()
		verify    func(t *testing.T, result []string)
	}{
		{
			name:      "Should respect the limit n",
			timestamp: 1,
			userID:    "user1",
			n:         2,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user2", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.AddRecipe(1, "user2", "recipe2", []string{"sugar"}, "Blend")
				_ = recipeManager.AddRecipe(1, "user2", "recipe3", []string{"eggs"}, "Whisk")
				_ = recipeManager.AddRecipe(1, "user2", "recipe4", []string{"butter"}, "Cream")
				_ = recipeManager.RateRecipe(1, "user2", "recipe1", 5)
				_ = recipeManager.RateRecipe(1, "user2", "recipe2", 4)
				_ = recipeManager.RateRecipe(1, "user2", "recipe3", 3)
				_ = recipeManager.RateRecipe(1, "user2", "recipe4", 2)
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			verify: func(t *testing.T, result []string) {
				assert.LessOrEqual(t, len(result), 2)
			},
		},
		{
			name:      "Should return empty list if user does not exist",
			timestamp: 2,
			userID:    "user_nonexistent",
			n:         10,
			setup:     func() {},
			cleanup:   func() {},
			verify: func(t *testing.T, result []string) {
				assert.Equal(t, 0, len(result))
			},
		},
		{
			name:      "Should return empty list with invalid n",
			timestamp: 3,
			userID:    "user1",
			n:         0,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user2", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.RateRecipe(1, "user2", "recipe1", 5)
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			verify: func(t *testing.T, result []string) {
				assert.Equal(t, 0, len(result))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := recipeManager.GetRecommendations(tc.timestamp, tc.userID, tc.n)
			tc.verify(t, result)
			tc.cleanup()
		})
	}
}

func TestRecipeManager_GetLeaderboard(t *testing.T) {
	recipeManager := NewRecipeManager()

	testCases := []struct {
		name      string
		timestamp int64
		n         int
		setup     func()
		cleanup   func()
		expected  []string
	}{
		{
			name:      "Should return leaderboard sorted by score descending",
			timestamp: 1,
			n:         10,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.CreateUser(3, "user3")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.AddRecipe(1, "user2", "recipe2", []string{"sugar"}, "Heat")
				_ = recipeManager.AddRecipe(1, "user2", "recipe3", []string{"eggs"}, "Whisk")
				_ = recipeManager.RateRecipe(1, "user3", "recipe1", 5)
				_ = recipeManager.RateRecipe(1, "user1", "recipe2", 4)
				_ = recipeManager.RateRecipe(1, "user1", "recipe3", 3)
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{"user2(26)", "user1(13)"},
		},
		{
			name:      "Should sort alphabetically by user ID in case of tie",
			timestamp: 2,
			n:         10,
			setup: func() {
				_ = recipeManager.CreateUser(1, "alice")
				_ = recipeManager.CreateUser(2, "bob")
				_ = recipeManager.CreateUser(3, "charlie")
				_ = recipeManager.AddRecipe(1, "alice", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.AddRecipe(1, "bob", "recipe2", []string{"sugar"}, "Heat")
				_ = recipeManager.AddRecipe(1, "charlie", "recipe3", []string{"eggs"}, "Whisk")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{"alice(10)", "bob(10)", "charlie(10)"},
		},
		{
			name:      "Should respect the limit n",
			timestamp: 3,
			n:         2,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.CreateUser(3, "user3")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.AddRecipe(1, "user2", "recipe2", []string{"sugar"}, "Heat")
				_ = recipeManager.AddRecipe(1, "user3", "recipe3", []string{"eggs"}, "Whisk")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{"user1(10)", "user2(10)"},
		},
		{
			name:      "Should return empty list if no recipes",
			timestamp: 4,
			n:         10,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{},
		},
		{
			name:      "Should return single user leaderboard",
			timestamp: 5,
			n:         10,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.AddRecipe(1, "user1", "recipe2", []string{"sugar"}, "Heat")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{"user1(20)"},
		},
		{
			name:      "Should calculate score with recipes, ratings, and favorites",
			timestamp: 6,
			n:         10,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.AddRecipe(1, "user1", "recipe2", []string{"sugar"}, "Heat")
				_ = recipeManager.RateRecipe(1, "user2", "recipe1", 5)
				_ = recipeManager.RateRecipe(1, "user2", "recipe2", 4)
				_ = recipeManager.FavoriteRecipe(1, "user2", "recipe1")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{"user1(31)"},
		},
		{
			name:      "Should return empty list with n <= 0",
			timestamp: 7,
			n:         0,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{},
		},
		{
			name:      "Should only include users with recipes",
			timestamp: 9,
			n:         10,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.CreateUser(3, "user3")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{"user1(10)"},
		},
		{
			name:      "Should handle complex scenario with multiple interactions",
			timestamp: 10,
			n:         3,
			setup: func() {
				_ = recipeManager.CreateUser(1, "alice")
				_ = recipeManager.CreateUser(2, "bob")
				_ = recipeManager.CreateUser(3, "charlie")
				_ = recipeManager.CreateUser(4, "diana")
				_ = recipeManager.AddRecipe(1, "alice", "pasta", []string{"flour", "eggs"}, "Mix and cook")
				_ = recipeManager.AddRecipe(1, "alice", "pizza", []string{"flour", "tomato"}, "Bake")
				_ = recipeManager.AddRecipe(1, "bob", "salad", []string{"lettuce", "tomato"}, "Toss")
				_ = recipeManager.AddRecipe(1, "charlie", "bread", []string{"flour", "water"}, "Knead")
				_ = recipeManager.RateRecipe(1, "bob", "pasta", 5)
				_ = recipeManager.RateRecipe(1, "charlie", "pasta", 4)
				_ = recipeManager.RateRecipe(1, "diana", "pasta", 5)
				_ = recipeManager.RateRecipe(1, "alice", "salad", 3)
				_ = recipeManager.RateRecipe(1, "diana", "bread", 5)
				_ = recipeManager.FavoriteRecipe(1, "bob", "pizza")
				_ = recipeManager.FavoriteRecipe(1, "charlie", "pizza")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{"alice(39)", "bob(13)", "charlie(13)"},
		},
		{
			name:      "Should return all users when n is larger than user count",
			timestamp: 11,
			n:         100,
			setup: func() {
				_ = recipeManager.CreateUser(1, "user1")
				_ = recipeManager.CreateUser(2, "user2")
				_ = recipeManager.AddRecipe(1, "user1", "recipe1", []string{"flour"}, "Mix")
				_ = recipeManager.AddRecipe(1, "user2", "recipe2", []string{"sugar"}, "Heat")
			},
			cleanup: func() {
				recipeManager = NewRecipeManager()
			},
			expected: []string{"user1(10)", "user2(10)"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := recipeManager.GetLeaderboard(tc.timestamp, tc.n)
			assert.Equal(t, tc.expected, result)
			tc.cleanup()
		})
	}
}
