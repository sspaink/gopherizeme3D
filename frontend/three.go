package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/divan/three"
	"github.com/hschendel/stl"
)

type GopherThree struct {
	Scene    *three.Scene
	Light    *three.DirectionalLight
	Camera   three.PerspectiveCamera
	Mesh     *three.Mesh
	Renderer three.WebGLRenderer
}

func newGopherThree(width, height float64) (GopherThree, error) {
	var g GopherThree

	g.Renderer = three.NewWebGLRenderer()
	g.Renderer.SetSize(width, height, true)

	// setup camera and scene
	g.Camera = three.NewPerspectiveCamera(70, width/height, 1, 500)
	g.Camera.Position.Set(0, 0, 100)

	g.Scene = three.NewScene()

	// lights
	g.Light = three.NewDirectionalLight(three.NewColor("white"), 1)
	g.Light.Position.Set(0, 256, 256)
	g.Scene.Add(g.Light)

	// material
	params := three.NewMaterialParameters()
	params.Color = three.NewColor("skyblue")
	mat := three.NewMeshLambertMaterial(params)

	geom := three.NewBufferGeometry()

	resp, err := http.Get("http://localhost:9000/fileserver/GopherPrintable.stl")
	if err != nil {
		return GopherThree{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GopherThree{}, fmt.Errorf("Invalid status code trying to retrieve stl: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GopherThree{}, err
	}

	var solid stl.Solid

	err = stl.CopyAll(bytes.NewReader(data), &solid)
	if err != nil {
		return GopherThree{}, err
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

	g.Mesh = three.NewMesh(geom, mat)
	g.Mesh.Rotation.Set("x", g.Mesh.Rotation.Get("x").Float()+4.5)
	g.Scene.Add(g.Mesh)
	g.Scene.Background = three.NewColor("gray")

	return g, nil
}

func (g *GopherThree) Render() {
	x, y, z := g.Camera.Position.Coords()
	g.Light.Position.Set(x, y, z)
	g.Renderer.Render(g.Scene, g.Camera)
}
