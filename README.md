# Pkger

[`github.com/markbates/pkger`](https://godoc.org/github.com/markbates/pkger) is a tool for embedding static files into Go binaries. It will, hopefully, be a replacement for [`github.com/gobuffalo/packr/v2`](https://godoc.org/github.com/gobuffalo/packr/v2).

### Requirements

* Go 1.13+
* Go Modules

## How it Works (Module Aware Pathing)

Pkger is powered by the dark magic of Go Modules, so they're like, totally required.

With Go Modules pkger can resolve packages with accuracy. No more guessing and trying to
figure out build paths, GOPATHS, etc... for this tired old lad.

With the module's path correctly resolved, it can serve as the "root" directory for that
module, and all files in that module's directory are available.

Paths:
* Paths should use UNIX style paths:
  `/cmd/pkger/main.go`
* If unspecified the path's package is assumed to be the current module.
* Packages can specified in at the beginning of a path with a `:` seperator.
github.com/markbates/pkger:/cmd/pkger/main.go
* There are no relative paths. All paths are absolute to the modules root.

```
"github.com/gobuffalo/buffalo:/go.mod" => $GOPATH/pkg/mod/github.com/gobuffalo/buffalo@v0.14.7/go.mod
```

## CLI

### Installation

```bash
$ go get github.com/markbates/pkger/cmd/pkger
$ pkger -h
```
## Usage

### In Code

The first step in using Packr is to create a new box. A box represents a folder on disk. Once you have a box you can get `string` or `[]byte` representations of the file.

```go
// set up a new box by giving it a (relative) path to a folder on disk:
box := packr.NewBox("./templates")

// Get the string representation of a file, or an error if it doesn't exist:
html, err := box.FindString("index.html")

// Get the []byte representation of a file, or an error if it doesn't exist:
html, err := box.FindBytes("index.html")
```

### What is a Box?

A box represents a folder, and any sub-folders, on disk that you want to have access to in your binary. When compiling a binary using the `packr` CLI the contents of the folder will be converted into Go files that can be compiled inside of a "standard" go binary. Inside of the compiled binary the files will be read from memory. When working locally the files will be read directly off of disk. This is a seamless switch that doesn't require any special attention on your part.

#### Example

Assume the follow directory structure:

```
├── main.go
└── templates
    ├── admin
    │   └── index.html
    └── index.html
```

The following program will read the `./templates/admin/index.html` file and print it out.

```go
package main

import (
  "fmt"

  "github.com/gobuffalo/packr"
)

func main() {
  box := packr.NewBox("./templates")

  s, err := box.FindString("admin/index.html")
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(s)
}
```

### Development Made Easy

In order to get static files into a Go binary, those files must first be converted to Go code. To do that, Packr, ships with a few tools to help build binaries. See below.

During development, however, it is painful to have to keep running a tool to compile those files.

Packr uses the following resolution rules when looking for a file:

1. Look for the file in-memory (inside a Go binary)
1. Look for the file on disk (during development)

Because Packr knows how to fall through to the file system, developers don't need to worry about constantly compiling their static files into a binary. They can work unimpeded.

Packr takes file resolution a step further. When declaring a new box you use a relative path, `./templates`. When Packr receives this call it calculates out the absolute path to that directory. By doing this it means you can be guaranteed that Packr can find your files correctly, even if you're not running in the directory that the box was created in. This helps with the problem of testing, where Go changes the `pwd` for each package, making relative paths difficult to work with. This is not a problem when using Packr.

---

### Usage

```bash
$ pkger
```

The result will be a `pkged.go` file in the **root** of the module with the embedded information and the package name of the module.

```go
// ./pkged.go
package <.>

// Pkger stuff here
```

The `-o` flag can be used specify the directory of the `pkged.go` file.

```bash
$ pkger -o cmd/reader
```

The result will be a `pkged.go` file in the **cmd/reader** folder with the embedded information and the package name of that folder.

```go
// cmd/reader/pkged.go
package <reader>

// Pkger stuff here
```

### Including Files at Package Time

There may be reasons where you don't reference a particular file, or folder, that you want embedded in your application, such as a build artifact.

To do this you may use either the [`github.com/markbates/pkger#Include`](https://godoc.org/github.com/markbates/pkger#Include) function to set a no-op parser directive to include the specified path.

Alternatively, you may use the `-include` flag with the `pkger` and `pkger list` commands.

```bash
$ pkger list -include /actions -include github.com/gobuffalo/buffalo:/app.go

app
 > app:/actions
 > app:/actions/actions.go
 > app:/assets
 > app:/assets/css
 > app:/assets/css/_buffalo.scss
 > app:/assets/css/application.scss
 > app:/assets/images
 > app:/assets/images/favicon.ico
 > app:/assets/images/logo.svg
 > app:/assets/js
 > app:/assets/js/application.js
 > app:/go.mod
 > app:/locales/all.en-us.yaml
 > app:/public
 > app:/public/assets
 > app:/public/assets/.keep
 > app:/public/assets/app.css
 > app:/public/images
 > app:/public/images/img1.png
 > app:/public/index.html
 > app:/public/robots.txt
 > app:/templates
 > app:/templates/_flash.plush.html
 > app:/templates/application.plush.html
 > app:/templates/index.plush.html
 > app:/web
 > app:/web/web.go
 > github.com/gobuffalo/buffalo:/app.go
 > github.com/gobuffalo/buffalo:/logo.svg
```

## Reference Application

The reference application for the `README` examples, as well as all testing, can be found at [https://github.com/markbates/pkger/tree/master/pkging/pkgtest/testdata/ref](https://github.com/markbates/pkger/tree/master/pkging/pkgtest/testdata/ref).

```
├── actions
│   └── actions.go
├── assets
│   ├── css
│   │   ├── _buffalo.scss
│   │   └── application.scss
│   ├── images
│   │   ├── favicon.ico
│   │   └── logo.svg
│   └── js
│       └── application.js
├── go.mod
├── go.sum
├── locales
│   └── all.en-us.yaml
├── main.go
├── mod
│   └── mod.go
├── models
│   └── models.go
├── public
│   ├── assets
│   │   └── app.css
│   ├── images
│   │   └── img1.png
│   ├── index.html
│   └── robots.txt
├── templates
│   ├── _flash.plush.html
│   ├── application.plush.html
│   └── index.plush.html
└── web
    └── web.go

13 directories, 20 files
```


## API Usage

Pkger's API is modeled on that of the [`os`](https://godoc.org/os) package in Go's standard library. This makes Pkger usage familiar to Go developers.

The two most important interfaces are [`github.com/markbates/pkger/pkging#Pkger`](https://godoc.org/github.com/markbates/pkger/pkging#Pkger) and [`github.com/markbates/pkger/pkging#File`](https://godoc.org/github.com/markbates/pkger/pkging#File).

```go
type Pkger interface {
  Parse(p string) (Path, error)
  Current() (here.Info, error)
  Info(p string) (here.Info, error)
  Create(name string) (File, error)
  MkdirAll(p string, perm os.FileMode) error
  Open(name string) (File, error)
  Stat(name string) (os.FileInfo, error)
  Walk(p string, wf filepath.WalkFunc) error
  Remove(name string) error
  RemoveAll(path string) error
}

type File interface {
  Close() error
  Info() here.Info
  Name() string
  Open(name string) (http.File, error)
  Path() Path
  Read(p []byte) (int, error)
  Readdir(count int) ([]os.FileInfo, error)
  Seek(offset int64, whence int) (int64, error)
  Stat() (os.FileInfo, error)
  Write(b []byte) (int, error)
}
```

These two interfaces, along with the [`os#FileInfo`](https://godoc.org/os#FileInfo), provide the bulk of the API surface area.

### Open

```go
func run() error {
	f, err := pkger.Open("/public/index.html")
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	fmt.Println("Name: ", info.Name())
	fmt.Println("Size: ", info.Size())
	fmt.Println("Mode: ", info.Mode())
	fmt.Println("ModTime: ", info.ModTime())

	if _, err := io.Copy(os.Stdout, f); err != nil {
		return err
	}
	return nil
}
```

### Stat

```go
func run() error {
	info, err := pkger.Stat("/public/index.html")
	if err != nil {
		return err
	}

	fmt.Println("Name: ", info.Name())
	fmt.Println("Size: ", info.Size())
	fmt.Println("Mode: ", info.Mode())
	fmt.Println("ModTime: ", info.ModTime())

	return nil
}
```

### Walk

```go
func run() error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug)
	defer w.Flush()

	return pkger.Walk("/public", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fmt.Fprintf(w,
			"%s \t %d \t %s \t %s \t\n",
			info.Name(),
			info.Size(),
			info.Mode(),
			info.ModTime().Format(time.RFC3339),
		)

		return nil
	})

}
```

## Understanding the Parser

The [`github.com/markbates/pkger/parser#Parser`](https://godoc.org/github.com/markbates/pkger/parser#Parser) works by statically analyzing the source code of your module using the [`go/parser`](https://godoc.org/go/parser) to find a selection of declarations.

The following declarations in your source code will tell the parser to embed files or folders.

* `pkger.Dir("<path>")` - Embeds all files under the specified path.
* `pkger.Open("<path>")` - Embeds the file, or folder, of the specified path.
* `pkger.Stat("<path>")` - Embeds the file, or folder, of the specified path.
* `pkger.Walk("<path>", filepath.WalkFunc)` - Embeds all files under the specified path.
* `pkger.Include("<path>")` - `Include` is a no-op that directs the pkger tool to include the desired file or folder.

### CLI Usage

To see what declarations the parser has found, you can use the `pkger parse` command to get a `JSON` list of the declarations.

```bash
$ pkger parse

{
 ".": [
  {
   "file": {
    "Abs": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref/foo/bar/baz",
    "Path": {
     "Pkg": "app",
     "Name": "/foo/bar/baz"
    },
    "Here": {
     "Dir": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref",
     "ImportPath": "app",
     "Module": {
      "Path": "app",
      "Main": true,
      "Dir": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref",
      "GoMod": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref/go.mod",
      "GoVersion": "1.13"
     },
     "Name": "main"
    }
   },
   "pos": {
    "Filename": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref/main.go",
    "Offset": 629,
    "Line": 47,
    "Column": 27
   },
   "type": "pkger.MkdirAll",
   "value": "/foo/bar/baz"
  },
  {
   "file": {
    "Abs": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref/foo/bar/baz/biz.txt",
    "Path": {
     "Pkg": "app",
     "Name": "/foo/bar/baz/biz.txt"
    },
    "Here": {
     "Dir": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref",
     "ImportPath": "app",
     "Module": {
      "Path": "app",
      "Main": true,
      "Dir": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref",
      "GoMod": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref/go.mod",
      "GoVersion": "1.13"
     },
     "Name": "main"
    }
   },
   "pos": {
    "Filename": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref/main.go",
    "Offset": 706,
    "Line": 51,
    "Column": 25
   },
   "type": "pkger.Create",
   "value": "/foo/bar/baz/biz.txt"
  },
  ...
 ]
}
```

For a module aware list use the `pkger list` command.

```bash
$ pkger list

app
 > app:/assets
 > app:/assets/css
 > app:/assets/css/_buffalo.scss
 > app:/assets/css/application.scss
 > app:/assets/images
 > app:/assets/images/favicon.ico
 > app:/assets/images/logo.svg
 > app:/assets/js
 > app:/assets/js/application.js
 > app:/go.mod
 > app:/locales/all.en-us.yaml
 > app:/public
 > app:/public/assets
 > app:/public/assets/.keep
 > app:/public/assets/app.css
 > app:/public/images
 > app:/public/images/img1.png
 > app:/public/index.html
 > app:/public/robots.txt
 > app:/templates
 > app:/templates/_flash.plush.html
 > app:/templates/application.plush.html
 > app:/templates/index.plush.html
 > app:/web
 > app:/web/web.go
 > github.com/gobuffalo/buffalo:/logo.svg
```

The `-json` flag can be used to get a more detailed list in JSON.

```bash
$ pkger list -json

{
 "ImportPath": "app",
 "Files": [
  {
   "Abs": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref/assets",
   "Path": {
    "Pkg": "app",
    "Name": "/assets"
   },
   "Here": {
    "Dir": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref/assets",
    "ImportPath": "",
    "Module": {
     "Path": "app",
     "Main": true,
     "Dir": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref",
     "GoMod": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref/go.mod",
     "GoVersion": "1.13"
    },
    "Name": "assets"
   }
  },
  {
   "Abs": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref/assets/css",
   "Path": {
    "Pkg": "app",
    "Name": "/assets/css"
   },
   "Here": {
    "Dir": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref/assets",
    "ImportPath": "",
    "Module": {
     "Path": "app",
     "Main": true,
     "Dir": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref",
     "GoMod": "$GOPATH/src/github.com/markbates/pkger/pkging/pkgtest/testdata/ref/go.mod",
     "GoVersion": "1.13"
    },
    "Name": "assets"
   }
  },
  ...
}
```

