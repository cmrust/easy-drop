// CMR 20171107
// Todo: add command line argument functionality
// Todo: add an alternate download path arg
// Todo: add a file-size upload-limit arg
// Todo: add a port number arg

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Globals
var port = "3767"

// Handlers
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Handle incoming multipart
	mp, header, err := r.FormFile("file")
	defer mp.Close()
	filename := header.Filename
	fmt.Printf("Receiving file: %s...", filename)
	defer fmt.Printf("\n")
	if err != nil {
		// Todo: check that this returns HTTP 400
		fmt.Fprintln(w, err)
		return
	}

	// Setup destination file
	path, _ := os.Getwd()
	path += "/uploads/" + filename
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		// Todo: check that this returns HTTP 500
		fmt.Fprintf(w, "Failed to open "+path+" for writing")
		return
	}

	// Copy contents to disk
	_, err = io.Copy(file, mp)
	if err != nil {
		// Todo: check that this returns HTTP 500
		fmt.Fprintln(w, err)
		return
	}

	// Respond with success
	fmt.Fprintf(w, "File %s uploaded successfully.", filename)
	fmt.Printf(" completed successfully.")
}

// Main
func main() {
	http.HandleFunc("/upload", downloadHandler)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	fmt.Printf("Listening on port: %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
