#!/usr/bin/env bash

# Add it to pre-commit hook:
# echo -e "#\!/usr/bin/env bash\n\n./misc/pre-commit.sh" > .git/hooks/pre-commit && chmod +x .git/hooks/pre-commit

set -e

trap 'echo -e "\033[0;31mFAILED\033[0m"' ERR

golangci-lint run --fix
golangci-lint run --fix --build-tags postgres
golangci-lint run --fix --build-tags redis

go test -race -v ./...

echo -e "\033[0;32mSUCCESS\033[0m"
