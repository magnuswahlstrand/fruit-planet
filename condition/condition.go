package condition

import (
	"log"
	"strings"
	"time"

	"github.com/kyeett/gomponents/components"

	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/ecs/entity"

	"github.com/hajimehoshi/ebiten"
)

type KeyJustPressed struct {
	ebiten.Key
}

func (kp KeyJustPressed) IsMet() bool {
	return inpututil.IsKeyJustPressed(kp.Key)
}

type KeyPressed struct {
	ebiten.Key
}

func (kp KeyPressed) IsMet() bool {
	return ebiten.IsKeyPressed(kp.Key)
}

type InArea struct {
	entityID, areaID string
	em               *entity.Manager
}

func NewInArea(em *entity.Manager, entityID, areaID string) InArea {
	return InArea{
		entityID: entityID,
		areaID:   areaID,
		em:       em,
	}
}

func (c InArea) IsMet() bool {
	pos, err := c.em.Get(c.entityID, components.PosType)
	if err != nil {
		return false
	}
	hitbox, err := c.em.Get(c.entityID, components.HitboxType)
	if err != nil {
		return false
	}

	r1 := hitbox.(*components.Hitbox).Rect.Moved(pos.(*components.Pos).Vec)
	r2 := c.em.Area(c.areaID)
	return r1.Overlaps(r2.Rect)
}

type WaitUntil struct {
	time.Time
}

func (c WaitUntil) IsMet() bool {
	return time.Now().After(c.Time)
}

type AnimationComplete struct {
	animationID string
	em          *entity.Manager
}

func NewAnimationComplete(em *entity.Manager, animationID string) AnimationComplete {
	return AnimationComplete{
		animationID: animationID,
		em:          em,
	}
}

func (c AnimationComplete) IsMet() bool {
	for _, e := range c.em.FilteredEntities(components.AnimationType) {
		if e == c.animationID {
			return false
		}
	}
	// No such animation, must have ended
	return true
}

func KeyNameToKey(name string) ebiten.Key {
	switch strings.ToLower(name) {
	case "0":
		return ebiten.Key0
	case "1":
		return ebiten.Key1
	case "2":
		return ebiten.Key2
	case "3":
		return ebiten.Key3
	case "4":
		return ebiten.Key4
	case "5":
		return ebiten.Key5
	case "6":
		return ebiten.Key6
	case "7":
		return ebiten.Key7
	case "8":
		return ebiten.Key8
	case "9":
		return ebiten.Key9
	case "a":
		return ebiten.KeyA
	case "b":
		return ebiten.KeyB
	case "c":
		return ebiten.KeyC
	case "d":
		return ebiten.KeyD
	case "e":
		return ebiten.KeyE
	case "f":
		return ebiten.KeyF
	case "g":
		return ebiten.KeyG
	case "h":
		return ebiten.KeyH
	case "i":
		return ebiten.KeyI
	case "j":
		return ebiten.KeyJ
	case "k":
		return ebiten.KeyK
	case "l":
		return ebiten.KeyL
	case "m":
		return ebiten.KeyM
	case "n":
		return ebiten.KeyN
	case "o":
		return ebiten.KeyO
	case "p":
		return ebiten.KeyP
	case "q":
		return ebiten.KeyQ
	case "r":
		return ebiten.KeyR
	case "s":
		return ebiten.KeyS
	case "t":
		return ebiten.KeyT
	case "u":
		return ebiten.KeyU
	case "v":
		return ebiten.KeyV
	case "w":
		return ebiten.KeyW
	case "x":
		return ebiten.KeyX
	case "y":
		return ebiten.KeyY
	case "z":
		return ebiten.KeyZ
	case "alt":
		return ebiten.KeyAlt
	case "apostrophe":
		return ebiten.KeyApostrophe
	case "backslash":
		return ebiten.KeyBackslash
	case "backspace":
		return ebiten.KeyBackspace
	case "capslock":
		return ebiten.KeyCapsLock
	case "comma":
		return ebiten.KeyComma
	case "control":
		return ebiten.KeyControl
	case "delete":
		return ebiten.KeyDelete
	case "down":
		return ebiten.KeyDown
	case "end":
		return ebiten.KeyEnd
	case "enter":
		return ebiten.KeyEnter
	case "equal":
		return ebiten.KeyEqual
	case "escape":
		return ebiten.KeyEscape
	case "f1":
		return ebiten.KeyF1
	case "f2":
		return ebiten.KeyF2
	case "f3":
		return ebiten.KeyF3
	case "f4":
		return ebiten.KeyF4
	case "f5":
		return ebiten.KeyF5
	case "f6":
		return ebiten.KeyF6
	case "f7":
		return ebiten.KeyF7
	case "f8":
		return ebiten.KeyF8
	case "f9":
		return ebiten.KeyF9
	case "f10":
		return ebiten.KeyF10
	case "f11":
		return ebiten.KeyF11
	case "f12":
		return ebiten.KeyF12
	case "graveaccent":
		return ebiten.KeyGraveAccent
	case "home":
		return ebiten.KeyHome
	case "insert":
		return ebiten.KeyInsert
	case "kp0":
		return ebiten.KeyKP0
	case "kp1":
		return ebiten.KeyKP1
	case "kp2":
		return ebiten.KeyKP2
	case "kp3":
		return ebiten.KeyKP3
	case "kp4":
		return ebiten.KeyKP4
	case "kp5":
		return ebiten.KeyKP5
	case "kp6":
		return ebiten.KeyKP6
	case "kp7":
		return ebiten.KeyKP7
	case "kp8":
		return ebiten.KeyKP8
	case "kp9":
		return ebiten.KeyKP9
	case "kpadd":
		return ebiten.KeyKPAdd
	case "kpdecimal":
		return ebiten.KeyKPDecimal
	case "kpdivide":
		return ebiten.KeyKPDivide
	case "kpenter":
		return ebiten.KeyKPEnter
	case "kpequal":
		return ebiten.KeyKPEqual
	case "kpmultiply":
		return ebiten.KeyKPMultiply
	case "kpsubtract":
		return ebiten.KeyKPSubtract
	case "left":
		return ebiten.KeyLeft
	case "leftbracket":
		return ebiten.KeyLeftBracket
	case "menu":
		return ebiten.KeyMenu
	case "minus":
		return ebiten.KeyMinus
	case "numlock":
		return ebiten.KeyNumLock
	case "pagedown":
		return ebiten.KeyPageDown
	case "pageup":
		return ebiten.KeyPageUp
	case "pause":
		return ebiten.KeyPause
	case "period":
		return ebiten.KeyPeriod
	case "printscreen":
		return ebiten.KeyPrintScreen
	case "right":
		return ebiten.KeyRight
	case "rightbracket":
		return ebiten.KeyRightBracket
	case "scrolllock":
		return ebiten.KeyScrollLock
	case "semicolon":
		return ebiten.KeySemicolon
	case "shift":
		return ebiten.KeyShift
	case "slash":
		return ebiten.KeySlash
	case "space":
		return ebiten.KeySpace
	case "tab":
		return ebiten.KeyTab
	case "up":
		return ebiten.KeyUp
	}

	log.Fatal("invalid key")
	return 0
}
