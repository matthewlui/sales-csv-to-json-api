package service_test

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"wilfred404/sales-csv-to-json-api/mocks"
	. "wilfred404/sales-csv-to-json-api/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

type receiveSalesRecordTestSuite struct {
	suite.Suite

	Service *Service
}

func TestCreateCaseTestSuite(t *testing.T) {
	suite.Run(t, new(receiveSalesRecordTestSuite))
}

func (suite *receiveSalesRecordTestSuite) SetupSuite() {
	ctrl := gomock.NewController(suite.Suite.T())
	defer ctrl.Finish()

	mockDB := mocks.NewMockSalesRecordDB(ctrl)
	mockDB.EXPECT().InsertMany(gomock.Any()).Return(&mongo.InsertManyResult{InsertedIDs: []interface{}{"InsertedID_1", "InsertedID_2"}}, nil).AnyTimes()

	suite.Service = &Service{mockDB}
}

func createPostRequest(fileName string) *http.Request {
	path, err := filepath.Abs("../testing_files/" + fileName)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", "/sales/record", body)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	return request
}

func (suite *receiveSalesRecordTestSuite) TestSmallCSVSuccess() {
	req := createPostRequest("small.csv")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Service.ReceiveSalesRecordHandler(w, r, nil)
	})
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal("Upload Successful!", rr.Body.String())
}

func (suite *receiveSalesRecordTestSuite) TestBrokenCSVHeaderFail() {
	req := createPostRequest("broken_header.csv")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Service.ReceiveSalesRecordHandler(w, r, nil)
	})
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusInternalServerError, rr.Code)
	suite.Suite.T().Log(rr.Body.String())
}

func (suite *receiveSalesRecordTestSuite) TestValidContentWrongFileExtensionSuccess() {
	req := createPostRequest("valid_csv.abc")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Service.ReceiveSalesRecordHandler(w, r, nil)
	})
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal("Upload Successful!", rr.Body.String())
}

func (suite *receiveSalesRecordTestSuite) TestCSVMissingColumnDataFail() {
	req := createPostRequest("missing_column_data.csv")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Service.ReceiveSalesRecordHandler(w, r, nil)
	})
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusInternalServerError, rr.Code)
	suite.Suite.T().Log(rr.Body.String())
}

func (suite *receiveSalesRecordTestSuite) TestCSVIncorrectDataTypeFail() {
	req := createPostRequest("column_incorrect_data_type.csv")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Service.ReceiveSalesRecordHandler(w, r, nil)
	})
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusInternalServerError, rr.Code)
	suite.Suite.T().Log(rr.Body.String())
}

func (suite *receiveSalesRecordTestSuite) TestFailToObtainFile() {
	req, _ := http.NewRequest("POST", "/sales/record", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Service.ReceiveSalesRecordHandler(w, r, nil)
	})
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusInternalServerError, rr.Code)
	suite.Suite.T().Log(rr.Body.String())
}
