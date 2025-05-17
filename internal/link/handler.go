package link

import (
	"fmt"
	"gorm.io/gorm"
	"links-service/configs"
	"links-service/pkg/event"
	"links-service/pkg/middleware"
	req "links-service/pkg/request"
	res "links-service/pkg/response"
	"net/http"
	"strconv"
)

type LinkHandler struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
}

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
	Config         *configs.Config
	EventBus       *event.EventBus
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}
	router.HandleFunc("GET /{hash}", handler.GoTo())
	router.HandleFunc("POST /link", handler.Create())
	router.Handle("PATCH /link/{id}", middleware.Auth(handler.Update(), *deps.Config))
	router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.Handle("GET /link", middleware.Auth(handler.Get(), *deps.Config))
}

func (h *LinkHandler) Get() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		params, err := req.GetParams(r.URL.Query())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		links := h.LinkRepository.GetAll(params.Limit, params.Offset, params.Order)
		count := h.LinkRepository.Count()
		res.JsonResponse(http.StatusOK, w, &GetAllLinksResponse{
			Links: links,
			Count: count,
		})
	}
}

func (h *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := req.HandleRq[LinkCreateRq](&w, r)
		if err != nil {
			return
		}
		link := NewLink(req.Url)
		for {
			byHash, _ := h.LinkRepository.GetByHash(link.Hash)
			if byHash == nil {
				break
			}
			link.GenerateHash()
			fmt.Println("With hash already exists: ", byHash.Hash, ". Generate new: ", link.Hash)
		}

		createdLink, err := h.LinkRepository.Create(link)
		if err != nil {
			res.JsonResponse(http.StatusBadRequest, w, err.Error())
			return
		}
		res.JsonResponse(http.StatusCreated, w, createdLink)
	}
}

func (h *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if email, ok := r.Context().Value(middleware.ContextEmailKey).(string); !ok {
			fmt.Println("Not found email")
		} else {
			fmt.Println("With email: ", email)
		}
		body, err := req.HandleRq[LinkUpdateRq](&w, r)
		if err != nil {
			res.JsonResponse(http.StatusBadRequest, w, err.Error())
			return
		}
		idRaw := r.PathValue("id")
		id, err := strconv.ParseUint(idRaw, 10, 64)
		if err != nil {
			res.JsonResponse(http.StatusBadRequest, w, err.Error())
			return
		}
		link, err := h.LinkRepository.Update(
			&Link{
				Model: gorm.Model{ID: uint(id)},
				Url:   body.Url,
				Hash:  body.Hash,
			})
		if err != nil {
			res.JsonResponse(http.StatusBadRequest, w, err.Error())
			return
		}
		res.JsonResponse(http.StatusCreated, w, link)

	}
}

func (h *LinkHandler) GoTo() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		link, err := h.LinkRepository.GetByHash(hash)
		if err != nil {
			res.JsonResponse(http.StatusNotFound, w, err.Error())
			return
		}

		go h.EventBus.Publish(event.Event{
			Type: event.LinkVisited,
			Data: link.ID,
		})

		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}

func (h *LinkHandler) Delete() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		idRaw := r.PathValue("id")
		id, err := strconv.ParseUint(idRaw, 10, 64)
		if err != nil {
			res.JsonResponse(http.StatusBadRequest, w, err.Error())
			return
		}
		err = h.LinkRepository.GetById(uint(id))
		if err != nil {
			res.JsonResponse(http.StatusNotFound, w, err.Error())
			return
		}

		err = h.LinkRepository.DeleteById(id)
		if err != nil {
			res.JsonResponse(http.StatusInternalServerError, w, err.Error())
			return
		}
		res.JsonResponse(http.StatusNoContent, w, nil)
	}
}
