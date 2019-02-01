package eventchain

import (
	"fmt"
	"log"
	"time"

	"github.com/fogleman/ease"
	"github.com/hajimehoshi/ebiten"
	"github.com/peterhellberg/gfx"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/gomponents/animation"
	"github.com/kyeett/gomponents/components"
)

type PlayerDeath struct {
	Basic
	ID                 string
	currentAnimationID string
}

type Basic struct {
	step         int
	waitingForID string
	done         bool
	ComponentStack
}

type ComponentStack struct {
	stack []interface{}
}

func (cs *ComponentStack) Push(c interface{}) {
	cs.stack = append(cs.stack, c)
}

func (cs *ComponentStack) Pop() interface{} {
	c := cs.stack[len(cs.stack)-1]
	cs.stack = cs.stack[:len(cs.stack)-1]
	return c
}

func NewPlayerDeath() *PlayerDeath {
	basic := Basic{}
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
		3: "2nd",
		4: "3rd",
		5: "4th",
	}

	target := gfx.V(200, 100)
	switch names[pd.step] {
	case "animate_out", "2nd":
		fmt.Println("1")
		// Save components
		pd.ComponentStack.Push(*em.Drawable(pd.ID))
		pd.ComponentStack.Push(components.Joystick{})

		em.Remove(pd.ID, components.JoystickType)
		em.Remove(pd.ID, components.DrawableType)

		pos := em.Pos(pd.ID)
		pd.currentAnimationID = NewAnimation(em, "death", pos.Vec)
		pd.waitingForID = NewAnimationCondition(em, pd.currentAnimationID)
	case "animate_in", "3rd":
		fmt.Println("2")
		em.RemoveEntity(pd.currentAnimationID)
		pd.currentAnimationID = NewAnimation(em, "death", target)
		pd.waitingForID = NewAnimationCondition(em, pd.currentAnimationID)

	case "restore", "4th":
		fmt.Println("rest")
		em.RemoveEntity(pd.currentAnimationID)
		pos := em.Pos(pd.ID)
		pos.Vec = target
		pd.currentAnimationID = ""

		em.Add(pd.ID, pd.ComponentStack.Pop())
		em.Add(pd.ID, pd.ComponentStack.Pop())
		pd.waitingForID = NewWaitCondition(em, 2*time.Second)

	default:
		fmt.Println("done")
		pd.done = true
	}
	fmt.Println(names[pd.step])
	pd.step++
}

func (pd *Basic) Next(em *entity.Manager) bool {
	log.Fatal("not implemented")
	return true
}

func (pd *Basic) Done() bool {
	return pd.done
}

func (pd *Basic) WaitingForID() string {
	return pd.waitingForID
}

func NewWaitCondition(em *entity.Manager, duration time.Duration) string {
	e := em.NewEntity("condition")
	cond := components.Condition{
		Name: e, // For debugging
	}
	targetTime := time.Now().Add(duration)
	cond.Conditions = append(cond.Conditions, []string{"wait_until", targetTime.Format(time.RFC3339Nano)})
	em.Add(e, cond)
	return e
}

func NewAnimation(em *entity.Manager, name string, v gfx.Vec) string {

	switch name {
	case "death":
		img, err := gfx.OpenPNG("assets/images/frames.png")
		if err != nil {
			log.Fatal(err)
		}
		eImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

		animID := em.NewEntity("animation")
		anim := animation.New(eImg, 64, 64)
		anim.Easing = ease.OutCubic
		anim.Direction = 1
		em.Add(animID, anim)
		em.Add(animID, components.Drawable{})
		em.Add(animID, components.Pos{v.Add(gfx.V(-22, -22))})
		return animID
	}
	return ""
}

func NewAnimationCondition(em *entity.Manager, ID string) string {
	e := em.NewEntity("condition")
	cond := components.Condition{
		Name:       ID,
		Conditions: [][]string{{"animation_complete", ID}},
	}
	em.Add(e, cond)
	return e
}
