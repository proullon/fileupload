package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const ImageDirectory = "./test"

func main() {

	// home page for example
	http.HandleFunc("/", home)

	// Upload handler
	http.HandleFunc("/upload", upload)

	// Directory where pictures are stored
	fs := http.FileServer(http.Dir(ImageDirectory))
	http.Handle("/images/", http.StripPrefix("/images/", fs))

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		fmt.Printf("cannot start server: %s\n", err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("upload.gtpl")

	var images []string
	fis, err := ioutil.ReadDir(ImageDirectory)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, fi := range fis {
		images = append(images, fi.Name())
	}

	err = t.Execute(w, images)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 302)
		return
	}
	fmt.Println("method:", r.Method)

	// Get file handler from http request
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Open file on disk
	f, err := os.OpenFile(ImageDirectory+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// We copy the file content from http request to disk
	io.Copy(f, file)

	// Done, now we redirect user to home page
	http.Redirect(w, r, "/", 302)
}
