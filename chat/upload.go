package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

func uploadHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: insecure, use cookie value instead. For now.
	userID := req.FormValue("userid")
	file, header, err := req.FormFile("avatarFile")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// TODO: stream data to tmpfile, then copy if successful instead of loading whole file in memory
	data, err := ioutil.ReadAll(file)

	if err != nil {
		// TODO: is this the correct error? What can go wrong reading the file?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := path.Join("avatars", userID+path.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0777)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Successfully added avatar file %s", filename)

	_, err = io.WriteString(w, "Successful")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
