package stat

import (
	"links-service/pkg/event"
	"log"
)

type StatServiceDeps struct {
	EventBus       *event.EventBus
	StatRepository *StatRepository
}

type StatService struct {
	EventBus       *event.EventBus
	StatRepository *StatRepository
}

func NewStatService(deps *StatServiceDeps) *StatService {
	return &StatService{
		EventBus:       deps.EventBus,
		StatRepository: deps.StatRepository,
	}
}

func (s *StatService) AddClickSubscriber() {
	for msg := range s.EventBus.Subscribe() {
		log.Printf("Receive msg: %s", msg.Type)
		if msg.Type != event.LinkVisited {
			continue
		}

		id, ok := msg.Data.(uint)
		if !ok {
			log.Println("Bad link visited. Data: ", msg.Data)
			continue
		}

		s.StatRepository.AddClick(id)
	}
}
