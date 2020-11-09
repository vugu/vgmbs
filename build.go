// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {

	err := vfsgen.Generate(http.Dir("scss"), vfsgen.Options{
		PackageName: "vgmbs",
		//BuildTags:    "!dev",
		VariableName: "assets",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Disable dist folder stuff for now, no current need for it.

	//os.MkdirAll("dist", 0755)
	//
	//scssFileList := []string{
	//	"scss/reset.scss",
	//	"scss/globals.scss",
	//	"scss/elements.scss",
	//}
	//
	//var cmd *exec.Cmd
	//var b []byte
	//
	//// compile sass
	//cmd = exec.Command("vgsassc", append([]string{"-o", "dist/material.css"}, scssFileList...)...)
	//b, err = cmd.CombinedOutput()
	//log.Printf("vgsassc output:\n%s", b)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// compile sass again minified
	//cmd = exec.Command("vgsassc", append([]string{"-o", "dist/material.min.css"}, scssFileList...)...)
	//b, err = cmd.CombinedOutput()
	//log.Printf("vgsassc output:\n%s", b)
	//if err != nil {
	//	panic(err)
	//}

}
