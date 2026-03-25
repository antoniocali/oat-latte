.PHONY: build test cover vet fmt check clean run-notes run-kanban run-tasklist \
        docs-install docs-start docs-build docs-serve docs-clean help

# ANSI colours
BOLD   := \033[1m
RESET  := \033[0m
CYAN   := \033[36m
YELLOW := \033[33m
GREEN  := \033[32m
GRAY   := \033[90m

# ── Go ─────────────────────────────────────────────────────────────

build:       ## Build all packages
	go build ./...

test:        ## Run all tests
	go test ./...

cover:       ## Run tests and open HTML coverage report
	go test -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt -o coverage.html

vet:         ## Run go vet
	go vet ./...

fmt:         ## Format all source files
	go fmt ./...

check: fmt vet test  ## Run fmt, vet, and test

# ── Examples ───────────────────────────────────────────────────────

run-notes:   ## Run the Notes example app
	go run ./cmd/example/notes

run-kanban:  ## Run the Kanban example app
	go run ./cmd/example/kanban

run-tasklist: ## Run the Task List example app
	go run ./cmd/example/tasklist

# ── Cleanup ────────────────────────────────────────────────────────

clean:       ## Remove Go build artifacts
	rm -f coverage.txt coverage.html
	rm -rf bin/ dist/

# ── Docs (Docusaurus) ──────────────────────────────────────────────

docs-install:  ## Install docs dependencies (run once)
	npm install --prefix docs-site

docs-start:    ## Start docs dev server with live reload → http://localhost:3000
	npm run start --prefix docs-site

docs-build:    ## Build docs for production → docs-site/build/
	npm run build --prefix docs-site

docs-serve: docs-build  ## Build then serve production docs → http://localhost:3000
	npm run serve --prefix docs-site

docs-clean:    ## Remove docs build cache and output
	npm run clear --prefix docs-site
	rm -rf docs-site/build/

# ── Help ───────────────────────────────────────────────────────────

help:  ## Show this help
	@printf '$(BOLD)oat-latte$(RESET)\n\n'
	@printf '$(BOLD)$(CYAN)Go$(RESET)\n'
	@awk '/^# ── Go/,/^# ── Examples/' $(MAKEFILE_LIST) \
		| grep -E '^[a-zA-Z_-]+:.*##' \
		| awk -F ':.*##' '{ printf "  $(GREEN)%-16s$(RESET) $(GRAY)%s$(RESET)\n", $$1, $$2 }'
	@printf '\n$(BOLD)$(CYAN)Examples$(RESET)\n'
	@awk '/^# ── Examples/,/^# ── Cleanup/' $(MAKEFILE_LIST) \
		| grep -E '^[a-zA-Z_-]+:.*##' \
		| awk -F ':.*##' '{ printf "  $(GREEN)%-16s$(RESET) $(GRAY)%s$(RESET)\n", $$1, $$2 }'
	@printf '\n$(BOLD)$(CYAN)Cleanup$(RESET)\n'
	@awk '/^# ── Cleanup/,/^# ── Docs/' $(MAKEFILE_LIST) \
		| grep -E '^[a-zA-Z_-]+:.*##' \
		| awk -F ':.*##' '{ printf "  $(GREEN)%-16s$(RESET) $(GRAY)%s$(RESET)\n", $$1, $$2 }'
	@printf '\n$(BOLD)$(CYAN)Docs$(RESET)\n'
	@awk '/^# ── Docs/,/^# ── Help/' $(MAKEFILE_LIST) \
		| grep -E '^[a-zA-Z_-]+:.*##' \
		| awk -F ':.*##' '{ printf "  $(GREEN)%-16s$(RESET) $(GRAY)%s$(RESET)\n", $$1, $$2 }'
	@printf '\n'

.DEFAULT_GOAL := help
