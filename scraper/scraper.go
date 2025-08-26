package scraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"simple-script-dino/utils"

	"github.com/PuerkitoBio/goquery"
)

func httpGet(url string) (string, error) {
	resp, err := http.Get("https://www.nhm.ac.uk/discover/dino-directory" + url)
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

// GetOnlyDinoListConcurrent scrapes the list of dinosaurs concurrently.
func GetOnlyDinoListConcurrent() ([]map[string]string, error) {
	data, err := httpGet("/name/name-az-all.html")
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
			name := utils.CleanText(s.Find("a .dinosaurfilter--name-unhyphenated").Text())
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

// GetDinoByName scrapes dino detail page by name.
func GetDinoByName(name string) (map[string]string, error) {
	data, err := httpGet(fmt.Sprintf("/%s.html", strings.ToLower(name)))
	if err != nil {
		log.Println(err)
		return map[string]string{}, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		log.Println(err)
		return map[string]string{}, err
	}

	meaning := utils.CleanText(doc.Find(".dinosaur--meaning").Text())
	if len(meaning) > 0 {
		meaning = meaning[1:]
	} else {
		meaning = "N/A"
	}

	dinoData := map[string]string{
		"name":          utils.CleanText(doc.Find(".dinosaur--name-unhyphenated").Text()),
		"pronunciation": utils.CleanText(doc.Find(".dinosaur--pronunciation").Text()),
		"meaning":       meaning,
		"picture":       doc.Find(".dinosaur--image").AttrOr("src", ""),
		"content":       utils.CleanText(doc.Find(".dinosaur--content-container p").Text()),
	}

	doc.Find(".dinosaur--list dt").Each(func(_ int, sel *goquery.Selection) {
		key := strings.TrimRight(utils.CleanText(sel.Text()), ":")
		val := utils.CleanText(sel.Next().Text())
		dinoData[utils.CamelCase(key)] = val
	})

	return dinoData, nil
}
