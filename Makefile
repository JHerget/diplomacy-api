build: format
	rm -rf dist
	mkdir -p dist/board dist/orders dist/test

	GOOS=linux GOARCH=arm64 go build -o dist/board/bootstrap ./cmd/board
	GOOS=linux GOARCH=arm64 go build -o dist/orders/bootstrap ./cmd/orders
	GOOS=linux GOARCH=arm64 go build -o dist/test/bootstrap ./cmd/test

	cd dist/board && zip ../board.zip bootstrap
	cd dist/orders && zip ../orders.zip bootstrap
	cd dist/test && zip ../test.zip bootstrap

publish: build
	terraform -chdir=terraform apply

start: format
	go build ./cmd/local -o dist
	AWS_REGION=us-west-2 go run ./cmd/local

format:
	gofmt -w .

tidy:
	go mod tidy
