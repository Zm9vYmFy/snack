package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type fileEntry struct {
	FileID           string
	ProcessingStatus string
}

type detail struct {
	FileID           string
	FileName         string
	MP3Path          string
	OriginalFilePath string
	SeriesTitle      string
}

type detailResponse struct {
	detail *detail
	err    error
}

type segment struct {
	FileSegmentID int64
	FileID        string
	SegmentText   string
	StartTime     int64
	EndTime       int64
}

type segmentResponse struct {
	segments []segment
	err      error
}

type client interface {
	RetrieveAllFiles() ([]fileEntry, error)

	RetrieveDetails(fileID string, response chan<- detailResponse)
	RetrieveSegments(fileID string, response chan<- segmentResponse)
}

func createClient(apiRoot string) client {
	return &defaultClient{apiRoot}
}

type defaultClient struct {
	apiRoot string
}

func (d defaultClient) RetrieveAllFiles() ([]fileEntry, error) {
	r, err := http.Get(fmt.Sprintf("%sall", d.apiRoot))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var files []fileEntry
	if err := json.NewDecoder(r.Body).Decode(&files); err != nil {
		return nil, err
	}

	return files, nil
}

func (d defaultClient) RetrieveDetails(fileID string, response chan<- detailResponse) {
	r, err := http.Get(fmt.Sprintf("%sdetails/%s", d.apiRoot, fileID))
	if err != nil {
		response <- detailResponse{nil, err}
	}
	defer r.Body.Close()

	var deets detail
	if err := json.NewDecoder(r.Body).Decode(&deets); err != nil {
		response <- detailResponse{nil, err}
	}

	response <- detailResponse{&deets, nil}
}

func (d defaultClient) RetrieveSegments(fileID string, response chan<- segmentResponse) {
	r, err := http.Get(fmt.Sprintf("%ssegments/%s", d.apiRoot, fileID))
	if err != nil {
		response <- segmentResponse{nil, err}
	}
	defer r.Body.Close()

	var segments []segment
	if err := json.NewDecoder(r.Body).Decode(&segments); err != nil {
		response <- segmentResponse{nil, err}
	}

	response <- segmentResponse{segments, nil}
}
