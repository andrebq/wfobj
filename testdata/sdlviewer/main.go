// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

// This code is mosdly based on the glfw samples from github.com/jteeuwen/glfw
// NEHE Tutorial 05: 3D shapes.
// http://nehe.gamedev.net/data/lessons/lesson.asp?lesson=05
package main

import (
	"flag"
	"fmt"
	"github.com/andrebq/wfobj"
	"github.com/banthar/gl"
	"github.com/banthar/glu"
	"github.com/banthar/Go-SDL/sdl"
	"math/rand"
	"os"
	"log"
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

type DragInfo struct {
	Start  wfobj.Vertex
	End    wfobj.Vertex
	IsDrag bool
}

type State struct {
	Keys     map[int]bool
	Rot      wfobj.Vertex
	speed    float32
	light    [4]float32
	wheel    float32
	Mouse    map[int]bool
	MousePos wfobj.Vertex
	Drag     DragInfo
}

var (
	globalState = &State{Keys: make(map[int]bool), speed: 1}
)

func main() {
	globalState.Mouse = make(map[int]bool)
	flag.Parse()
	var err error
	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		log.Printf("Unable to init SDL")
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

	defer sdl.Quit()
	
	screen := sdl.SetVideoMode(1024, 600, 32, sdl.OPENGL|sdl.RESIZABLE)
	if screen == nil {
		log.Printf("Unable to init sdl Screen")
	}

	sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1)

	initGL()
	
	sdl.WM_SetCaption("sdlviewer", "sdlviewer")

	running = true
	for running {
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

/*
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
	globalState.MousePos.Z = float32(delta)
}

func onMouseButton(button, state int) {
	globalState.Mouse[button] = state == glfw.KeyPress
	if globalState.Mouse[glfw.Mouse1] {
		x, y := glfw.MousePos()
		globalState.Drag.Start.X, globalState.Drag.Start.Y = float32(x), float32(y)
		globalState.Drag.IsDrag = true
	} else {
		if globalState.Drag.IsDrag {
			globalState.Drag.IsDrag = false
			x, y := glfw.MousePos()
			globalState.Drag.End.X, globalState.Drag.End.Y = float32(x), float32(y)
		}
	}
}

func onMousePos(x, y int) {
	globalState.MousePos.X = float32(x)
	globalState.MousePos.Y = float32(y)
	if globalState.Drag.IsDrag {
		//		println("Start: ", fmt.Sprintf("%v", globalState.Drag.Start))
		//		println("End: ", fmt.Sprintf("%v", globalState.Drag.End))
		//		println("Sub: ", fmt.Sprintf("%v", globalState.Drag.End.Sub(&globalState.Drag.Start)))
		globalState.Drag.End.X = float32(x)
		globalState.Drag.End.Y = float32(y)
	}
}

func onChar(key, state int) {
}
*/

func initGL() {
	gl.Init()
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
	for e := sdl.PollEvent(); e != nil; {
		// just ignore for the moment
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

	gl.Translatef(0, 0, -20+globalState.MousePos.Z*globalState.speed)

	gl.Rotatef(globalState.Rot.X, 1, 0, 0)
	gl.Rotatef(globalState.Rot.Y, 0, 1, 0)
	gl.Rotatef(globalState.Rot.Z, 0, 0, 1)

	if globalState.speed != 1 {
		gl.Scalef(globalState.speed, globalState.speed, globalState.speed)
	}
	
	gl.RenderMode(gl.RENDER)

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
	gl.Finish()
	gl.Flush()
	
	sdl.GL_SwapBuffers()
}
