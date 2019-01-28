package entityloader

import (
	"fmt"
	"image"
	"log"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/gomponents/components"
	tiled "github.com/lafriks/go-tiled"
	"github.com/peterhellberg/gfx"
)

var headless bool

func Hitbox(em *entity.Manager, o *tiled.Object) {
	e := em.NewEntity("hitbox")
	hitbox := gfx.R(0, 0, o.Width, o.Height)
	em.Add(e, components.Pos{Vec: gfx.V(o.X, o.Y)})
	em.Add(e, components.NewHitbox(hitbox))
}

func Area(em *entity.Manager, o *tiled.Object) {
	e := em.NewEntity("area")
	hitbox := gfx.R(o.X, o.Y, o.X+o.Width, o.Y+o.Height)
	em.Add(e, components.Area{
		Rect: hitbox,
		Name: o.Name,
	})
}

func Player(em *entity.Manager, o *tiled.Object) string {
	e := em.NewEntity("player")
	hitbox := gfx.R(4, 8, 18, 22)

	em.Add(e, components.Pos{Vec: gfx.V(o.X, o.Y)})
	em.Add(e, components.Velocity{Vec: gfx.V(0, 0)})
	em.Add(e, components.Joystick{})
	if !headless {
		tmp, err := gfx.OpenPNG("assets/images/platformer.png")
		if err != nil {
			log.Fatal(err)
		}
		pImage, _ := ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)
		em.Add(e, components.Drawable{pImage.SubImage(image.Rect(5, 10, 27, 32)).(*ebiten.Image)})
	}
	em.Add(e, components.NewHitbox(hitbox))
	return e
}

func parseInAreaCondition(em *entity.Manager, prop *tiled.Property) []string {
	params := strings.Split(prop.Value, ",")

	// Find area matching name
	for _, e := range em.FilteredEntities(components.AreaType) {
		if em.Area(e).Name == params[1] {
			return []string{prop.Name, params[0], e}
		}
	}

	log.Fatal("no area exists with name", params[1])
	return nil
}

func Condition(em *entity.Manager, o *tiled.Object) {
	e := em.NewEntity("condition")
	cond := components.Condition{
		Name: o.Name,
	}
	for _, p := range o.Properties {
		switch p.Name {
		case "key_pressed":
			cond.Conditions = append(cond.Conditions, []string{p.Name, p.Value})
		case "in_area":
			cond.Conditions = append(cond.Conditions, parseInAreaCondition(em, p))
		default:
			fmt.Println("Unknown condition property", o)
		}
	}
	em.Add(e, cond)
}

func Text(em *entity.Manager, o *tiled.Object) {
	fmt.Println("Found a text!", o.Text)
	e := em.NewEntity("text")
	img, _ := ebiten.NewImage(100, 100, ebiten.FilterDefault)
	ebitenutil.DebugPrint(img, o.Text.Text)
	if !headless {
		em.Add(e, components.Pos{Vec: gfx.V(o.X, o.Y)})

		conditional := o.Properties.GetString("conditional")
		if conditional != "" {
			maxTransitions := 0 //Default, infinite number of transitions
			i, err := strconv.Atoi(o.Properties.GetString("max_transitions"))
			if err == nil {
				maxTransitions = i
				fmt.Println("loaded", maxTransitions, conditional)
			}

			em.Add(e, components.ConditionalDrawable{ConditionName: conditional, MaxTransitions: maxTransitions})
		}

	}
	em.Add(e, components.Drawable{img})
	// e := em.NewEntity("condition")
	// em.Add(e, cond)
}
