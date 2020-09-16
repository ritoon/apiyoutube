build:
	go build -tags=jsoniter .

run:
	docker-compose up -d
	go run main.go

