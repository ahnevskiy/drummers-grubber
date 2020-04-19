package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Router struct {
	responsesList []getResponse
	archives      []archiveObject
	LinksList     string
	TimeNow       string
}

type getResponse struct {
	response string
	handler  func(w http.ResponseWriter, r *http.Request)
}

func (r *Router) New(archives *[]archiveObject) {
	r.archives = *archives
	r.responsesList = []getResponse{
		getResponse{"/", r.getIndex},
		getResponse{"/getArchives", r.getArchives},
		getResponse{"/addData", r.addData}}
	r.LinksList = "/getArchives"
	for _, response := range r.responsesList {
		http.HandleFunc(response.response, response.handler)
	}

}

func (r *Router) Start(port int) {
	fmt.Println("Server is listening...")
	http.ListenAndServe(":80", nil)
}

func (r *Router) getArchives(w http.ResponseWriter, req *http.Request) {
	var archivesList []Archive
	for _, archive := range r.archives {
		archivesList = append(archivesList, getArchive(archive.handler))
	}
	fmt.Fprintf(w, convertToJSON(archivesList))
}

func (r *Router) addData(w http.ResponseWriter, req *http.Request) {
	var newData mainData
	var err error
	password := req.FormValue("password")
	if GetMD5Hash(password) != "245831bff5fe7a508af3c7278aad1200" {
		return
	}
	newData.Count, err = strconv.ParseInt(req.FormValue("count"), 10, 64)
	if err != nil {
		panic(err)
	}
	newData.Date = req.FormValue("date")
	targetJSON := req.FormValue("target")
	for _, archive := range r.archives {
		if archive.jsonFile == targetJSON {
			archive.handler.newDataChannel <- newData
			fmt.Fprintf(w, "Data added in:\n")
			fmt.Fprintf(w, "File  : [%s]\n", targetJSON)
			fmt.Fprintf(w, "Date  : [%s]\n", newData.Date)
			fmt.Fprintf(w, "Count : [%d]\n", newData.Count)
		}
	}
}

func getArchive(h ArchiveHandler) Archive {
	h.getDataChannel <- true
	select {
	case msg := <-result:
		return msg
	}
}

func (r *Router) getIndex(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}
	r.TimeNow = fmt.Sprint(time.Now().Format("2006-01-02 15:04:05"))
	tmpl.Execute(w, r)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
