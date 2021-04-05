package main

import (
	"fmt"
	"log"
	"net/http"
	. "wilfred404/sales-csv-to-json-api/service"

	"github.com/julienschmidt/httprouter"
)

func main() {

	s := ProvideService()

	router := httprouter.New()
	router.POST("/sales/record", s.ReceiveSalesRecordHandler)
	router.GET("/sales/report", s.DummyHandler)

	fmt.Printf("Starting server at port 3000\n")
	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}
