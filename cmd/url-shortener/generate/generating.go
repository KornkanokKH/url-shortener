package generate

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	_ "github.com/gin-gonic/gin"
	_ "github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"
	"url-shortener/internal/generate/encode"
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

	// Step : validate request
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

	// step : generate shorten url
	urlHashBytes := encode.Sha256Of(request.FullURL)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := encode.Base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))

	// step : set url on redis
	additional := make(map[string]interface{})
	additional[fmt.Sprintf("%v%v:%v", s.RedisConfig.Key, finalString, "expire")] = request.ExpireDate
	additional[fmt.Sprintf("%v%v:%v", s.RedisConfig.Key, finalString, "full")] = request.FullURL
	additional[fmt.Sprintf("%v%v:%v", s.RedisConfig.Key, finalString, "hits")] = request.NumberOfHits

	for k, v := range additional {
		err = s.RedisHandler.Set(k, v, 0, txn)
		if err != nil {
			msg := fmt.Sprintf("error redis (%v)", err)
			log.Error().Msgf(fmtError, msg)
			respErr.Error.Code = rest.ErrCodeRedis["Code"].(int)
			respErr.Error.Message = fmt.Sprintf(rest.ErrCodeRedis["Message"].(string), err)
			rest.WriteResponse(w, http.StatusBadRequest, respErr)
		}
	}

	host := fmt.Sprintf("http://%v/", r.Host)
	data := &ShortenerResponse{
		ShortCode: request.ShortCode,
		FullURL:   request.FullURL,
		ShortURL:  host + finalString,
	}
	_ = rest.WriteResponse(w, http.StatusCreated, &rest.Response{
		Code:    302,
		Message: "Success",
		Data:    data,
	})

	return
}
