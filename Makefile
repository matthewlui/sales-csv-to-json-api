.PHONY: start

start:
	docker build -t sales-csv-to-json-api .
	docker-compose up -d

test:
	go test -v -cover ./...

mod-tidy:
	go mod tidy

gen-mock:
	mockgen -destination=mocks/sales_record_db_mock.go -package=mocks -source=./repo/sales_record_db.go
	mockgen -destination=mocks/cursor.go -package=mocks -source=./repo/cursor.go

build-docker:
	docker build -t sales-csv-to-json-api .
