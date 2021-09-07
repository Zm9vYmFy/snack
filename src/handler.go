package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type response struct {
	Details  *detail
	Segments []segment
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}

type endpointHandler struct {
	client client
}

func getFile(files []fileEntry, fileID string) *fileEntry {
	for _, file := range files {
		if file.FileID == fileID && file.ProcessingStatus == "FINISHED" {
			return &file
		}
	}
	return nil
}

func (e endpointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 { // starts with a / as well in the path
		writeError(w, 400, "Invalid path for endpoint. Include (only) fileID")
		return
	}

	var fileID = pathParts[2]
	if files, err := e.client.RetrieveAllFiles(); err != nil {
		log.Printf("Failed to retrieve list of files: %v", err)
		writeError(w, 500, "Failed to retrieve files")
		return
	} else if file := getFile(files, fileID); file != nil && file.ProcessingStatus != "FINISHED" {
		writeError(w, 418, "File pending")
		return
	} else if file == nil {
		writeError(w, 404, "File not found")
		return
	}

	var (
		detailResponses  = make(chan detailResponse, 1)
		segmentResponses = make(chan segmentResponse, 1)
		deets            detailResponse
		segmentR         segmentResponse
		endpointResponse response
	)

	// Parallelize the API calls, expect responses to come back over a channel for each
	go e.client.RetrieveDetails(fileID, detailResponses)
	go e.client.RetrieveSegments(fileID, segmentResponses)

	deets = <-detailResponses
	if deets.err != nil {
		log.Printf("Failed to query details: %v\n", deets.err)
	} else {
		endpointResponse.Details = deets.detail
	}

	segmentR = <-segmentResponses
	if segmentR.err != nil {
		log.Printf("Failed to query segments: %v\n", segmentR.err)
	} else {
		endpointResponse.Segments = segmentR.segments
	}

	if err := json.NewEncoder(w).Encode(endpointResponse); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
