##
## EPITECH PROJECT, 2026
## ~/epitech/hub
## File description:
## Makefile
##

build:
	@mkdir -p bin
	@echo "🏗️  Build PLF-cli" && go build -o bin/plf main.go
	@chmod +x bin/plf

run: build
	@echo "🔐 Run PLF-cli"
	@./bin/plf

clean:
	@echo "🧹 Clean PLF"
	@rm -rf bin

re: clean build

.PHONY: build clean re