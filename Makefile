build:
	go build -o shopify-challenge-2022 ./cmd

test:
	go test -v ./models
	go test -v ./handlers