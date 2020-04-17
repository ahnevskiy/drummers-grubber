package main

import (
	"fmt"
	"path/filepath"
)

var (
	result chan string
	// Variables for input arguments
	serverPort int
)

type archiveObject struct {
	parseMethod func([]byte) int64
	jsonFile    string
	handler     ArchiveHandler
}

func init() {
	result = make(chan string)
}

func serveArchives(archives *[]archiveObject) {
	archs := *archives
	for i := range archs {
		archs[i].handler.New(archs[i].parseMethod, filepath.Join("data", archs[i].jsonFile))
		go archs[i].handler.Serve()
		fmt.Printf("Start serve [%s]\n", archs[i].jsonFile)
	}
}

func main() {
	archives := []archiveObject{
		archiveObject{parseResumeHH, "hh_resume.json", ArchiveHandler{}},
		archiveObject{parseVacantHH, "hh_vacant.json", ArchiveHandler{}},
	}
	serveArchives(&archives)
	var r router
	r.New(&archives)
	r.Start(serverPort)
}
