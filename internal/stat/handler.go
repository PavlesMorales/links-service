package stat

import (
	"links-service/configs"
	"links-service/pkg/middleware"
	req "links-service/pkg/request"
	res "links-service/pkg/response"
	"net/http"
)

type StatHandler struct {
	StatRepository *StatRepository
}

type StatHandlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}
	router.Handle("GET /stat", middleware.Auth(handler.Stat(), *deps.Config))

}

func (h StatHandler) Stat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statParams, err := req.GetStatParams(r.URL.Query())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stat := h.StatRepository.GetStats(*statParams)
		res.JsonResponse(http.StatusOK, w, stat)
	}
}
