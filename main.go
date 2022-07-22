package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
)

var (
	storageClient *storage.Client
	project       string
)

func main() {
	flag.StringVar(&project, "project", "example", "Google Cloud project (default: example)")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/download/{bucket}/{filename:.*}", Download).Methods("GET")
	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		log.Fatal(err)
	}
}

func Download(w http.ResponseWriter, r *http.Request) {
	bucket := mux.Vars(r)["bucket"]
	filename := mux.Vars(r)["filename"]
	log.Printf("Request to download file %s from bucket %s (project: %s)", filename, bucket, project)

	clientCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reader, err := storageClient.Bucket(bucket).UserProject(project).Object(filename).NewReader(clientCtx)
	if err != nil {
		log.Printf("Error creating file reader: %v", err)
		status := http.StatusInternalServerError
		if errors.Is(err, storage.ErrObjectNotExist) {
			status = http.StatusNotFound
		}
		w.WriteHeader(status)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", reader.Attrs.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", strconv.FormatInt(reader.Attrs.Size, 10))
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, reader); err != nil {
		log.Printf("Error downloading file: %v", err)
	}
}
