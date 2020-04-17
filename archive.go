package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type mainData struct {
	Count int64  `json:"count"`
	Date  string `json:"date"`
	Delta int64  `json:"delta"`
}

//Archive is a structure with all data for grubber
type Archive struct {
	AbsoluteDelta int64      `json:"absolute_delta"`
	Description   string     `json:"description"`
	SourceURL     string     `json:"source_url"`
	Data          []mainData `json:"data"`
}

type archiveDataMap map[string]interface{}

func castArchiveData(data archiveDataMap) Archive {

	absoluteDelta := int64(data["absolute_delta"].(float64))
	description := data["description"].(string)
	sourceURL := data["source_url"].(string)

	ddata := data["data"].([]interface{})
	var dddata []mainData

	for _, element := range ddata {
		el := element.(map[string]interface{})
		count := int64(el["count"].(float64))
		date := el["date"].(string)
		delta := int64(el["delta"].(float64))
		dddata = append(dddata, mainData{count, date, delta})
	}

	return Archive{absoluteDelta, description, sourceURL, dddata}
}

func readArchive(fileName string) Archive {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}
	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}

	var archiveData archiveDataMap
	if err := json.Unmarshal(data, &archiveData); err != nil {
		panic(fmt.Errorf("%s", err))
	}

	return castArchiveData(archiveData)
}

func (arch *Archive) addData(newCount int64) {
	firstData := arch.Data[0]
	lastData := arch.Data[len(arch.Data)-1]
	delta := newCount - lastData.Count

	date := fmt.Sprint(time.Now().Format("2006-01-02"))
	arch.AbsoluteDelta = newCount - firstData.Count
	arch.Data = append(arch.Data, mainData{newCount, date, delta})
}

func (arch *Archive) getLastData() mainData {
	return arch.Data[len(arch.Data)-1]
}

func (arch *Archive) saveData(resultsFileName string) {
	// println(resultsFileName)
	f, err := os.Create(resultsFileName)
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}
	defer f.Close()

	f.WriteString(arch.convertToJSON())
}

func (arch *Archive) convertToJSON() string {
	data, err := json.MarshalIndent(arch, "", "    ")
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}
	return string(data)
}
