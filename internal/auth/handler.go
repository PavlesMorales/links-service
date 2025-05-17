package auth

import (
	"fmt"
	"links-service/configs"
	"links-service/pkg/jwt"
	req "links-service/pkg/request"
	res "links-service/pkg/response"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}
type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}

	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rq, err := req.HandleRq[LoginRq](&w, r)
		if err != nil {
			return
		}

		byEmail, err := handler.AuthService.Login(rq.Email, rq.Password)
		if err != nil {
			fmt.Println(err.Error())
			res.JsonResponse(http.StatusBadRequest, w, err.Error())
			return
		}

		jwt, err := jwt.NewJwt(handler.Config.Auth.Secret).Create(jwt.JwtData{Email: byEmail})
		if err != nil {
			fmt.Println(err.Error())
			res.JsonResponse(http.StatusInternalServerError, w, err.Error())
			return
		}

		response := LoginRs{
			Token: jwt,
		}
		res.JsonResponse(http.StatusCreated, w, response)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rq, err := req.HandleRq[RegisterRq](&w, r)
		if err != nil {
			fmt.Println(err.Error())
			res.JsonResponse(http.StatusBadRequest, w, err.Error())
			return
		}

		if err != nil {
			res.JsonResponse(http.StatusBadRequest, w, err.Error())
			return
		}
		email, err := handler.AuthService.Register(
			rq.Name,
			rq.Email,
			rq.Password,
		)
		if err != nil {
			fmt.Println(err.Error())
			res.JsonResponse(http.StatusBadRequest, w, err.Error())
			return
		}

		jwt, err := jwt.NewJwt(handler.Config.Auth.Secret).Create(jwt.JwtData{Email: email})
		if err != nil {
			fmt.Println(err.Error())
			res.JsonResponse(http.StatusInternalServerError, w, err.Error())
			return
		}

		response := LoginRs{
			Token: jwt,
		}
		res.JsonResponse(http.StatusCreated, w, response)
	}
}
