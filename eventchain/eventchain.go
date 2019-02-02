package eventchain

import (
	"log"
	"time"

	"github.com/kyeett/fruit-planet/condition"

	"github.com/fogleman/ease"
	"github.com/hajimehoshi/ebiten"
	"github.com/peterhellberg/gfx"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/eventsystem"
	"github.com/kyeett/gomponents/animation"
	"github.com/kyeett/gomponents/components"
)

type PlayerDeath struct {
	eventsystem.Basic
	ID                 string
	currentAnimationID string
}

func NewPlayerDeath() *PlayerDeath {
	basic := eventsystem.Basic{}
	return &PlayerDeath{
		Basic: basic,
		ID:    "player_1",
	}
}

func (pd *PlayerDeath) Next(em *entity.Manager) {
	names := map[int]string{
		0: "animate_out",
		1: "animate_in",
		2: "restore",
	}

	switch names[pd.CurrentStep] {
	case "animate_out":
		// Save components
		pd.ComponentStack.Push(*em.Drawable(pd.ID))
		pd.ComponentStack.Push(*em.Hitbox(pd.ID))
		pd.ComponentStack.Push(components.Joystick{})

		em.Remove(pd.ID, components.DrawableType)
		em.Remove(pd.ID, components.HitboxType)
		em.Remove(pd.ID, components.JoystickType)

		pos := em.Pos(pd.ID)
		pd.currentAnimationID = NewAnimation(em, "death_out", pos.Vec)
		pd.SetWaitingForID(NewAnimationCondition(em, pd.currentAnimationID))
	case "animate_in":
		target := em.Pos("initial_1")
		em.RemoveEntity(pd.currentAnimationID)
		pd.currentAnimationID = NewAnimation(em, "death_in", target.Vec)
		pd.SetWaitingForID(NewAnimationCondition(em, pd.currentAnimationID))

	case "restore":
		target := em.Pos("initial_1")
		em.RemoveEntity(pd.currentAnimationID)
		pos := em.Pos(pd.ID)
		pos.Vec = target.Vec
		pd.currentAnimationID = ""

		po := pd.ComponentStack.Pop()
		po2 := pd.ComponentStack.Pop()
		po3 := pd.ComponentStack.Pop()
		em.Add(pd.ID, po, po2, po3)
		pd.SetDone(true)
	default:
		pd.SetDone(true)
	}
	pd.CurrentStep++
}

func NewWaitCondition(em *entity.Manager, duration time.Duration) string {
	e := em.NewEntity("condition")
	em.Add(e, components.Trigger{
		Name:       e,
		Conditions: []components.Condition{condition.WaitUntil{time.Now().Add(duration)}}, // [][]string{{"animation_complete", ID}},
	})
	return e
}

func NewAnimation(em *entity.Manager, name string, v gfx.Vec) string {

	switch name {
	case "death_in", "death_out":
		img, err := gfx.OpenPNG("assets/images/frames.png")
		if err != nil {
			log.Fatal(err)
		}
		eImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		animID := em.NewEntity("animation")
		anim := animation.New(eImg, 64, 64)
		anim.Easing = ease.OutCubic
		if name == "death_in" {
			anim.Direction = -1
		} else {
			anim.Direction = 1
		}

		em.Add(animID, anim)
		em.Add(animID, components.Drawable{})
		em.Add(animID, components.Pos{v.Add(gfx.V(-22, -22))})
		return animID
	}
	return ""
}

func NewAnimationCondition(em *entity.Manager, ID string) string {
	e := em.NewEntity("condition")
	em.Add(e, components.Trigger{
		Name:       ID,
		Conditions: []components.Condition{condition.NewAnimationComplete(em, ID)}, // [][]string{{"animation_complete", ID}},
	})
	return e
}
