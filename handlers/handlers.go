package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"simple-script-dino/scraper"
)

func GetAllDinoList(res http.ResponseWriter, req *http.Request) {
	dinoList, err := scraper.GetOnlyDinoListConcurrent()
	if err != nil {
		log.Println(err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resJSON, err := json.Marshal(map[string]interface{}{
		"results": len(dinoList),
		"data":    dinoList,
	})
	if err != nil {
		log.Println(err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resJSON)
}

// GetAllDinoListWithDetails uses goroutines to fetch all dino details concurrently.
func GetAllDinoListWithDetails(res http.ResponseWriter, req *http.Request) {
	var noDataDinosaurs int16
	dinoList, err := scraper.GetOnlyDinoListConcurrent()
	if err != nil {
		log.Println(err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	wg := sync.WaitGroup{}
	dinoChan := make(chan map[string]string, len(dinoList))
	for _, dino := range dinoList {
		wg.Add(1)
		go func(dino map[string]string) {
			defer wg.Done()
			dinoData, err := scraper.GetDinoByName(dino["name"])

			if dinoData["meaning"] == "N/A" {
				noDataDinosaurs++
			}
			if err != nil {
				log.Println(err)
				return
			}
			dinoData["name"] = dino["name"]
			dinoData["link"] = dino["link"]
			dinoChan <- dinoData
		}(dino)
	}
	wg.Wait()
	close(dinoChan)

	processedDinoList := []map[string]string{}
	for dinoData := range dinoChan {
		processedDinoList = append(processedDinoList, dinoData)
	}

	resJSON, err := json.Marshal(map[string]interface{}{
		"results": len(processedDinoList),
		"data":    processedDinoList,
		"noData":  noDataDinosaurs,
	})
	if err != nil {
		log.Println(err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resJSON)
}

func GetDinoDataByName(res http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")

	dinoData, err := scraper.GetDinoByName(name)

	resJSON, err := json.Marshal(dinoData)
	if err != nil {
		log.Println(err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resJSON)
}
