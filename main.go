package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
)

//go:embed static
var static embed.FS

func main() {
	addr := "localhost:3000"
	fmt.Printf("open http://%s/\n", addr)
	fmt.Println("press ctrl-c to stop")

	root, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}
	http.Handle("/", http.FileServer(http.FS(root)))
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
