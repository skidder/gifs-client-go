package gifs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	// "github.com/davecgh/go-spew/spew"
)

const (
	defaultAPIEndpoint = "https://api.gifs.com"
)

// GIFSClient has functions for adding content to gifs
type GIFSClient interface {
	Import(*ImportRequest) error
	Upload(*UploadRequest, io.Reader) error
}

type gifsClient struct {
	baseURL string
	apiKey  string
}

// ImportRequest represents a request to import externally hosted content
type ImportRequest struct {
	SourceURL   string              `json:"source"`
	Title       string              `json:"title,omitempty"`
	Tags        []string            `json:"tags,omitempty"`
	NSFW        bool                `json:"nsfw"`
	Attribution *AttributionDetails `json:"attribution"`
}

// UploadRequest represents a request to add content that will be uploaded to gifs
type UploadRequest struct {
	Filename    string
	Title       string
	Tags        []string
	NSFW        bool
	Attribution *AttributionDetails
}

// AttributionDetails describes source content attribution info
type AttributionDetails struct {
	Site string `json:"site"`
	User string `json:"user"`
	URL  string `json:"url"`
}

// NewGIFSClient interacts with production gifs API endpoint
func NewGIFSClient(apiKey string) GIFSClient {
	return &gifsClient{
		baseURL: defaultAPIEndpoint,
		apiKey:  apiKey,
	}
}

// NewGIFSClientWithURL interacts with an arbitrary API endpoint
func NewGIFSClientWithURL(baseURL string, apiKey string) GIFSClient {
	return &gifsClient{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (g *gifsClient) Import(request *ImportRequest) error {
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(body)
	var httpRequest *http.Request
	httpRequest, err = http.NewRequest("POST", fmt.Sprintf("%s/media/import", g.baseURL), bodyReader)
	httpRequest.Header.Add("Gifs-API-Key", g.apiKey)
	httpRequest.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	var response *http.Response
	response, err = client.Do(httpRequest)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var responseBody []byte
	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		fmt.Printf("Response body: %s\n", string(responseBody))
		return fmt.Errorf("Status code indicated error: %d", response.StatusCode)
	}
	return nil
}

func (g *gifsClient) Upload(request *UploadRequest, fileReader io.Reader) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if request.Title != "" {
		writer.WriteField("title", request.Title)
	}
	if request.NSFW {
		writer.WriteField("nsfw", "true")
	}
	if len(request.Tags) > 0 {
		tags := ""
		for i, v := range request.Tags {
			if i > 0 {
				tags = v
			} else {
				tags = fmt.Sprintf("%s, %s", tags, v)
			}
		}
		writer.WriteField("tags", tags)
	}
	if request.Attribution != nil {
		if request.Attribution.Site != "" {
			writer.WriteField("attribution_site", request.Attribution.Site)
		}
		if request.Attribution.URL != "" {
			writer.WriteField("attribute_url", request.Attribution.URL)
		}
		if request.Attribution.User != "" {
			writer.WriteField("attribution_user", request.Attribution.User)
		}
	}

	part, err := writer.CreateFormFile("file", request.Filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, fileReader)

	err = writer.Close()
	if err != nil {
		return err
	}

	// fmt.Println("Request Body:\n%s\n", string(body.Bytes()))
	bodyReader := bytes.NewReader(body.Bytes())
	var httpRequest *http.Request
	httpRequest, err = http.NewRequest("POST", fmt.Sprintf("%s/media/upload", g.baseURL), bodyReader)
	httpRequest.Header.Add("Gifs-API-Key", g.apiKey)
	httpRequest.Header.Add("Content-Type", writer.FormDataContentType())

	// spew.Dump(httpRequest)
	client := &http.Client{}
	var response *http.Response
	response, err = client.Do(httpRequest)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var responseBody []byte
	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		fmt.Printf("Response body: %s\n", string(responseBody))
		return fmt.Errorf("Status code indicated error: %d", response.StatusCode)
	}
	return nil
}
