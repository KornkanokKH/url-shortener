package generate

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

func (s *TSuite) TestDeleteSubscription_GetValidError() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	StorageService := setUpServiceMocking(ctrl)
	mockRedis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))

	w := httptest.NewRecorder()
	testRequest := httptest.NewRequest(http.MethodPost, "/something", nil)
	StorageService.GenerateUrlShortener(w, testRequest)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	s.Assert().Equal(http.StatusBadRequest, w.Code)
	s.Assert().Contains(string(body), `"code":1004`)
}
