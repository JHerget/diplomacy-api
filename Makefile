build: format
	rm -rf dist
	mkdir -p dist/board dist/orders dist/test

	GOOS=linux GOARCH=arm64 go build -o dist/turns/bootstrap ./cmd/turns
	GOOS=linux GOARCH=arm64 go build -o dist/games/bootstrap ./cmd/games
	GOOS=linux GOARCH=arm64 go build -o dist/maps/bootstrap ./cmd/maps
	GOOS=linux GOARCH=arm64 go build -o dist/phases/bootstrap ./cmd/phases
	GOOS=linux GOARCH=arm64 go build -o dist/players/bootstrap ./cmd/players
	GOOS=linux GOARCH=arm64 go build -o dist/board/bootstrap ./cmd/board
	GOOS=linux GOARCH=arm64 go build -o dist/orders/bootstrap ./cmd/orders
	GOOS=linux GOARCH=arm64 go build -o dist/test/bootstrap ./cmd/test

	cd dist/turns && zip ../turns.zip bootstrap
	cd dist/games && zip ../games.zip bootstrap
	cd dist/maps && zip ../maps.zip bootstrap
	cd dist/phases && zip ../phases.zip bootstrap
	cd dist/players && zip ../players.zip bootstrap
	cd dist/board && zip ../board.zip bootstrap
	cd dist/orders && zip ../orders.zip bootstrap
	cd dist/test && zip ../test.zip bootstrap

publish: build
	terraform -chdir=terraform apply

start: build
	go build -o dist/local ./cmd/local
	AWS_REGION=us-west-2 go run ./cmd/local

format:
	gofmt -w .

tidy:
	go mod tidy
