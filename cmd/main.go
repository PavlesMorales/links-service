package main

import (
	"fmt"
	"links-service/configs"
	"links-service/internal/auth"
	"links-service/internal/link"
	"links-service/internal/stat"
	"links-service/internal/user"
	"links-service/pkg/db"
	"links-service/pkg/event"
	"links-service/pkg/middleware"
	"net/http"
)

func main() {

	conf := configs.LoadConfig()
	dbConf := db.NewDb(conf)
	eventBus := event.NewEventBus()

	lRep := link.NewLinkRepository(dbConf)
	uRep := user.NewUserRepository(dbConf)
	ser := auth.NewAuthService(uRep)
	clickRep := stat.NewStatRepository(dbConf)

	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: clickRep,
	})

	router := http.NewServeMux()
	ahd := auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: ser,
	}
	auth.NewAuthHandler(router, ahd)

	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: lRep,
		Config:         conf,
		EventBus:       eventBus,
	})

	stat.NewStatHandler(router, stat.StatHandlerDeps{StatRepository: clickRep, Config: conf})

	stack := middleware.Chain(
		middleware.Cors,
		middleware.Logging,
	)

	go statService.AddClickSubscriber()
	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("Server started on port", server.Addr)

	err := server.ListenAndServe()
	fmt.Println("Server started after ListenAndServe")

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
