package deleting

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
	"net/http"
	"url-shortener/internal/http/rest"
	"url-shortener/internal/repository/redis"
)

const fmtError = "%v"

type Service interface {
	DeleteUrlShortener(w http.ResponseWriter, r *http.Request)
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type StorageService struct {
	RedisHandler redis.HandlerInterface
	RedisConfig  redis.Config
}

func NewService(redisHandler *redis.Handler, redisConfig redis.Config) Service {
	return &StorageService{
		RedisHandler: redisHandler,
		RedisConfig:  redisConfig,
	}
}

func (s *StorageService) DeleteUrlShortener(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteUrlShortener")

	ctx := r.Context()
	txn := newrelic.FromContext(ctx)
	code := mux.Vars(r)["code"]

	respErr := &rest.ErrorResponse{
		Error: rest.Response{},
	}
	// step: delete key
	key := fmt.Sprintf("%v%v%v", s.RedisConfig.Key, code, "*")
	err := s.RedisHandler.Del(key, txn)
	if err != nil {
		msg := fmt.Sprintf("redis error (%v)", err)
		log.Error().Msgf(fmtError, msg)
		respErr.Error.Code = rest.ErrCodeBadRequest["Code"].(int)
		respErr.Error.Message = msg
		rest.WriteResponse(w, http.StatusBadRequest, respErr)
		return
	}

	fmt.Println("DeleteUrlShortener : Success")
	_ = rest.WriteResponse(w, http.StatusCreated, &rest.Response{
		Code:    200,
		Message: "Success",
	})

	return
}
