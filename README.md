# gifs-client-go

[![GoDoc](https://godoc.org/github.com/skidder/gifs-client-go/gifs?status.svg)](https://godoc.org/github.com/skidder/gifs-client-go/gifs) [![Circle CI](https://circleci.com/gh/skidder/gifs-client-go.svg?style=svg)](https://circleci.com/gh/skidder/gifs-client-go)

Go Client for the [GIFS API](http://docs.gifs.com/v1.0/docs).

# Installation

```shell
go get -u github.com/skidder/gifs-client-go/gifs
```

# Operations

## Import from URL
Invokes the [Import](http://docs.gifs.com/docs/mediaimport) endpoint in the GIFS API.  Sample client code follows:

```go
package main

import (
	"fmt"
	"os"

	"github.com/skidder/gifs-client-go/gifs"
)

func main() {
	client := gifs.NewGIFSClient(os.Getenv("GIFS_API_KEY"))

	request := &gifs.ImportRequest{
		Title:     os.Args[1],
		SourceURL: os.Args[2],
	}
	err := client.Import(request)
	if err != nil {
		fmt.Printf("Error importing file: %s\n", err.Error())
		return
	}
}
```


## Upload File
Invokes the [Upload](http://docs.gifs.com/docs/mediaupload) endpoint on the GIFS API.  Sample code follows:

```go
package main

import (
	"fmt"
	"os"

	"github.com/skidder/gifs-client-go/gifs"
)

func main() {
	client := gifs.NewGIFSClient(os.Getenv("GIFS_API_KEY"))

	title := os.Args[1]
	filename := os.Args[2]
	inputFile, err := os.OpenFile(filename, os.O_RDONLY, os.FileMode(666))
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		return
	}
	defer inputFile.Close()

	request := &gifs.UploadRequest{
		Title:    title,
		Filename: filename,
	}
	err = client.Upload(request, inputFile)
	if err != nil {
		fmt.Printf("Error uploading file: %s\n", err.Error())
		return
	}
}
```
