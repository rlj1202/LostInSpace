package lostinspace

import "github.com/go-gl/glfw/v3.2/glfw"

const (
	KEY_UNKNOWN    Key = Key(glfw.KeyUnknown)
	KEY_SPACE          = Key(glfw.KeySpace)
	KEY_APOSTROPHE     = Key(glfw.KeyApostrophe)
	KEY_COMMA          = Key(glfw.KeyComma)
	KEY_MINUS          = Key(glfw.KeyMinus)
	KEY_PERIOD         = Key(glfw.KeyPeriod)
	KEY_SLASH          = Key(glfw.KeySlash)

	KEY_0 = Key(glfw.Key0)
	KEY_1 = Key(glfw.Key1)
	KEY_2 = Key(glfw.Key2)
	KEY_3 = Key(glfw.Key3)
	KEY_4 = Key(glfw.Key4)
	KEY_5 = Key(glfw.Key5)
	KEY_6 = Key(glfw.Key6)
	KEY_7 = Key(glfw.Key7)
	KEY_8 = Key(glfw.Key8)
	KEY_9 = Key(glfw.Key9)

	KEY_SEMICOLON = Key(glfw.KeySemicolon)
	KEY_EQUAL     = Key(glfw.KeyEqual)

	KEY_A = Key(glfw.KeyA)
	KEY_B = Key(glfw.KeyB)
	KEY_C = Key(glfw.KeyC)
	KEY_D = Key(glfw.KeyD)
	KEY_E = Key(glfw.KeyE)
	KEY_F = Key(glfw.KeyF)
	KEY_G = Key(glfw.KeyG)
	KEY_H = Key(glfw.KeyH)
	KEY_I = Key(glfw.KeyI)
	KEY_J = Key(glfw.KeyJ)
	KEY_K = Key(glfw.KeyK)
	KEY_L = Key(glfw.KeyL)
	KEY_M = Key(glfw.KeyM)
	KEY_N = Key(glfw.KeyN)
	KEY_O = Key(glfw.KeyO)
	KEY_P = Key(glfw.KeyP)
	KEY_Q = Key(glfw.KeyQ)
	KEY_R = Key(glfw.KeyR)
	KEY_S = Key(glfw.KeyS)
	KEY_T = Key(glfw.KeyT)
	KEY_U = Key(glfw.KeyU)
	KEY_V = Key(glfw.KeyV)
	KEY_W = Key(glfw.KeyW)
	KEY_X = Key(glfw.KeyX)
	KEY_Y = Key(glfw.KeyY)
	KEY_Z = Key(glfw.KeyZ)

	KEY_BRACKET_LEFT  = Key(glfw.KeyLeftBracket)
	KEY_BRACKET_RIGHT = Key(glfw.KeyRightBracket)
	KEY_BACK_SLASH    = Key(glfw.KeyBackslash)
	KEY_GRAVE_ACCENT  = Key(glfw.KeyGraveAccent)
	KEY_WORLD_1       = Key(glfw.KeyWorld1)
	KEY_WORLD_2       = Key(glfw.KeyWorld2)
	KEY_ESCAPE        = Key(glfw.KeyEscape)
	KEY_ENTER         = Key(glfw.KeyEnter)
	KEY_TAB           = Key(glfw.KeyTab)
	KEY_BACKSPACE     = Key(glfw.KeyBackspace)
	KEY_INSERT        = Key(glfw.KeyInsert)
	KEY_DELETE        = Key(glfw.KeyDelete)
	KEY_RIGHT         = Key(glfw.KeyRight)
	KEY_LEFT          = Key(glfw.KeyLeft)
	KEY_DOWN          = Key(glfw.KeyDown)
	KEY_UP            = Key(glfw.KeyUp)
	KEY_PAGE_UP       = Key(glfw.KeyPageUp)
	KEY_PAGE_DOWN     = Key(glfw.KeyPageDown)
	KEY_HOME          = Key(glfw.KeyHome)
	KEY_END           = Key(glfw.KeyEnd)
	KEY_CAPS_LOCK     = Key(glfw.KeyCapsLock)
	KEY_SCROLL_LOCK   = Key(glfw.KeyScrollLock)
	KEY_NUM_LOCK      = Key(glfw.KeyNumLock)
	KEY_PRINT_SCREEN  = Key(glfw.KeyPrintScreen)
	KEY_PAUSE         = Key(glfw.KeyPause)

	KEY_F1  = Key(glfw.KeyF1)
	KEY_F2  = Key(glfw.KeyF2)
	KEY_F3  = Key(glfw.KeyF3)
	KEY_F4  = Key(glfw.KeyF4)
	KEY_F5  = Key(glfw.KeyF5)
	KEY_F6  = Key(glfw.KeyF6)
	KEY_F7  = Key(glfw.KeyF7)
	KEY_F8  = Key(glfw.KeyF8)
	KEY_F9  = Key(glfw.KeyF9)
	KEY_F10 = Key(glfw.KeyF10)
	KEY_F11 = Key(glfw.KeyF11)
	KEY_F12 = Key(glfw.KeyF12)
	KEY_F13 = Key(glfw.KeyF13)
	KEY_F14 = Key(glfw.KeyF14)
	KEY_F15 = Key(glfw.KeyF15)
	KEY_F16 = Key(glfw.KeyF16)
	KEY_F17 = Key(glfw.KeyF17)
	KEY_F18 = Key(glfw.KeyF18)
	KEY_F19 = Key(glfw.KeyF19)
	KEY_F20 = Key(glfw.KeyF20)
	KEY_F21 = Key(glfw.KeyF21)
	KEY_F22 = Key(glfw.KeyF22)
	KEY_F23 = Key(glfw.KeyF23)
	KEY_F24 = Key(glfw.KeyF24)
	KEY_F25 = Key(glfw.KeyF25)

	KEY_KP_0     = Key(glfw.KeyKP0)
	KEY_KP_1     = Key(glfw.KeyKP1)
	KEY_KP_2     = Key(glfw.KeyKP2)
	KEY_KP_3     = Key(glfw.KeyKP3)
	KEY_KP_4     = Key(glfw.KeyKP4)
	KEY_KP_5     = Key(glfw.KeyKP5)
	KEY_KP_6     = Key(glfw.KeyKP6)
	KEY_KP_7     = Key(glfw.KeyKP7)
	KEY_KP_8     = Key(glfw.KeyKP8)
	KEY_KP_9     = Key(glfw.KeyKP9)
	KEY_KP_DEC   = Key(glfw.KeyKPDecimal)
	KEY_KP_DIV   = Key(glfw.KeyKPDivide)
	KEY_KP_MUL   = Key(glfw.KeyKPMultiply)
	KEY_KP_SUB   = Key(glfw.KeyKPSubtract)
	KEY_KP_ADD   = Key(glfw.KeyKPAdd)
	KEY_KP_ENTER = Key(glfw.KeyKPEnter)
	KEY_KP_EQUAL = Key(glfw.KeyKPEqual)

	KEY_SHIFT_LEFT    = Key(glfw.KeyLeftShift)
	KEY_SHIFT_RIGHT   = Key(glfw.KeyRightShift)
	KEY_CONTROL_LEFT  = Key(glfw.KeyLeftControl)
	KEY_CONTROL_RIGHT = Key(glfw.KeyRightControl)
	KEY_ALT_LEFT      = Key(glfw.KeyLeftAlt)
	KEY_ALT_RIGHT     = Key(glfw.KeyRightAlt)
	KEY_SUPER_LEFT    = Key(glfw.KeyLeftSuper)
	KEY_SUPER_RIGHT   = Key(glfw.KeyRightSuper)
	KEY_MENU          = Key(glfw.KeyMenu)
	KEY_LAST          = Key(glfw.KeyLast)
)

const (
	ACTION_PRESS   Action = Action(glfw.Press)
	ACTION_RELEASE        = Action(glfw.Release)
	ACTION_REPEAT         = Action(glfw.Repeat)
)

const (
	MOUSE_BUTTON_1      MouseButton = MouseButton(glfw.MouseButton1)
	MOUSE_BUTTON_2                  = MouseButton(glfw.MouseButton2)
	MOUSE_BUTTON_3                  = MouseButton(glfw.MouseButton3)
	MOUSE_BUTTON_4                  = MouseButton(glfw.MouseButton4)
	MOUSE_BUTTON_5                  = MouseButton(glfw.MouseButton5)
	MOUSE_BUTTON_6                  = MouseButton(glfw.MouseButton6)
	MOUSE_BUTTON_7                  = MouseButton(glfw.MouseButton7)
	MOUSE_BUTTON_8                  = MouseButton(glfw.MouseButton8)
	MOUSE_BUTTON_LAST               = MouseButton(glfw.MouseButtonLast)
	MOUSE_BUTTON_LEFT               = MouseButton(glfw.MouseButtonLeft)
	MOUSE_BUTTON_RIGHT              = MouseButton(glfw.MouseButtonRight)
	MOUSE_BUTTON_MIDDLE             = MouseButton(glfw.MouseButtonMiddle)
)

const (
	MOD_ALT     ModifierKey = ModifierKey(glfw.ModAlt)
	MOD_SHIFT               = ModifierKey(glfw.ModShift)
	MOD_SUPER               = ModifierKey(glfw.ModSuper)
	MOD_CONTROL             = ModifierKey(glfw.ModControl)
)

type Key glfw.Key
type Action glfw.Action
type MouseButton glfw.MouseButton
type ModifierKey glfw.ModifierKey

type inputListenerImpl struct{}

var inputListener inputListenerImpl
var keyActionStates map[Key]Action

func init() {
	keyActionStates = make(map[Key]Action)
	inputListener = inputListenerImpl{}
	RegisterEventListener(inputListener)
}

func (listener inputListenerImpl) OnEvent(event Event) {
	switch event.(type) {
	case KeyboardEvent:
		keyboardEvent := event.(KeyboardEvent)
		action := keyboardEvent.Action
		key := keyboardEvent.Key
		keyActionStates[key] = action
	}
}

func GetKeyActionState(key Key) Action {
	action, exist := keyActionStates[key]
	if !exist {
		action = ACTION_RELEASE
	}

	return action
}
