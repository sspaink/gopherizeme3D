package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	. "github.com/dave/jennifer/jen"
	"github.com/gin-gonic/gin"
)

// FOR LOCAL TESTING
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func GenerateProgram() error {
	err := os.Mkdir("generated", 0777)
	if err != nil {
		return err
	}

	f := NewFile("main")
	f.Func().Id("main").Params().Block(
		List(Id("body"), Id("_")).Op(":=").Qual("github.com/deadsy/sdfx/sdf", "Sphere3D").Call(Lit(45)),
		Qual("github.com/deadsy/sdfx/render", "ToSTL").Call(Id("body"), Lit("gopher.stl"), Qual("github.com/deadsy/sdfx/render", "NewMarchingCubesOctree").Call(Lit(300))),
	)

	err = os.WriteFile("generated/main.go", []byte(fmt.Sprintf("%#v", f)), 0777)
	if err != nil {
		return err
	}
	err = os.Chdir("generated")
	if err != nil {
		return err
	}
	err = exec.Command("go", "mod", "init", "generated").Run()
	if err != nil {
		return err
	}
	err = exec.Command("go", "mod", "tidy").Run()
	if err != nil {
		return err
	}
	out, _ := exec.Command("ls").Output()
	fmt.Println(string(out))
	cmd := exec.Command("go", "run", "main.go")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.StaticFile("fileserver/js/three.js", "./js/three.js")
	r.StaticFile("fileserver/js/OrbitControls.js", "./js/OrbitControls.js")
	r.StaticFile("fileserver/GopherPrintable.stl", "./GopherPrintable.stl")

	r.POST("stl", func(c *gin.Context) {
		err := GenerateProgram()
		if err != nil {
			fmt.Print(err)
			return
		}
	})

	err := r.Run(":9000")
	if err != nil {
		log.Panic(err)
	}
}
