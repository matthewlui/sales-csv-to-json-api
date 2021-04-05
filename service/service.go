package service

import (
	"wilfred404/sales-csv-to-json-api/repo"
)

type Service struct {
	DB repo.SalesRecordDB
}

func ProvideService() *Service {
	host := "host.docker.internal"
	dbEndPoint := "mongodb://" + host + ":27017"
	db := "salesDB"
	collection := "sales_records"

	return &Service{repo.NewSalesRecordDB(dbEndPoint, db, collection)}
}
