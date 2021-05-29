package main

import (
	"fmt"
	"hello-gogl/render"
	"hello-gogl/resources/format"
	"io/ioutil"
	"math"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil
}

func main() {
	// initialize
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetOutput(colorable.NewColorableStdout())

	err := render.Initialize()
	if err != nil {
		logrus.Panicln(err)
	}

	// load
	file, err := os.Open("sample.obj")
	if err != nil {
		logrus.Panicln(err)
	}
	obj, err := format.LoadOBJ(file)
	if err != nil {
		logrus.Panicln(err)
	}

	triangle := []float32{}
	index := []uint32{}
	for _, group := range obj.Groups {
		for _, v := range group.Vertex {
			triangle = append(triangle, []float32{float32(v[0]), float32(v[1]), float32(v[2])}...)
		}
		for _, p := range group.Polygon {
			tmp := []uint32{uint32(p[0].VertexIndex - 1), uint32(p[1].VertexIndex - 1), uint32(p[2].VertexIndex - 1)}
			index = append(index, tmp...)
		}
	}
	logrus.Infoln(len(triangle) / 3)
	logrus.Infoln(len(index) / 3)

	// shaders
	vertexShaderSource, err := ioutil.ReadFile("vertex.glsl")
	if err != nil {
		logrus.Panicln(err)
	}

	fragmentShaderSource, err := ioutil.ReadFile("pixel.glsl")
	if err != nil {
		logrus.Panicln(err)
	}

	vertexShader, err := compileShader(string(vertexShaderSource)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		logrus.Panicln(err)
	}
	fragmentShader, err := compileShader(string(fragmentShaderSource)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		logrus.Panicln(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	gl.UseProgram(prog)

	// matrix
	projectionLocation := gl.GetUniformLocation(prog, gl.Str("projection\x00"))
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), 800/600, 0.1, 10.0)
	gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])

	camera := mgl32.LookAtV(
		mgl32.Vec3{2, 3, 6},
		mgl32.Vec3{0, 2, 0},
		mgl32.Vec3{0, 1, 0},
	)
	cameraLocation := gl.GetUniformLocation(prog, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraLocation, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelLocation := gl.GetUniformLocation(prog, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelLocation, 1, false, &model[0])

	// buffers
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var triangleVBO uint32
	gl.GenBuffers(1, &triangleVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, triangleVBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(triangle), gl.Ptr(triangle), gl.STATIC_DRAW)

	var indexBuffer uint32
	gl.GenBuffers(1, &indexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(index), gl.Ptr(index), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(prog, gl.Str("vp\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 12, gl.PtrOffset(0))

	angleX := 0.0
	angleY := 0.0
	angleZ := 0.0
	period := 2.0
	previousTime := glfw.GetTime()

	render.Mainloop(func() {
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angleY += math.Sin((elapsed / period) / 6.0 * math.Pi * 2.0)
		model = mgl32.HomogRotate3DY(float32(angleY)).Mul4(mgl32.HomogRotate3DX(float32(angleX))).Mul4(mgl32.HomogRotate3DZ(float32(angleZ)))
		gl.UseProgram(prog)
		gl.UniformMatrix4fv(modelLocation, 1, false, &model[0])

		gl.DrawElements(gl.TRIANGLES, int32(len(index)), gl.UNSIGNED_INT, nil)
	})
}
