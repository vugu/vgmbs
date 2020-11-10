# Vugu Material/Bootstrap

A Bootstrap build with Material-ish extensions, made for integration into Vugu applications. 

## Setup

You'll need a SASS compiler to build the SASS files from this project.  We recommend `vgsassc`, a wrapper
around [libsass](https://sass-lang.com/libsass) designed for this purpose.

### Windows

To build vgsassc on Windows you'll need a GCC-compatible C compiler.  We recommend: https://jmeubank.github.io/tdm-gcc/
Once installed you should be able to type `gcc --help` at the Command Prompt and get help info instead of an error, at
which point you're good to move to the next step.

### Installing vgsassc 

Install vgsassc with:
```bash
go get -u github.com/vugu/vgsassc
```

## Your SCSS File(s)

You can include the main bootstrap SCSS into your own file with an `@import`.  For example, a main.scss that contains
Bootstrap plus your own CSS could look like:

### main.scss
```scss
@import "vgmbs/bootstrap";

#something {
  // etc
}
```

## Create build.go

The recommended approach to compiling your SASS files into CSS is to make a standalone Go program in a file called
build.go and then running this via `go run` whenever needed.

Example build.go:
```go
// +build ignore

package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"

    "github.com/vugu/vgmbs"
    "github.com/vugu/vugu/distutil"
)

func main() {
    
    // write out vgmbs files
    os.MkdirAll("scss/vgmbs", 0755)
    vgmbs.NewFileWriter("scss/vgmbs").MustWrite()

    // then build everything
    os.MkdirAll("dist/css", 0755)
    cmd = exec.Command("vgsassc",
        "-o", "dist/css/main.css",
        "-I", "web/scss",
        "web/scss/main.scss",
    )
    distutil.Must(err)
    out, err = cmd.CombinedOutput()
    fmt.Printf("%s", out)
    distutil.Must(err)

}
``` 

## Call build.go

For many applications, you'll want to call build.go from two places: wherever your overall project build/distribution
is done, and in an appropriate HTTP handler in order to automatically reload pages.

Calling build.go from the command line or for example a shell script is trivial:

```
go run build.go
```

You can also set up an HTTP handler to run build.go when pages are loaded during development, so you get a fast
and easy to use workflow where only a browser refresh is required.  Example:
 
```go
type MyHandler struct {
    BaseDir string // set this to the project directory
    DevMode bool   // set to true in development mode
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    distDir := http.Dir(filepath.Join(h.BaseDir, "dist")) // 

    // when a page with no file extension is requested, run build.go first
    ext := path.Ext(r.URL.Path)
    if h.DevMode && ext == "" {

        cmd := exec.Command("go", "run", "build.go")
        cmd.Dir = h.BaseDir
        b, err := cmd.CombinedOutput()
        log.Printf("go run build.go - err: %v\n%s", err, b)
        // send the error back in the page request if build.go fails
        if err != nil { 
            w.Header().Set("Content-Type", "text/plain")
            w.WriteHeader(500)
            fmt.Fprintf(w, "build-frontend.go error: %v\n%s", err, b)
            return
        }

    }

    // more serving logic here...
}
```
