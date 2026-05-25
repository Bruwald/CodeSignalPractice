run-banking-system:
	go run assessments/industry_coding/bankingsystem/main/main.go

test-banking-system:
	go test -v -race -cover assessments/industry_coding/bankingsystem/banking_system_test.go assessments/industry_coding/bankingsystem/banking_system.go

run-recipe-manager:
	go run assessments/industry_coding/recipemanager/main/main.go

test-recipe-manager:
	go test -v -race -cover assessments/industry_coding/recipemanager/recipe_manager_test.go assessments/industry_coding/recipemanager/recipe_manager.go
