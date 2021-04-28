package getting

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
	"url-shortener/internal/http/rest"
	"url-shortener/internal/repository/redis"
)

const fmtError = "%v"

type Service interface {
	GetUrlShortener(w http.ResponseWriter, r *http.Request)
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

func (s *StorageService) GetUrlShortener(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUrlShortener")

	ctx := r.Context()
	txn := newrelic.FromContext(ctx)
	code := mux.Vars(r)["code"]

	respErr := &rest.ErrorResponse{
		Error: rest.Response{},
	}
	// step: check data from redis
	// check exp
	exp := fmt.Sprintf("%v%v:%v", s.RedisConfig.Key, code, "expire")
	expireTime, err := s.RedisHandler.Get(exp, txn)
	if err != nil {
		msg := fmt.Sprintf("redis error (%v)", err)
		log.Error().Msgf(fmtError, msg)
		respErr.Error.Code = rest.ErrCodeBadRequest["Code"].(int)
		respErr.Error.Message = msg
		rest.WriteResponse(w, http.StatusBadRequest, respErr)
		return
	}

	times, _ := strconv.ParseInt(expireTime, 0, 64)
	if expireTime != "" && times <= time.Now().Unix() {
		msg := fmt.Sprintf("redis error (%v)", err)
		log.Error().Msgf(fmtError, msg)
		respErr.Error.Code = rest.ErrCodeUrlExp["Code"].(int)
		respErr.Error.Message = rest.ErrCodeUrlExp["Message"].(string)
		rest.WriteResponse(w, http.StatusBadRequest, respErr)
		return
	}

	// check found
	url := fmt.Sprintf("%v%v:%v", s.RedisConfig.Key, code, "full")
	urlRes, err := s.RedisHandler.Get(url, txn)
	if err != nil {
		msg := fmt.Sprintf("redis error (%v)", err)
		log.Error().Msgf(fmtError, msg)
		respErr.Error.Code = rest.ErrCodeBadRequest["Code"].(int)
		respErr.Error.Message = msg
		rest.WriteResponse(w, http.StatusBadRequest, respErr)
		return
	}

	if urlRes == "" {
		msg := fmt.Sprintf("url not found (%v)", err)
		log.Error().Msgf(fmtError, msg)
		respErr.Error.Code = rest.ErrCodeNotfound["Code"].(int)
		respErr.Error.Message = rest.ErrCodeNotfound["Message"].(string)
		rest.WriteResponse(w, http.StatusNotFound, respErr)
		return
	}

	fmt.Println("GetUrlShortener : Success")
	// Redirect
	http.Redirect(w, r, urlRes, 302)

	return
}
