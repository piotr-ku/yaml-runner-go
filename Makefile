test:
	@go test -cover ./app ./system
	@rm -f coverage.out

coverage:
	@~/.go/bin/courtney ./app ./system
	@go tool cover -html=coverage.out
	@rm -f coverage.out