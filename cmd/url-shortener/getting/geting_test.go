package getting

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
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
