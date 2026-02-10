##
## EPITECH PROJECT, 2026
## ~/epitech/hub
## File description:
## Makefile
##

build:
	@mkdir -p bin
	@printf "📦 Fetching dependencies...  "
	@go mod tidy & PID=$$!; \
	while kill -0 $$PID 2>/dev/null; do \
		for s in / - \\\\ \|; do \
			printf "\b$$s"; \
			sleep 0.1; \
		done; \
	done; \
	wait $$PID; \
	printf "\r📦 Fetching dependencies...\n"
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