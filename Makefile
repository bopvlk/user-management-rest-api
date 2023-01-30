.SILENT:

# build:
# 	go mod download && go build -o bin/bot ./cmd/telegrambot/

# docker-build: build
# 	docker build --rm -t telegram-bot .
	
# run: docker-build
# 	docker run -d --name go-telegram-bot telegram-bot



# lint:
# 	golangci-lint run ./...

# fix-imports:
# 	gogroup -order std,other,prefix=git.foxminded.com.ua/3_REST_API --rewrite $(find . -type f -name "*.go" | grep -v /vendor/ |grep -v /.git/)

# format:
# 	go vet ./...
# 	go fmt ./...