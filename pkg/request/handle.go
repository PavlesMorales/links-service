package req

import (
	res "links-service/pkg/response"
	"net/http"
)

func HandleRq[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)

	if err != nil {
		res.JsonResponse(http.StatusBadRequest, *w, err.Error())
		return nil, err
	}
	err = IsValid(body)
	if err != nil {
		res.JsonResponse(http.StatusBadRequest, *w, err.Error())
		return nil, err
	}
	return &body, nil
}
