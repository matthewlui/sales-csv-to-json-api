package service

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"wilfred404/sales-csv-to-json-api/model"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Service) RetrieveSalesRecordHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	startDateStr := ps.ByName("startDate")
	endDateStr := ps.ByName("endDate")
	var err error
	var startDate time.Time
	var endDate time.Time
	filter := bson.M{}
	salesRecords := make([]*model.SalesRecord, 0)

	if len(startDateStr) > 0 {
		startDate, err = time.Parse(time.RFC3339, strings.Replace(startDateStr, "\n", "", -1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filter = bson.M{"lastpurchasedate": bson.M{
			"$gte": startDate,
		}}
	}
	if len(endDateStr) > 0 {
		endDate, err = time.Parse(time.RFC3339, strings.Replace(endDateStr, "\n", "", -1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filter = bson.M{"lastpurchasedate": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		}}
	}

	cursor, err := s.DB.Find(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer cursor.Close()
	for cursor.Next() {
		var salesRecord *model.SalesRecord
		if err = cursor.Decode(&salesRecord); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		salesRecords = append(salesRecords, salesRecord)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(salesRecords)
}
