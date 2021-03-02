package generate

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"url-shortener/internal/http/rest"
)

const fmtError = "%v"

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Service interface {
	GenerateUrlShortener(w http.ResponseWriter, r *http.Request)
}

type service struct {
}

func NewService() Service {
	return &service{
	}
}

func (s *service) GenerateUrlShortener(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GenerateUrlShortener")

	respErr := &rest.ErrorResponse{
		Error: rest.Response{},
	}

	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("ioutil.ReadAll (%v)", err)
		log.Error().Msgf(fmtError, msg)
		respErr.Error.Code = rest.ErrCodeBadRequest["Code"].(int)
		respErr.Error.Message = msg
		return
	}

	_ = rest.WriteResponse(w, http.StatusCreated, &rest.Response{
		Code:            200,
		Message:         "Success",
	})

	return
}
