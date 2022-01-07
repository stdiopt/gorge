package input

import "github.com/stdiopt/gorge"

type keyManager struct {
	gorge    *gorge.Context
	keyState map[Key]ActionState
}

func (m *keyManager) update() {
	for k, v := range m.keyState {
		switch v {
		case ActionDown:
			m.keyState[k] = ActionHold
		case ActionUp:
			delete(m.keyState, k)
		}
	}
}

func (m *keyManager) SetKeyState(key Key, s ActionState) {
	if m.keyState == nil {
		m.keyState = map[Key]ActionState{}
	}
	m.keyState[key] = s
	switch s {
	case ActionUp:
		gorge.Trigger(m.gorge, EventKeyUp{key})
	case ActionDown:
		gorge.Trigger(m.gorge, EventKeyDown{key})
	}
}

func (m *keyManager) KeyUp(k Key) bool {
	return m.getKey(k) == ActionUp
}

// GetKey checks if a key was pressed
func (m *keyManager) KeyPress(k Key) bool {
	return m.getKey(k) == ActionUp
}

func (m *keyManager) KeyDown(k Key) bool {
	s := m.getKey(k)
	return s == ActionDown || s == ActionHold
}

func (m *keyManager) getKey(k Key) ActionState {
	if m.keyState == nil {
		return ActionState(0)
	}
	return m.keyState[k]
}

// Key represents a keyboard key.
type Key int

// Known input keys
const (
	KeyUnknown = Key(iota)
	KeySpace
	KeyApostrophe
	KeyComma
	KeyMinus
	KeyPeriod
	KeySlash
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeySemicolon
	KeyEqual
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyLeftBracket
	KeyBackslash
	KeyRightBracket
	KeyGraveAccent
	KeyWorld1
	KeyWorld2
	KeyEscape
	KeyEnter
	KeyTab
	KeyBackspace
	KeyInsert
	KeyDelete
	KeyArrowRight
	KeyArrowLeft
	KeyArrowDown
	KeyArrowUp
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyCapsLock
	KeyScrollLock
	KeyNumLock
	KeyPrintScreen
	KeyPause
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyF25
	KeyKP0
	KeyKP1
	KeyKP2
	KeyKP3
	KeyKP4
	KeyKP5
	KeyKP6
	KeyKP7
	KeyKP8
	KeyKP9
	KeyKPDecimal
	KeyKPDivide
	KeyKPMultiply
	KeyKPSubtract
	KeyKPAdd
	KeyKPEnter
	KeyKPEqual
	KeyLeftShift
	KeyLeftControl
	KeyLeftAlt
	KeyLeftSuper
	KeyRightShift
	KeyRightControl
	KeyRightAlt
	KeyRightSuper
	KeyMenu
	KeyLast
)

func (k Key) String() string {
	kstr := map[Key]string{
		KeyUnknown:      "KeyUnknown",
		KeySpace:        "KeySpace",
		KeyApostrophe:   "KeyApostrophe",
		KeyComma:        "KeyComma",
		KeyMinus:        "KeyMinus",
		KeyPeriod:       "KeyPeriod",
		KeySlash:        "KeySlash",
		Key0:            "Key0",
		Key1:            "Key1",
		Key2:            "Key2",
		Key3:            "Key3",
		Key4:            "Key4",
		Key5:            "Key5",
		Key6:            "Key6",
		Key7:            "Key7",
		Key8:            "Key8",
		Key9:            "Key9",
		KeySemicolon:    "KeySemicolon",
		KeyEqual:        "KeyEqual",
		KeyA:            "KeyA",
		KeyB:            "KeyB",
		KeyC:            "KeyC",
		KeyD:            "KeyD",
		KeyE:            "KeyE",
		KeyF:            "KeyF",
		KeyG:            "KeyG",
		KeyH:            "KeyH",
		KeyI:            "KeyI",
		KeyJ:            "KeyJ",
		KeyK:            "KeyK",
		KeyL:            "KeyL",
		KeyM:            "KeyM",
		KeyN:            "KeyN",
		KeyO:            "KeyO",
		KeyP:            "KeyP",
		KeyQ:            "KeyQ",
		KeyR:            "KeyR",
		KeyS:            "KeyS",
		KeyT:            "KeyT",
		KeyU:            "KeyU",
		KeyV:            "KeyV",
		KeyW:            "KeyW",
		KeyX:            "KeyX",
		KeyY:            "KeyY",
		KeyZ:            "KeyZ",
		KeyLeftBracket:  "KeyLeftBracket",
		KeyBackslash:    "KeyBackslash",
		KeyRightBracket: "KeyRightBracket",
		KeyGraveAccent:  "KeyGraveAccent",
		KeyWorld1:       "KeyWorld1",
		KeyWorld2:       "KeyWorld2",
		KeyEscape:       "KeyEscape",
		KeyEnter:        "KeyEnter",
		KeyTab:          "KeyTab",
		KeyBackspace:    "KeyBackspace",
		KeyInsert:       "KeyInsert",
		KeyDelete:       "KeyDelete",
		KeyArrowRight:   "KeyArrowRight",
		KeyArrowLeft:    "KeyArrowLeft",
		KeyArrowDown:    "KeyArrowDown",
		KeyArrowUp:      "KeyArrowUp",
		KeyPageUp:       "KeyPageUp",
		KeyPageDown:     "KeyPageDown",
		KeyHome:         "KeyHome",
		KeyEnd:          "KeyEnd",
		KeyCapsLock:     "KeyCapsLock",
		KeyScrollLock:   "KeyScrollLock",
		KeyNumLock:      "KeyNumLock",
		KeyPrintScreen:  "KeyPrintScreen",
		KeyPause:        "KeyPause",
		KeyF1:           "KeyF1",
		KeyF2:           "KeyF2",
		KeyF3:           "KeyF3",
		KeyF4:           "KeyF4",
		KeyF5:           "KeyF5",
		KeyF6:           "KeyF6",
		KeyF7:           "KeyF7",
		KeyF8:           "KeyF8",
		KeyF9:           "KeyF9",
		KeyF10:          "KeyF10",
		KeyF11:          "KeyF11",
		KeyF12:          "KeyF12",
		KeyF13:          "KeyF13",
		KeyF14:          "KeyF14",
		KeyF15:          "KeyF15",
		KeyF16:          "KeyF16",
		KeyF17:          "KeyF17",
		KeyF18:          "KeyF18",
		KeyF19:          "KeyF19",
		KeyF20:          "KeyF20",
		KeyF21:          "KeyF21",
		KeyF22:          "KeyF22",
		KeyF23:          "KeyF23",
		KeyF24:          "KeyF24",
		KeyF25:          "KeyF25",
		KeyKP0:          "KeyKP0",
		KeyKP1:          "KeyKP1",
		KeyKP2:          "KeyKP2",
		KeyKP3:          "KeyKP3",
		KeyKP4:          "KeyKP4",
		KeyKP5:          "KeyKP5",
		KeyKP6:          "KeyKP6",
		KeyKP7:          "KeyKP7",
		KeyKP8:          "KeyKP8",
		KeyKP9:          "KeyKP9",
		KeyKPDecimal:    "KeyKPDecimal",
		KeyKPDivide:     "KeyKPDivide",
		KeyKPMultiply:   "KeyKPMultiply",
		KeyKPSubtract:   "KeyKPSubtract",
		KeyKPAdd:        "KeyKPAdd",
		KeyKPEnter:      "KeyKPEnter",
		KeyKPEqual:      "KeyKPEqual",
		KeyLeftShift:    "KeyLeftShift",
		KeyLeftControl:  "KeyLeftControl",
		KeyLeftAlt:      "KeyLeftAlt",
		KeyLeftSuper:    "KeyLeftSuper",
		KeyRightShift:   "KeyRightShift",
		KeyRightControl: "KeyRightControl",
		KeyRightAlt:     "KeyRightAlt",
		KeyRightSuper:   "KeyRightSuper",
		KeyMenu:         "KeyMenu",
		KeyLast:         "KeyLast",
	}

	if str, ok := kstr[k]; ok {
		return str
	}
	return "<unknown key>"
}

// Char returns the default character for the specific key
// us layout based?
func (k Key) Char() string {
	kstr := map[Key]string{
		KeyUnknown:      "",
		KeySpace:        " ",
		KeyApostrophe:   "'",
		KeyComma:        ",",
		KeyMinus:        "-",
		KeyPeriod:       ".",
		KeySlash:        "/",
		Key0:            "0",
		Key1:            "1",
		Key2:            "2",
		Key3:            "3",
		Key4:            "4",
		Key5:            "5",
		Key6:            "6",
		Key7:            "7",
		Key8:            "8",
		Key9:            "9",
		KeySemicolon:    ";",
		KeyEqual:        "=",
		KeyA:            "a",
		KeyB:            "b",
		KeyC:            "c",
		KeyD:            "d",
		KeyE:            "e",
		KeyF:            "f",
		KeyG:            "g",
		KeyH:            "h",
		KeyI:            "i",
		KeyJ:            "j",
		KeyK:            "k",
		KeyL:            "l",
		KeyM:            "m",
		KeyN:            "n",
		KeyO:            "o",
		KeyP:            "p",
		KeyQ:            "q",
		KeyR:            "r",
		KeyS:            "s",
		KeyT:            "t",
		KeyU:            "u",
		KeyV:            "v",
		KeyW:            "w",
		KeyX:            "x",
		KeyY:            "y",
		KeyZ:            "z",
		KeyLeftBracket:  "[",
		KeyBackslash:    "\\",
		KeyRightBracket: "]",
		KeyGraveAccent:  "`",
		KeyEnter:        "\n",
		KeyTab:          "\t",
		KeyKPDecimal:    ",",
		KeyKPDivide:     "/",
		KeyKPMultiply:   "*",
		KeyKPSubtract:   "-",
		KeyKPAdd:        "+",
		KeyKPEnter:      "\n",
		KeyKPEqual:      "=",
	}

	if str, ok := kstr[k]; ok {
		return str
	}
	return ""
}
