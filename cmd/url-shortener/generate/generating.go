package generate

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	_ "github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	_ "strings"
	"time"
	_ "time"
	"url-shortener/internal/http/rest"
	"url-shortener/internal/repository/redis"
)

const fmtError = "%v"

type Service interface {
	GenerateUrlShortener(w http.ResponseWriter, r *http.Request)
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

func (s *StorageService) GenerateUrlShortener(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GenerateUrlShortener")

	ctx := r.Context()
	txn := newrelic.FromContext(ctx)

	respErr := &rest.ErrorResponse{
		Error: rest.Response{},
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("ioutil.ReadAll (%v)", err)
		log.Error().Msgf(fmtError, msg)
		respErr.Error.Code = rest.ErrCodeBadRequest["Code"].(int)
		respErr.Error.Message = msg
		return
	}

	request := new(ShortenerRequest)
	err = json.Unmarshal(bodyBytes, request)
	if err != nil {
		msg := fmt.Sprintf("json.Unmarshal (%v)", err)
		log.Error().Msgf(fmtError, msg)
		respErr.Error.Code = rest.ErrCodeBadRequest["Code"].(int)
		respErr.Error.Message = msg
		return
	}

	// Step : validate url input
	_, err = govalidator.ValidateStruct(request)
	if err != nil {
		var errList []string
		for i := range govalidator.ErrorsByField(err) {
			errList = append(errList, i)
		}
		msg := fmt.Sprintf("Invalid parameter (%v)", err)
		log.Error().Msgf(fmtError, msg)
		respErr.Error.Code = rest.ErrCodeBadRequest["Code"].(int)
		respErr.Error.Message = fmt.Sprintf(rest.ErrCodeBadRequest["Message"].(string), strings.Join(errList, ","))
		rest.WriteResponse(w, http.StatusBadRequest, respErr)
		return
	}

	uri, err1 := url.ParseRequestURI(request.FullURL)
	if err1 != nil && uri != nil && uri.Scheme != "http" && uri.Scheme != "https" {
		msg := fmt.Sprintf("Invalid URL (%v)", err1)
		log.Error().Msgf(fmtError, msg)
		respErr.Error.Code = rest.ErrCodeURLInvalid["Code"].(int)
		respErr.Error.Message = fmt.Sprintf(rest.ErrCodeBadRequest["Message"].(string), err1)
		rest.WriteResponse(w, http.StatusBadRequest, respErr)
		return

	}

	if request.ExpireDate <= time.Now().Unix() {
		msg := fmt.Sprintf("Invalid expire_date :%v <= %v", request.ExpireDate, time.Now().Unix())
		log.Error().Msgf(fmtError, msg)
		respErr.Error.Code = rest.ErrCodeDate["Code"].(int)
		respErr.Error.Message = rest.ErrCodeDate["Message"].(string)
		rest.WriteResponse(w, http.StatusBadRequest, respErr)
		return
	}

	// step : counter url on redis
	err = s.RedisHandler.Set(s.RedisConfig.Key, request.ShortCode, 0, txn)
	if err != nil {
		msg := fmt.Sprintf("error redis (%v)", err)
		log.Error().Msgf(fmtError, msg)
		respErr.Error.Code = rest.ErrCodeRedis["Code"].(int)
		respErr.Error.Message = fmt.Sprintf(rest.ErrCodeRedis["Message"].(string), err)
		rest.WriteResponse(w, http.StatusBadRequest, respErr)
		return

	}
	// todo : generate shorten url

	_ = rest.WriteResponse(w, http.StatusCreated, &rest.Response{
		Code:    200,
		Message: "Success",
	})

	return
}
