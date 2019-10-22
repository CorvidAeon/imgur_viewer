package main

import (
  "fmt"
  "bytes"
  "mime/multipart"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36"

//AlbumResponseWrapper cause why not make life more difficult
type AlbumResponseWrapper struct {
	Data Album `json:"data"`
	Success bool `json:"success"`
	Status int `json:"status"`
}

//Album is an imgur API model
type Album struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Datetime int `json:"datetime"`
	Cover string `json:"cover"`
	CoverWidth int `json:"cover_width"`
	CoverHeight int `json:"cover_height"`
	AccountURL string `json:"account_url"`
	AccountID int `json:"account_id"`
	Privacy string `json:"privacy"`
	Layout string `json:"layout"`
	Views int `json:"views"`
	Link string `json:"link"`
	Favorite bool `json:"favorite"`
	Nsfw bool `json:"nsfw"`
	Section string `json:"section"`
	Order int `json:"order"`
	DeleteHash *string `json:"deletehash"`
	ImagesCount int `json:"images_count"`
	Images []Image `json:"images"`
	InGallery bool `json:"in_gallery"`
}

//Image is an imgur API model
type Image struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Datetime int `json:"datetime"`
	MimeType string `json:"type"`
	Animated bool `json:"animated"`
	Width int `json:"width"`
	Height int `json:"height"`
	Size int `json:"size"`
	Views int `json:"views"`
	Bandwidth int `json:"bandwidth"`
	DeleteHash *string `json:"deletehash"`
	Name string `json:"name"`
	Section string `json:"section"`
	Link string `json:"link"`
	Gifv string `json:"gifv"`
	Mp4 string `json:"mp4"`
	Looping bool `json:"looping"`
	Favorite bool `json:"favorite"`
	Nsfw bool `json:"nsfw"`
	Vote string `json:"vote"`
	InGallery bool `json:"in_gallery"`
}
func handler(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "Testing: %s", r.URL.Path[1:])
	if r.URL.Path[1:] == "" {
		fmt.Fprintf(w, "Enter an imgur hash to load album.")
		return
	}
	images := imgurAPIAlbumHandler(r.URL.Path[1:])

	var body string
	//images
	body = `<!doctype html>
	<html>
	
	<head>
	<meta name="referrer" content="no-referrer" />
	</head>
	<body>`
	for _, image :=range images {
		body = body + fmt.Sprintf("<img src=\"%s\"></img>", image)
	}
	body = body + `
	</body>

	</html>`
	w.Header().Set("Accept", "application/json")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, body)
}

func main() {
	http.HandleFunc("/", handler)

	http.ListenAndServe(":8080", nil)

}

func imgurAPIAlbumHandler(hash string) (images []string) {
	url := "https://api.imgur.com/3/album/" + hash
	method := "GET"
  
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	err := writer.Close()
	if err != nil {  fmt.Println(err)}
  
  
	client := &http.Client {
	  CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	  },
	}
	req, err := http.NewRequest(method, url, payload)
  
	if err != nil {
	  fmt.Println(err)
	}
	req.Header.Set("User-Agent", userAgent)
  
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
  
	var wrapper AlbumResponseWrapper
	var album Album
	json.Unmarshal(body, &wrapper)
	album = wrapper.Data
	for _, image := range album.Images {
		images = append(images, image.Link)
	}
	return 
}
