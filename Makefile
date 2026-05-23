run-banking-system:
	go run assessments/industry_coding/bankingsystem/main/main.go

test-banking-system:
	go test -v -race -cover assessments/industry_coding/bankingsystem/banking_system_test.go assessments/industry_coding/bankingsystem/banking_system.go
