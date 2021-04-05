package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	m "wilfred404/sales-csv-to-json-api/model"

	"github.com/julienschmidt/httprouter"
)

func (s *Service) ReceiveSalesRecordHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	const insertThreshold = 1000
	collection := s.DB

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	name := header.Filename

	fmt.Printf("File name %s\n", name)

	reader := bufio.NewReader(file)
	isFirstLine := true
	salesRecords := make([]interface{}, 0)
	for {
		line, readErr := reader.ReadString('\n')
		if readErr != nil && readErr != io.EOF {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if isFirstLine && err == io.EOF {
			return
		}

		if isFirstLine && !isHeaderValid(line) {
			http.Error(w, errors.New("invalid csv header").Error(), http.StatusInternalServerError)
			return
		}

		if isFirstLine {
			isFirstLine = false
			continue
		}

		salesRecord, err := ParseSalesRecord(line)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		salesRecords = append(salesRecords, salesRecord)

		if len(salesRecords) >= insertThreshold || readErr == io.EOF {
			insertManyResult, err := collection.InsertMany(salesRecords)
			if err != nil {
				//TODO: decide to rollback or not
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			salesRecords = make([]interface{}, 0)
			fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
		}

		if readErr == io.EOF {
			break
		}
	}
	fmt.Fprintf(w, "Upload Successful!")

}

func isHeaderValid(h string) bool {
	return strings.ToUpper(h) == "USER_NAME,AGE,HEIGHT,GENDER,SALE_AMOUNT,LAST_PURCHASE_DATE\n"
}

func ParseSalesRecord(l string) (*m.SalesRecord, error) {
	fmt.Println((l))
	line := strings.Split(l, ",")
	if len(line) != 6 {
		return nil, errors.New("insufficient column")
	}

	age, err := strconv.Atoi(line[1])
	if err != nil {
		return nil, err
	}

	height, err := strconv.ParseFloat(line[2], 32)
	if err != nil {
		return nil, err
	}

	gender := strings.ToUpper(line[3])
	if gender != "M" && gender != "F" {
		return nil, errors.New("error converting gender")
	}

	saleAmount, err := strconv.ParseFloat(line[4], 32)
	if err != nil {
		return nil, err
	}

	lastPurchaseDate, err := time.Parse(time.RFC3339, strings.Replace(line[5], "\n", "", -1))
	if err != nil {
		return nil, err
	}

	return &m.SalesRecord{UserName: line[0], Age: age, Height: float32(height), Gender: gender, SaleAmount: float32(saleAmount), LastPurchaseDate: lastPurchaseDate}, nil
}

func (s *Service) DummyHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello!")
}
