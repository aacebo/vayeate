ifneq (,$(wildcard ./.env))
    include .env
    export
endif

clean:
	rm -rf ./bin

build:
	go build -o bin/main main.go

clean_build: clean build

run:
	go run main.go

compile:
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 main.go

migrate_up:
	migrate -source file://database/migrations -database $(POSTGRES_CONNECTION_STRING) up

migrate_down:
	migrate -source file://database/migrations -database $(POSTGRES_CONNECTION_STRING) down

migrate_new:
	migrate create -ext sql -dir database/migrations $(name)
