# gifs-client-go
Go Client for the gifs API

# Operations

## Import from URL

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

```go
package main

import (
	"fmt"
	"os"

	"github.com/skidder/gifs-client-go/gifs"
)

func main() {
	client := gifs.NewGIFSClient(os.Getenv("GIFS_API_KEY"))

	inputFile, err := os.OpenFile(os.Args[2], os.O_RDONLY, os.FileMode(666))
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		return
	}
	defer inputFile.Close()

	request := &gifs.UploadRequest{
		Title:    os.Args[1],
		Filename: os.Args[2],
	}
	err = client.Upload(request, inputFile)
	if err != nil {
		fmt.Printf("Error uploading file: %s\n", err.Error())
		return
	}
}
```
