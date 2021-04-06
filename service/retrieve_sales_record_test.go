package service_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"wilfred404/sales-csv-to-json-api/mocks"
	"wilfred404/sales-csv-to-json-api/model"
	. "wilfred404/sales-csv-to-json-api/service"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/suite"
)

type retrieveSalesRecordTestSuite struct {
	suite.Suite

	Service *Service
}

func TestRetrieveSalesRecordTestSuite(t *testing.T) {
	suite.Run(t, new(retrieveSalesRecordTestSuite))
}

func getSalesRecord() *model.SalesRecord {
	lastPurchaseDate, _ := time.Parse(time.RFC3339, "2020-11-05T13:15:30Z")
	return &model.SalesRecord{UserName: "UserName", Age: 18, Height: 180, Gender: "M", SaleAmount: 12345, LastPurchaseDate: lastPurchaseDate}
}

func (suite *retrieveSalesRecordTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.Suite.T())
	defer ctrl.Finish()

	mockCursor := mocks.NewMockCursor(ctrl)

	mockCursor.EXPECT().Decode(gomock.Any()).SetArg(0, getSalesRecord()).Return(nil).AnyTimes()

	i := 0
	mockCursor.EXPECT().Next().DoAndReturn(func() bool {
		if i < 1 {
			i++
			return true
		}
		return false
	}).AnyTimes()

	mockCursor.EXPECT().Close().Return(nil).AnyTimes()

	mockDB := mocks.NewMockSalesRecordDB(ctrl)
	mockDB.EXPECT().Find(gomock.Any()).Return(mockCursor, nil).AnyTimes()

	suite.Service = &Service{mockDB}
}

func (suite *retrieveSalesRecordTestSuite) TestNoStartEndDateSuccess() {
	req, err := http.NewRequest("GET", "/sales/report", nil)
	if err != nil {
		log.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Service.RetrieveSalesRecordHandler(w, r, httprouter.Params{})
	})
	handler.ServeHTTP(rr, req)

	expectedRspBody, err := json.Marshal([]*model.SalesRecord{getSalesRecord()})
	if err != nil {
		log.Fatal(err)
	}
	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal(string(expectedRspBody)+"\n", rr.Body.String())
}

func (suite *retrieveSalesRecordTestSuite) TestValidStartDateSuccess() {
	const startDate = "2000-11-05T13:15:30Z"
	req, err := http.NewRequest("GET", "/sales/report/"+startDate, nil)
	if err != nil {
		log.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Service.RetrieveSalesRecordHandler(w, r, httprouter.Params{httprouter.Param{Key: "startDate", Value: startDate}})
	})
	handler.ServeHTTP(rr, req)

	expectedRspBody, err := json.Marshal([]*model.SalesRecord{getSalesRecord()})
	if err != nil {
		log.Fatal(err)
	}
	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal(string(expectedRspBody)+"\n", rr.Body.String())
}

func (suite *retrieveSalesRecordTestSuite) TestValidStartEndDateSuccess() {
	const startDate = "2000-11-05T13:15:30Z"
	const endDate = "2030-11-05T13:15:30Z"
	req, err := http.NewRequest("GET", "/sales/report/"+startDate+"/"+endDate, nil)
	if err != nil {
		log.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Service.RetrieveSalesRecordHandler(w, r, httprouter.Params{httprouter.Param{Key: "startDate", Value: startDate}, httprouter.Param{Key: "endDate", Value: endDate}})
	})
	handler.ServeHTTP(rr, req)

	expectedRspBody, err := json.Marshal([]*model.SalesRecord{getSalesRecord()})
	if err != nil {
		log.Fatal(err)
	}
	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal(string(expectedRspBody)+"\n", rr.Body.String())
}

func (suite *retrieveSalesRecordTestSuite) TestInvalidStartDateFail() {
	const startDate = "invalid start date"
	req, err := http.NewRequest("GET", "/sales/report/"+startDate, nil)
	if err != nil {
		log.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Service.RetrieveSalesRecordHandler(w, r, httprouter.Params{httprouter.Param{Key: "startDate", Value: startDate}})
	})
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusInternalServerError, rr.Code)
	suite.Suite.T().Log(rr.Body.String())
}
