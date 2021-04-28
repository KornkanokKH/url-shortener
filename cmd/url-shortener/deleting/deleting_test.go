package deleting

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func (s *TSuite) TestDelete_ValidRedisError() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	StorageService := setUpServiceMocking(ctrl)
	mockRedis.EXPECT().Del(gomock.Any(), gomock.Any()).Return(errors.New("redis error"))

	w := httptest.NewRecorder()
	testRequest := httptest.NewRequest(http.MethodDelete, "/code", nil)
	StorageService.DeleteUrlShortener(w, testRequest)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	s.Assert().Equal(http.StatusBadRequest, w.Code)
	s.Assert().Contains(string(body), `"code":1004`)
}

func (s *TSuite) TestDelete_Success() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	StorageService := setUpServiceMocking(ctrl)
	mockRedis.EXPECT().Del(gomock.Any(), gomock.Any()).Return(nil)

	w := httptest.NewRecorder()
	testRequest := httptest.NewRequest(http.MethodDelete, "/code", nil)
	StorageService.DeleteUrlShortener(w, testRequest)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	s.Assert().Equal(http.StatusGone, w.Code)
	s.Assert().Contains(string(body), `"code":200`)
}
