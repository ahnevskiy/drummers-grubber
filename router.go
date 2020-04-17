package main

import (
	"fmt"

	"github.com/go-martini/martini"
)

var responsesList string

type router struct {
	server        *martini.ClassicMartini
	responsesList []getResponse
	archives      []archiveObject
}

type getResponse struct {
	response string
	handler  func() string
}

func (r *router) New(archives *[]archiveObject) {
	r.archives = *archives
	r.server = martini.Classic()
	r.responsesList = []getResponse{
		getResponse{"/", r.getIndex},
		getResponse{"/getResponses", r.getResponses},
		getResponse{"/getSummaryStat", r.getSummaryStat},
		getResponse{"/getResumeStat", r.getResumeStat},
		getResponse{"/getVacantStat", r.getVacantStat},
		getResponse{"/getResumeData", r.getResumeData},
		getResponse{"/getVacantData", r.getVacantData}}
	for _, response := range r.responsesList {
		r.server.Get(response.response, response.handler)
		responsesList += fmt.Sprintf("%s\n", response.response)
	}

}

func (r *router) Start(port int) {
	r.server.Run()

}

func getStatistic(h ArchiveHandler) string {
	h.getStatisticChannel <- true
	select {
	case msg := <-result:
		return msg
	}
}

func getData(h ArchiveHandler) string {
	h.getDataChannel <- true
	select {
	case msg := <-result:
		return msg
	}
}

func (r *router) getIndex() string {
	s := "Grabber for HH\n"
	s += "--------------------------------\n"
	s += "Available GET responses:\n"
	s += responsesList
	s += "--------------------------------\n"
	s += "Based on GO\n"
	s += "Author @ahnevskiy"
	return s
}
func (r *router) getResponses() string {
	return responsesList
}

func (r *router) getSummaryStat() string {
	s := ""
	for i := range r.archives {
		s += getStatistic(r.archives[i].handler)
	}
	return s
}

func (r *router) getResumeStat() string {
	return getStatistic(r.archives[0].handler)
}

func (r *router) getVacantStat() string {
	return getStatistic(r.archives[1].handler)
}

func (r *router) getResumeData() string {
	return getData(r.archives[0].handler)
}

func (r *router) getVacantData() string {
	return getData(r.archives[1].handler)
}
