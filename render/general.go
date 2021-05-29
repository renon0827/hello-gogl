package render

import (
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/sirupsen/logrus"
)

var window *glfw.Window

func createWindow() error {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var err error
	window, err = glfw.CreateWindow(800, 600, "unko", nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()
	return nil
}

func Mainloop(mainloop func()) {
	defer glfw.Terminate()
	for !window.ShouldClose() {
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LESS)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		mainloop()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func Initialize() error {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		return err
	}

	if err := createWindow(); err != nil {
		return err
	}

	if err := gl.Init(); err != nil {
		return err
	}

	logrus.Infoln("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))

	return nil
}
