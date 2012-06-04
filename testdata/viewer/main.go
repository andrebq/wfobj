// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

// This code is mostly based on the glfw samples from github.com/jteeuwen/glfw
// NEHE Tutorial 05: 3D shapes.
// http://nehe.gamedev.net/data/lessons/lesson.asp?lesson=05
package main

import (
	"fmt"
	"github.com/jteeuwen/glfw"
	"github.com/banthar/gl"
	"github.com/banthar/glu"
	"github.com/andrebq/wfobj"
	"math/rand"
	"os"
	"flag"
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
	mesh *wfobj.Mesh
	faceColor = map[int][]float32{}
)

type State struct {
	Left bool
	Right bool
}

var (
	globalState = &State{}
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

	initGL()

	running = true
	for running && glfw.WindowParam(glfw.Opened) == 1 {
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
	case glfw.KeyLeft:
		globalState.Left = state == glfw.KeyPress
	case glfw.KeyRight:
		globalState.Right = state == glfw.KeyPress
	}
}

func initGL() {
	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.LoadIdentity()
	gl.Translatef(0, 0, -20)
	gl.Rotatef(quadAngle, 1, 1, 1)

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
			v := &face.Vertices[j]
			gl.Vertex3f(v.X, v.Y, v.Z)
		}
	}
	gl.End()
	
	if globalState.Left {
		quadAngle += 0.20
	}
	
	if globalState.Right {
		quadAngle -= 0.20
	}

	glfw.SwapBuffers()
}
