package world

import (
	"errors"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/ecs/camera"
	"github.com/kyeett/ecs/constants"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/ecs/rendersystem"
	"github.com/kyeett/ecs/system"
	"github.com/sirupsen/logrus"
)

// World holds the ECS system, cameras and the current scene
type World struct {
	eventCh       chan events.Event
	camera        *camera.Camera
	systems       []system.System
	renderSystems []rendersystem.System
	em            *entity.Manager
	mapName       string
	canvas        Canvas
}

// New returns an initiated world, with camera width x height
func New(m string, width, height int) *World {
	em := entity.NewManager(logging.NewLogger())
	eventCh := make(chan events.Event, 100)
	w := World{
		eventCh: eventCh,
		em:      em,
		mapName: m,
		canvas:  NewCanvas(),
	}

	err := w.LoadScene("default")
	if err != nil {
		log.Fatal(err)
	}
	w.systems = []system.System{
		system.NewInput(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
		system.NewControls(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
		system.NewMovement(em, logging.NewLogger(logrus.InfoLevel)),
	}

	w.renderSystems = []rendersystem.System{
		rendersystem.NewRenderImage(w.canvas.renderers["background"], logging.NewLogger()),
		rendersystem.NewRender(em, logging.NewLogger()),
	}
	return &w
}

// Reset the world and it's entitites to its original state
func (w *World) Reset() {
	w.em.Reset()
	w.camera.Reset()
	w.LoadScene("default")
}

// StartEventQueue starts the event queue for the ECS system
func (w *World) StartEventQueue() {
	go func() {
		for {
			ev := <-w.eventCh
			for _, s := range w.systems {
				s.Send(ev)
			}
		}
	}()
}

var timeStep = 1.0

// Update and redraw the world
func (w *World) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		return errors.New("exit game")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		w.Reset()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		timeStep = constants.DefaultTimeStep - timeStep
	}

	for _, s := range w.systems {
		s.Update(timeStep)
	}

	for _, s := range w.renderSystems {
		s.Update(screen)
	}

	// r, op := w.camera.View(timeStep)
	// screen.DrawImage(screen, op)
	return nil
}
