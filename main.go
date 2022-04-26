package main

import (
	"log"

	"github.com/gopherjs/gopherjs/js"
)

func main() {
	listener := func(e *js.Object) {
		renderButton := js.Global.Get("document").Call("getElementById", "render")
		click := func(event *js.Object) {
			event.Call("stopPropagation")
			eyesSelector := js.Global.Get("document").Call("getElementById", "eyes")
			index := eyesSelector.Get("selectedIndex")
			log.Println(eyesSelector.Get("options").Index(index.Int()).Get("value").String())
		}
		renderButton.Call("addEventListener", "click", click)
	}

	js.Global.Get("document").Call("addEventListener", "DOMContentLoaded", listener)

	eyesSelector := js.Global.Get("document").Call("getElementById", "eyes")
	index := eyesSelector.Get("selectedIndex")
	log.Println(eyesSelector.Get("options").Index(index.Int()).Get("value").String())

	width := js.Global.Get("innerWidth").Float() - 300
	height := js.Global.Get("innerHeight").Float() - 50

	g, err := newGopherThree(width, height)
	if err != nil {
		log.Panicf("err")
	}

	js.Global.Get("document").Call("getElementById", "scene").Call("appendChild", g.Renderer.Get("domElement"))
	t := js.Global.Get("THREE").Get("OrbitControls").New(g.Camera, g.Renderer.Get("domElement"))

	// start animation
	var animate func()
	animate = func() {
		t.Call("update")
		js.Global.Call("requestAnimationFrame", animate)
		g.Render()
	}
	animate()
}
