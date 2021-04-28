package generate

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	mockredis "url-shortener/internal/repository/redis/mocks"
)

type TSuite struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TSuite))
}

var (
	mockRedis *mockredis.MockHandlerInterface
)

func setUpServiceMocking(ctrl *gomock.Controller) StorageService {
	mockRedis = mockredis.NewMockHandlerInterface(ctrl)

	return StorageService{
		RedisHandler: mockRedis,
	}
}

func (s *TSuite) TestGenerate_GetValidReqError() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	StorageService := setUpServiceMocking(ctrl)

	mockReqBody := `{
		"short_code": "test"
	}`

	w := httptest.NewRecorder()
	testRequest := httptest.NewRequest(http.MethodPost, "/generate", strings.NewReader(mockReqBody))
	testRequest.Header.Set("Content-Type", "application/json")
	StorageService.GenerateUrlShortener(w, testRequest)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	s.Assert().Equal(http.StatusBadRequest, w.Code)
	s.Assert().Contains(string(body), `"code":1001`)
}

func (s *TSuite) TestGenerate_GetValidRequireField() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	StorageService := setUpServiceMocking(ctrl)

	mockReqBody := `{
		"short_code": "",
		"full_url": "",
		"expire_date":0,
		"number_of_hits": 0
	}`

	w := httptest.NewRecorder()
	testRequest := httptest.NewRequest(http.MethodPost, "/generate", strings.NewReader(mockReqBody))
	testRequest.Header.Set("Content-Type", "application/json")
	StorageService.GenerateUrlShortener(w, testRequest)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	s.Assert().Equal(http.StatusBadRequest, w.Code)
	s.Assert().Contains(string(body), `"code":1001`)
}

func (s *TSuite) TestGenerate_GetValidateURLFormatField() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	StorageService := setUpServiceMocking(ctrl)

	mockReqBody := `{
		"short_code": "test",
		"full_url": "test",
		"expire_date":123123,
		"number_of_hits": 10
	}`

	w := httptest.NewRecorder()
	testRequest := httptest.NewRequest(http.MethodPost, "/generate", strings.NewReader(mockReqBody))
	testRequest.Header.Set("Content-Type", "application/json")
	StorageService.GenerateUrlShortener(w, testRequest)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	s.Assert().Equal(http.StatusBadRequest, w.Code)
	s.Assert().Contains(string(body), `"code":1001`)
}

func (s *TSuite) TestGenerate_GetValidUnmarshalError() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	StorageService := setUpServiceMocking(ctrl)

	mockReqBody := `{
		"a": "test",
		"b": "https://www.speedtest.net",
		"c":1619766384,
		"d": 10
	}`

	w := httptest.NewRecorder()
	testRequest := httptest.NewRequest(http.MethodPost, "/generate", strings.NewReader(mockReqBody))
	testRequest.Header.Set("Content-Type", "application/json")
	StorageService.GenerateUrlShortener(w, testRequest)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	s.Assert().Equal(http.StatusBadRequest, w.Code)
	s.Assert().Contains(string(body), `"code":1001`)
}

func (s *TSuite) TestGenerate_GetValidExpireError() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	StorageService := setUpServiceMocking(ctrl)

	mockReqBody := `{
		"short_code": "test",
		"full_url": "https://www.speedtest.net",
		"expire_date":1219766384,
		"number_of_hits": 10
	}`

	w := httptest.NewRecorder()
	testRequest := httptest.NewRequest(http.MethodPost, "/generate", strings.NewReader(mockReqBody))
	testRequest.Header.Set("Content-Type", "application/json")
	StorageService.GenerateUrlShortener(w, testRequest)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	s.Assert().Equal(http.StatusBadRequest, w.Code)
	s.Assert().Contains(string(body), `"code":1003`)
}

func (s *TSuite) TestGenerate_GetValidError() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	StorageService := setUpServiceMocking(ctrl)
	mockRedis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("redis error"))
	mockRedis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("redis error"))
	mockRedis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("redis error"))

	mockReqBody := `{
		"short_code": "test",
		"full_url": "https://www.speedtest.net",
		"expire_date":1619766384,
		"number_of_hits": 10
	}`

	w := httptest.NewRecorder()
	testRequest := httptest.NewRequest(http.MethodPost, "/generate", strings.NewReader(mockReqBody))
	testRequest.Header.Set("Content-Type", "application/json")
	StorageService.GenerateUrlShortener(w, testRequest)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	s.Assert().Equal(http.StatusBadRequest, w.Code)
	s.Assert().Contains(string(body), `"code":1004`)
}

func (s *TSuite) TestGenerate_URLSuccess() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	StorageService := setUpServiceMocking(ctrl)
	mockRedis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	mockReqBody := `{
		"short_code": "test",
		"full_url": "https://www.testlongtestlongtestlongtestlongtestlongtestlong.net",
		"expire_date":1619766384,
		"number_of_hits": 10
	}`

	w := httptest.NewRecorder()
	testRequest := httptest.NewRequest(http.MethodPost, "/generate", strings.NewReader(mockReqBody))
	testRequest.Header.Set("Content-Type", "application/json")
	StorageService.GenerateUrlShortener(w, testRequest)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	s.Assert().Equal(http.StatusCreated, w.Code)
	s.Assert().Contains(string(body), `"code":302`)
}
