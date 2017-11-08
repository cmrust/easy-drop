// CMR 20171107
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Globals
var port = "8114"

// Handlers
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Handle incoming multipart
	mp, header, err := r.FormFile("file")
	defer mp.Close()
	filename := header.Filename
	fmt.Printf("Receiving file: %s...", filename)
	defer fmt.Printf("\n")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	// Setup destination file
	path, _ := os.Getwd()
	path += "/uploads/" + filename
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		fmt.Fprintf(w, "Failed to open "+path+" for writing")
		return
	}

	// Copy contents to disk
	_, err = io.Copy(file, mp)
	if err != nil {
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
