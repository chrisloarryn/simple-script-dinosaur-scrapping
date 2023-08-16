package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func cleanText(input string) string {
	cleaned := strings.TrimSpace(strings.ReplaceAll(input, "\n", ""))

	// also clean \t
	cleaned = strings.TrimSpace(strings.ReplaceAll(cleaned, "\t", ""))

	if len(cleaned) > 0 {
		return cleaned
	}
	return input
}
func getOnlyDinoListConcurrent() ([]map[string]string, error) {
	data, err := get("/name/name-az-all.html")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	dinoChan := make(chan map[string]string)
	var wg sync.WaitGroup

	doc.Find(".dinosaurfilter--dinosaur").Each(func(index int, sel *goquery.Selection) {
		wg.Add(1)
		go func(s *goquery.Selection) {
			defer wg.Done()
			name := cleanText(s.Find("a .dinosaurfilter--name-unhyphenated").Text())
			link, _ := s.Find("a").Attr("href")
			dinoData := map[string]string{"name": name, "link": link}
			dinoChan <- dinoData
		}(sel)
	})

	go func() {
		wg.Wait()
		close(dinoChan)
	}()

	dinoList := []map[string]string{}
	for dinoData := range dinoChan {
		dinoList = append(dinoList, dinoData)
	}

	return dinoList, nil
}

func getAllDinoList(res http.ResponseWriter, req *http.Request) {
	dinoList, err := getOnlyDinoListConcurrent()
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

// use goroutine to get all dino data

func getAllDinoListWithDetails(res http.ResponseWriter, req *http.Request) {
	var noDataDinosaurs int16
	dinoList, err := getOnlyDinoListConcurrent()
	if err != nil {
		log.Println(err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	wg := sync.WaitGroup{}
	dinoChan := make(chan map[string]string, len(dinoList)) // Create a channel to receive processed dino data
	for _, dino := range dinoList {
		wg.Add(1)
		go func(dino map[string]string) {
			defer wg.Done()
			dinoData, err := getDinoByName(dino["name"])

			if dinoData["meaning"] == "N/A" {
				noDataDinosaurs++
			}
			if err != nil {
				log.Println(err)
				return
			}
			dinoData["name"] = dino["name"] // Store the dino name in the processed data
			dinoData["link"] = dino["link"] // Store the dino link in the processed data
			dinoChan <- dinoData
		}(dino)
	}
	wg.Wait()
	close(dinoChan) // Close the channel once all goroutines are done sending data

	// Collect data from the channel into a new dinoList
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

func getDinoByName(name string) (map[string]string, error) {
	data, err := get(fmt.Sprintf("/%s.html", strings.ToLower(name)))
	if err != nil {
		log.Println(err)
		return map[string]string{}, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		log.Println(err)
		return map[string]string{}, err
	}

	meaning := cleanText(doc.Find(".dinosaur--meaning").Text())
	if len(meaning) > 0 {
		meaning = meaning[1:]
	} else {
		meaning = "N/A"
	}

	dinoData := map[string]string{
		"name":          cleanText(doc.Find(".dinosaur--name-unhyphenated").Text()),
		"pronunciation": cleanText(doc.Find(".dinosaur--pronunciation").Text()),
		"meaning":       meaning,
		"picture":       doc.Find(".dinosaur--image").AttrOr("src", ""),
		"content":       cleanText(doc.Find(".dinosaur--content-container p").Text()),
	}

	doc.Find(".dinosaur--list dt").Each(func(_ int, sel *goquery.Selection) {
		key := strings.TrimRight(cleanText(sel.Text()), ":")
		val := cleanText(sel.Next().Text())
		dinoData[camelCase(key)] = val
	})

	return dinoData, nil
}

func getDinoDataByName(res http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")

	dinoData, err := getDinoByName(name)

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

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/getAllDinoList", getAllDinoList)
	http.HandleFunc("/getDinoDataByName", getDinoDataByName)
	http.HandleFunc("/getAllDinoListWithDetails", getAllDinoListWithDetails)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

	fmt.Println("Server is running on port 8080")
}

func get(url string) (string, error) {
	resp, err := http.Get("https://www.nhm.ac.uk/discover/dino-directory" + url) // Replace with your actual URL
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func camelCase(input string) string {
	parts := strings.Fields(input)
	for i, part := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(part)
		} else {
			parts[i] = strings.Title(strings.ToLower(part))
		}
	}
	return strings.Join(parts, "")
}
