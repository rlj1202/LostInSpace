package lostinspace

import (
	"image"
	"log"

	"github.com/go-gl/gl/v4.1-compatibility/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	nativeWindow *glfw.Window
}

func init() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
}

func NewWindow(width, height int, title string, icons []image.Image, resizable bool) *Window {
	if resizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.SetIcon(icons)

	window.SetKeyCallback(keyInput)
	window.SetMouseButtonCallback(mouseInput)
	window.SetScrollCallback(scrollInput)
	window.SetCursorEnterCallback(cursorEnterInput)
	window.SetCursorPosCallback(cursorPosInput)
	window.SetSizeCallback(sizeInput)
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("Opengl version", version)

	gl.ClearColor(0, 0, 0.1, 1)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.LINE_SMOOTH)
	gl.Enable(gl.MULTISAMPLE)
	gl.Hint(gl.LINE_SMOOTH_HINT, gl.NICEST)

	newWindow := &Window{
		nativeWindow: window,
	}

	return newWindow
}

func (window *Window) ShouldClose() bool {
	return window.nativeWindow.ShouldClose()
}

func (window *Window) Update() {
	window.nativeWindow.SwapBuffers()
}

func (window *Window) GetTitle() string {
	return window.GetTitle()
}

func (window *Window) SetTitle(title string) {
	window.nativeWindow.SetTitle(title)
}

func (window *Window) GetSize() (int, int) {
	return window.nativeWindow.GetSize()
}

func keyInput(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	PushEvent(KeyboardEvent{
		Key:      Key(key),
		Scancode: scancode,
		Action:   Action(action),
		Mods:     ModifierKey(mods),
	})
}

func mouseInput(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	xoff, yoff := w.GetCursorPos()
	PushEvent(MouseEvent{
		XPos:   xoff,
		YPos:   yoff,
		Button: MouseButton(button),
		Action: Action(action),
		Mod:    ModifierKey(mod),
	})
}

func scrollInput(w *glfw.Window, xoff, yoff float64) {
	PushEvent(ScrollEvent{
		XOff: xoff,
		YOff: yoff,
	})
}

func cursorEnterInput(w *glfw.Window, entered bool) {
	PushEvent(CursorEnterEvent{
		Entered: entered,
	})
}

func cursorPosInput(w *glfw.Window, xpos, ypos float64) {
	PushEvent(CursorPosEvent{
		XPos: xpos,
		YPos: ypos,
	})
}

func sizeInput(w *glfw.Window, width, height int) {
	PushEvent(WindowSizeEvent{
		Width:  width,
		Height: height,
	})
}
