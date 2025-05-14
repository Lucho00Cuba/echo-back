# Define the color codes for output
YELLOW = \033[33m
GREEN = \033[32m
CYAN = \033[36m
RESET = \033[0m

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

TARGETS := $(basename $(notdir $(wildcard scripts/*)))
TARGET := $(firstword $(MAKECMDGOALS))
TARGETS_FILES := $(wildcard scripts/*)
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

.PHONY: help

all: help

### Help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo "  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}"
	@echo ''
	@echo "${CYAN}Targets${RESET}"

	@# Print targets from scripts with comments
	@for target in $(TARGETS_FILES); do \
		script_name=$$(basename $$target); \
		description=$$(awk '/^## makefile:fmt/ {sub(/^## makefile:fmt /, ""); print $$0}' $$target); \
		if [ -n "$$description" ]; then \
			printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n" "$$script_name" "$$description"; \
		fi; \
	done

	@# Print inline help from the Makefile itself
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) { \
			sub(/^## makefile:fmt /, ""); \
			printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2; \
		} \
	}' $(MAKEFILE_LIST)

### Ensure all scripts have execution permissions
ensure-permissions:
	@for script in scripts/*; do \
	    if [ ! -x $$script ]; then \
	        echo "Adding execute permissions to $$script"; \
	        chmod +x $$script; \
	    fi \
	done

%::
	@true

$(TARGETS): ensure-permissions
	@./scripts/$@ $(ARGS)
