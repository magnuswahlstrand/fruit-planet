package events

import (
	"github.com/kyeett/fruit-planet/eventchain"
	"github.com/kyeett/gomponents/components"

	"github.com/kyeett/ecs/eventsystem"

	"github.com/kyeett/ecs/entity"
)

type Mapper struct{}

var eventQueue []interface{}

func (m Mapper) Handle(em *entity.Manager, event interface{}) {
	eventQueue = append(eventQueue, event)
}

func HandleEvents(em *entity.Manager) {
	for _, event := range eventQueue {
		switch v := event.(type) {
		case eventsystem.Collision:
			if em.HasComponents(v.ID2, components.HazardType) {
				pd := eventchain.NewPlayerDeath()
				eventsystem.AddEventChain(pd)
				pd.Next(em)
			}
		}
	}
	eventQueue = []interface{}{}
}
