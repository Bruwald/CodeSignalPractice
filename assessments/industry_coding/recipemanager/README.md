# Recipe Management System

Your task is to implement a simplified **Recipe Management System**. All operations that should be supported are listed below.

Solving this task consists of several levels. In real test, a subsequent level is opened when the current level tests are correct. You always have access to the data for the current and all previous levels.

## Requirements

Plan your design according to the level specifications below. Your code should be **clean, concise, and extensible** for future levels.

*   **Level 1:** Basic user creation, recipe addition, and recipe retrieval.
*   **Level 2:** Rating system and discovery features.
*   **Level 3:** Favorites and personalized recommendations.
*   **Level 4:** User merging and historical analysis.

**Note:** All operations will have a `timestamp` parameter — a stringified timestamp in milliseconds. It is guaranteed that all timestamps are unique and are in a range from `1` to `10^9`. Operations will be given in order of strictly increasing timestamps.

---

## Level 1: Basic Recipe Management

Initially, the system does not contain any users or recipes.

*   `create_user(self, timestamp: int, user_id: str) -> bool` — Creates a new user. Returns `True` if the user was successfully created, `False` if the user already exists.
*   `add_recipe(self, timestamp: int, user_id: str, recipe_id: str, ingredients: list[str], instructions: str) -> bool` — Adds a new recipe. Returns `True` if successful, `False` if the user doesn't exist or the recipe already exists.
*   `get_recipe(self, timestamp: int, recipe_id: str) -> dict | None` — Returns recipe information as a dictionary with the following keys:  
  `{"recipe_id", "user_id", "ingredients", "instructions", "average_rating"}`.  
  At this level, `average_rating` should always be `0`. Returns `None` if the recipe doesn't exist.

### Example

| Queries | Explanations |
| --- | --- |
| create_user(1, "alice") | returns True |
| create_user(2, "alice") | returns False |
| add_recipe(3, "alice", "pasta", ["pasta", "tomato", "basil"], "Boil the pasta...") | returns True |
| get_recipe(4, "pasta") | returns `{"recipe_id": "pasta", "user_id": "alice", "ingredients": ["pasta", "tomato", "basil"], "instructions": "Boil the pasta...", "average_rating": 0}` |

---

## Level 2: Rating System & Discovery

*   `rate_recipe(self, timestamp: int, user_id: str, recipe_id: str, score: int) -> bool` — Allows a user to rate a recipe (score must be between 1 and 5 inclusive).  
  Returns `False` if user doesn't exist, recipe doesn't exist, user is trying to rate their own recipe, or they already rated it.
*   `get_top_recipes(self, timestamp: int, n: int) -> list[str]` — Returns the top `n` recipes by average rating (descending). In case of ties, sort by `recipe_id` alphabetically ascending.  
  Format: `["recipe_id1(4.5)", "recipe_id2(4.0)", ...]`
*   `search_by_ingredient(self, timestamp: int, ingredient: str) -> list[str]` — Returns list of `recipe_id`s that contain the given ingredient (case-sensitive), sorted alphabetically.

**Note:** Average rating should be calculated as the arithmetic mean of all scores.

---

## Level 3: Favorites & Recommendations

*   `favorite_recipe(self, timestamp: int, user_id: str, recipe_id: str) -> bool` — Adds a recipe to the user's favorites. Returns `True` if successful.
*   `get_user_favorites(self, timestamp: int, user_id: str) -> list[str]` — Returns list of recipe ids favorited by the user, sorted alphabetically.
*   `get_recommendations(self, timestamp: int, user_id: str, n: int) -> list[str]` — Returns up to `n` recommended recipes for the user.  
  **Recommendation Logic**: Highest-rated recipes that the user has **not rated** and **not favorited**, preferably from users with similar taste (shared at least one recipe with rating difference ≤ 1).

---

## Level 4: Leaderboard

*   `get_leaderboard(self, timestamp: int, n: int) -> list[str]` — Returns top `n` users by engagement score:  
  `score = (recipes_created * 10) + (total_ratings_received * 3) + (total_favorites_received * 5)`  
  Format: `["user1(245)", "user2(180)", ...]`

---

## General Notes

- All operations must respect the strictly increasing timestamp order.
- Design your data structures carefully from the beginning to minimize refactoring.
- Average ratings should be displayed with one decimal place when required.
- Think about efficient ways to calculate averages, search ingredients, and generate recommendations.

**Test Commands:**
- Level 1: `pytest test_recipe_system.py::TestLevel1 -v`
- Level 2: `pytest test_recipe_system.py::TestLevel2 -v`
- Level 3: `pytest test_recipe_system.py::TestLevel3 -v`
- Level 4: `pytest test_recipe_system.py::TestLevel4 -v`

**Execution Constraints:**
- **[execution time limit]** 3 seconds
- **[memory limit]** 1 GB