package main

import (
	render "github.com/deadsy/sdfx/render"
	sdf "github.com/deadsy/sdfx/sdf"
)

func main() {
	body, _ := sdf.Sphere3D(45)
	render.RenderSTL(body, 300, "gopher.stl")
}
