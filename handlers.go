package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ArchiveHandler is a struct to serve an archive
type ArchiveHandler struct {
	parseMethod         func([]byte) int64
	resultsFileName     string
	newValueChannel     chan int64
	getStatisticChannel chan bool
	getDataChannel      chan bool
	sourceURL           string
}

// New creates a new handler object
func (h *ArchiveHandler) New(parseMethod func([]byte) int64, archiveFileName string) {
	h.parseMethod = parseMethod
	h.resultsFileName = archiveFileName
	h.newValueChannel = make(chan int64)
	h.getStatisticChannel = make(chan bool)
	h.getDataChannel = make(chan bool)
	h.sourceURL = readArchive(h.resultsFileName).SourceURL
}

// GetHTML returns a slice of bytes with html context of url
func GetHTML(url string) []byte {
	req, err := http.Get(url)
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}
	defer req.Body.Close()

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}
	return b
}

func grabContent(url string, parser func([]byte) int64) int64 {
	html := GetHTML(url)
	return parser(html)
}

//Serve is a infinite loop for goroutine, where we're collect data every day and form ansver for GET response
func (h *ArchiveHandler) Serve() {
	go h.grubLoop()
	for {
		select {
		case newValue := <-h.newValueChannel:
			archive := readArchive(h.resultsFileName)
			// currentDate := fmt.Sprint(time.Now().Format("2006-01-02"))
			// lastData := archive.getLastData()
			// if currentDate != lastData.Date {
			archive.addData(newValue)
			archive.saveData(h.resultsFileName)
			// }
		case response := <-h.getStatisticChannel:
			if response {
				archive := readArchive(h.resultsFileName)
				lastData := archive.getLastData()
				s := ""
				s += fmt.Sprintf("%s\n", archive.Description)
				s += fmt.Sprintf("\nStart date     : [%s]", archive.Data[0].Date)
				s += fmt.Sprintf("\nLast date      : [%s]\n", lastData.Date)
				s += fmt.Sprintf("\nLast count     : [%d]", lastData.Count)
				s += fmt.Sprintf("\nAbsolute delta : [%d]", archive.AbsoluteDelta)
				s += fmt.Sprintf("\nLast day delta : [%d]\n\n\n", lastData.Delta)
				result <- s
			}
		case response := <-h.getDataChannel:
			if response {
				archive := readArchive(h.resultsFileName)
				result <- archive.convertToJSON()
			}
		}
	}

}

func (h *ArchiveHandler) grubLoop() {
	for {
		currentCount := grabContent(h.sourceURL, h.parseMethod)
		h.newValueChannel <- currentCount
		time.Sleep(time.Minute)
		// lastDate := fmt.Sprint(time.Now().Format("2006-01-02"))
		// currentDate := fmt.Sprint(time.Now().Format("2006-01-02"))
		// for lastDate == currentDate {
		// 	time.Sleep(time.Minute)
		// 	currentDate = fmt.Sprint(time.Now().Format("2006-01-02"))
		// }
	}
}
