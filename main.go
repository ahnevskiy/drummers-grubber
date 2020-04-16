package main

import (
	"flag"
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
func argsParse() {
	// Parse arguments
	flag.IntVar(&serverPort, "p", 8000, "Number of port")
	flag.Parse()
	// Print arguments in terminal
	fmt.Printf("Number of port: [%d]\n", serverPort)
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
	argsParse()
	archives := []archiveObject{
		archiveObject{parseResumeHH, "hh_resume.json", ArchiveHandler{}},
		archiveObject{parseVacantHH, "hh_vacant.json", ArchiveHandler{}},
	}
	serveArchives(&archives)
	var r router
	r.New(&archives)
	r.Start(serverPort)
}
