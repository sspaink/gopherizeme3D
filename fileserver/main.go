// package main

// import (
// 	"log"
// 	"net/http"
// )

// func main() {

// 	// create file server handler
// 	fs := http.FileServer(http.Dir("/home/sspaink/Sandbox/gopherizeme3D/fileserver"))

// 	// start HTTP server with `fs` as the default handler
// 	log.Fatal(http.ListenAndServe(":9000", fs))

// }

package main

import (
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("/home/sspaink/Sandbox/gopherizeme3D/fileserver"))

	var wrapped = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		fs.ServeHTTP(w, r)
	})
	http.ListenAndServe(":9000", wrapped)
}
