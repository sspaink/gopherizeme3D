package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/divan/three"
	"github.com/gopherjs/gopherjs/js"
	"github.com/hschendel/stl"
)

func main() {
	listener := func(e *js.Object) {
		dropdown := js.Global.Get("document").Call("querySelector", ".dropdown")
		click := func(event *js.Object) {
			event.Call("stopPropagation")
			dropdown.Get("classList").Call("toggle", "is-active")
		}
		dropdown.Call("addEventListener", "click", click)
	}

	js.Global.Get("document").Call("addEventListener", "DOMContentLoaded", listener)

	width := js.Global.Get("innerWidth").Float() - 300
	height := js.Global.Get("innerHeight").Float() - 50

	renderer := three.NewWebGLRenderer()
	renderer.SetSize(width, height, true)

	js.Global.Get("document").Call("getElementById", "scene").Call("appendChild", renderer.Get("domElement"))

	// setup camera and scene
	camera := three.NewPerspectiveCamera(70, width/height, 1, 500)
	camera.Position.Set(0, 0, 100)

	scene := three.NewScene()

	// lights
	light := three.NewDirectionalLight(three.NewColor("white"), 1)
	light.Position.Set(0, 256, 256)
	scene.Add(light)

	// material
	params := three.NewMaterialParameters()
	params.Color = three.NewColor("skyblue")
	mat := three.NewMeshLambertMaterial(params)

	geom := three.NewBufferGeometry()

	resp, err := http.Get("http://localhost:9000/GopherPrintable.stl")
	if err != nil {
		log.Panicf("Oh no! %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Oh no! %s", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Oh no! %s", err)
	}

	var solid stl.Solid

	err = stl.CopyAll(bytes.NewReader(data), &solid)
	if err != nil {
		log.Panicf("Oh no! %s", err)
	}

	var normals []float32
	var vertices []float32

	for _, t := range solid.Triangles {

		for _, v := range t.Vertices {
			normals = append(normals, t.Normal[0])
			normals = append(normals, t.Normal[1])
			normals = append(normals, t.Normal[2])

			vertices = append(vertices, v[0])
			vertices = append(vertices, v[1])
			vertices = append(vertices, v[2])
		}
	}

	log.Printf("len of vertices: %d", len(vertices))
	log.Printf("len of normals: %d", len(normals))

	geom.Attributes.Set("position", three.NewBufferAttribute(vertices, 3))
	geom.Attributes.Set("normal", three.NewBufferAttribute(normals, 3))

	mesh := three.NewMesh(geom, mat)
	mesh.Rotation.Set("x", mesh.Rotation.Get("x").Float()+4.5)
	scene.Add(mesh)
	scene.Background = three.NewColor("gray")
	t := js.Global.Get("THREE").Get("OrbitControls").New(camera, renderer.Get("domElement"))
	// start animation
	var animate func()
	animate = func() {
		t.Call("update")
		js.Global.Call("requestAnimationFrame", animate)
		// mesh.Rotation.Set("x", mesh.Rotation.Get("x").Float()+0.01)
		// mesh.Rotation.Set("y", mesh.Rotation.Get("y").Float()+0.01)
		renderer.Render(scene, camera)
	}
	animate()
}
