// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

// This code is mostly based on the glfw samples from github.com/jteeuwen/glfw
// NEHE Tutorial 05: 3D shapes.
// http://nehe.gamedev.net/data/lessons/lesson.asp?lesson=05
package main

import (
	"flag"
	"fmt"
	"github.com/andrebq/wfobj"
	"github.com/banthar/gl"
	"github.com/banthar/glu"
	"github.com/jteeuwen/glfw"
	"math/rand"
	"os"
)

const (
	Title  = "Nehe 05"
	Width  = 640
	Height = 480
)

var (
	trisAngle float32
	quadAngle float32
	running   bool
	mesh      *wfobj.Mesh
	faceColor         = map[int][]float32{}
	speed     float32 = 0.5
)

type State struct {
	Keys             map[int]bool
	yRot, xRot, zRot float32
	speed            float32
	light            [4]float32
	wheel            float32
}

var (
	globalState = &State{Keys: make(map[int]bool), speed: 1}
)

func main() {
	flag.Parse()
	var err error
	if err = glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

	if len(flag.Args()) == 0 {
		old := flag.Usage
		flag.Usage = func() {
			old()
			fmt.Fprintf(os.Stderr, "You MUST pass the name of the file to view\n")
		}
		flag.Usage()
		return
	}

	mesh, err = wfobj.LoadMeshFromFile(flag.Args()[0])
	if err != nil {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

	defer glfw.Terminate()

	if err = glfw.OpenWindow(Width, Height, 8, 8, 8, 8, 0, 8, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle(Title)
	glfw.SetWindowSizeCallback(onResize)
	glfw.SetKeyCallback(onKey)
	glfw.SetCharCallback(onChar)
	glfw.SetMouseWheelCallback(onWheel)

	initGL()

	running = true
	for running && glfw.WindowParam(glfw.Opened) == 1 {
		handleInput()
		drawScene()
	}
}

func onResize(w, h int) {
	if h == 0 {
		h = 1
	}

	gl.Viewport(0, 0, w, h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	glu.Perspective(45.0, float64(w)/float64(h), 0.1, 100.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func onKey(key, state int) {
	switch key {
	case glfw.KeyEsc:
		running = false
	case glfw.KeyLeft, glfw.KeyRight, glfw.KeyUp, glfw.KeyDown, KeyW, KeyS, KeyPlus, KeyMinus, glfw.KeyPagedown, glfw.KeyPageup:
		globalState.Keys[key] = state == glfw.KeyPress
	default:
		print("Key: ", key, " is pressed? ", state == glfw.KeyPress)
	}
}

func onWheel(delta int) {
	globalState.wheel = float32(delta)
}

func onChar(key, state int) {
}

func initGL() {
	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)

	globalState.light[0] = 0
	globalState.light[1] = 20
	globalState.light[2] = -10
	globalState.light[3] = 1

	gl.Lightfv(gl.LIGHT1, gl.AMBIENT, []float32{1, 1, 1})
	gl.Lightfv(gl.LIGHT1, gl.DIFFUSE, []float32{1, 1, 1})
	gl.Lightfv(gl.LIGHT1, gl.POSITION, globalState.light[:])
	gl.Enable(gl.LIGHT1)

	gl.Enable(gl.LIGHTING)
}

func handleInput() {
	if globalState.Keys[glfw.KeyLeft] {
		globalState.yRot -= speed
	}
	if globalState.Keys[glfw.KeyRight] {
		globalState.yRot += speed
	}
	if globalState.Keys[glfw.KeyUp] {
		globalState.xRot += speed
	}
	if globalState.Keys[glfw.KeyDown] {
		globalState.xRot -= speed
	}
	if globalState.Keys[KeyW] {
		globalState.zRot += speed
	}
	if globalState.Keys[KeyS] {
		globalState.zRot -= speed
	}
	if globalState.Keys[KeyPlus] {
		speed += 0.5
	}
	if globalState.Keys[KeyMinus] {
		speed -= 0.5
	}
	if globalState.Keys[glfw.KeyPageup] {
		globalState.speed += 0.5
	}
	if globalState.Keys[glfw.KeyPagedown] {
		globalState.speed -= 0.5
	}
	if globalState.speed < 1 {
		globalState.speed = 1
	}
	if speed < 0.5 {
		speed = 0.5
	} else if speed > 2 {
		speed = 2
	}
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.LoadIdentity()

	gl.Translatef(0, 0, -20+globalState.wheel*globalState.speed)

	gl.Rotatef(globalState.xRot, 1, 0, 0)
	gl.Rotatef(globalState.yRot, 0, 1, 0)
	gl.Rotatef(globalState.zRot, 0, 0, 1)

	if globalState.speed != 1 {
		gl.Scalef(globalState.speed, globalState.speed, globalState.speed)
	}

	gl.Begin(gl.QUADS)
	for i, _ := range mesh.Faces {
		if colors, ok := faceColor[i]; ok {
			gl.Color3f(colors[0], colors[1], colors[2])
		} else {
			faceColor[i] = make([]float32, 3)
			faceColor[i][0] = rand.Float32()
			faceColor[i][1] = rand.Float32()
			faceColor[i][2] = rand.Float32()
			gl.Color3f(faceColor[i][0], faceColor[i][1], faceColor[i][2])
		}

		face := &mesh.Faces[i]
		for j, _ := range face.Vertices {
			var v *wfobj.Vertex
			if len(face.Normals) > 0 {
				v = &face.Normals[j]
				gl.Normal3f(v.X, v.Y, v.Z)
			}
			v = &face.Vertices[j]
			gl.Vertex3f(v.X, v.Y, v.Z)
		}
	}
	gl.End()

	glfw.SwapBuffers()
}
