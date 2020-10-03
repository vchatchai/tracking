// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {

	var fs http.FileSystem = http.Dir("../../webapp/build")

	err := vfsgen.Generate(fs, vfsgen.Options{
		PackageName:  "generate",
		BuildTags:    "!dev",
		VariableName: "Assets",
	})

	if err != nil {
		log.Fatalln(err)
	}
}
