package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
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

// stripHTML removes HTML tags from a string (best-effort, minimal).
var htmlTagRe = regexp.MustCompile("<[^>]*>")

func stripHTML(s string) string {
	return strings.TrimSpace(htmlTagRe.ReplaceAllString(s, ""))
}

// GetOnlyDinoListConcurrent scrapes the list of dinosaurs concurrently.
func GetOnlyDinoListConcurrent() ([]map[string]string, error) {
	// The old A–Z page no longer exists. Use the landing page and collect links to detail pages.
	data, err := httpGet("")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Collect unique links that match /discover/dino-directory/<slug>.html
	linkSet := make(map[string]struct{})
	var links []string
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}
		if strings.HasPrefix(href, "/discover/dino-directory/") && strings.HasSuffix(href, ".html") { // absolute path to a dinosaur page
			// ensure it's exactly one segment after base
			// e.g., /discover/dino-directory/triceratops.html
			parts := strings.Split(strings.TrimPrefix(href, "/discover/dino-directory/"), "/")
			if len(parts) == 1 && strings.HasSuffix(parts[0], ".html") {
				if _, ok := linkSet[href]; !ok {
					linkSet[href] = struct{}{}
					links = append(links, href)
				}
			}
		}
	})

	dinoChan := make(chan map[string]string, len(links))
	var wg sync.WaitGroup

	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			slug := strings.TrimSuffix(strings.TrimPrefix(link, "/discover/dino-directory/"), ".html")
			nameParts := strings.Split(strings.ReplaceAll(slug, "-", " "), " ")
			for i, p := range nameParts {
				if len(p) > 0 {
					nameParts[i] = strings.Title(p)
				}
			}
			name := strings.Join(nameParts, " ")
			dinoChan <- map[string]string{"name": name, "link": link}
		}(link)
	}

	wg.Wait()
	close(dinoChan)

	dinoList := []map[string]string{}
	for d := range dinoChan {
		dinoList = append(dinoList, d)
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

	// Extract __NEXT_DATA__ JSON
	jsonText := strings.TrimSpace(doc.Find("script#__NEXT_DATA__").Text())
	if jsonText == "" {
		// Fallback: return minimal data with name
		return map[string]string{
			"name":          name,
			"pronunciation": "",
			"meaning":       "N/A",
			"picture":       "",
			"content":       "",
		}, nil
	}

	// Minimal struct to parse needed fields
	type Media struct {
		Identifier    string `json:"identifier"`
		MediaTypeName string `json:"mediaTypeName"`
		MediaTypePath string `json:"mediaTypePath"`
	}
	type BodyShape struct {
		BodyShape string `json:"bodyShape"`
	}
	type Period struct {
		Period string `json:"period"`
	}
	type TextBlock struct {
		TextBlock string `json:"textBlock"`
	}
	type Dino struct {
		Genus               string      `json:"genus"`
		NamePronounciation  string      `json:"namePronounciation"`
		NameMeaning         string      `json:"nameMeaning"`
		DietTypeName        string      `json:"dietTypeName"`
		BodyShape           *BodyShape  `json:"bodyShape"`
		Period              *Period     `json:"period"`
		LengthFrom          *float64    `json:"lengthFrom"`
		LengthTo            *float64    `json:"lengthTo"`
		MassFrom            *float64    `json:"massFrom"`
		MassTo              *float64    `json:"massTo"`
		MediaCollection     []Media     `json:"mediaCollection"`
		TextBlockCollection []TextBlock `json:"textBlockCollection"`
	}
	type PageProps struct {
		Dinosaur             *Dino   `json:"dinosaur"`
		DinosaurErrorMessage *string `json:"dinosaurErrorMessage"`
	}
	type NextData struct {
		Props struct {
			PageProps PageProps `json:"pageProps"`
		} `json:"props"`
	}

	var nd NextData
	if err := json.Unmarshal([]byte(jsonText), &nd); err != nil {
		log.Println(err)
		return map[string]string{}, err
	}

	d := nd.Props.PageProps.Dinosaur
	if d == nil {
		return map[string]string{}, fmt.Errorf("dinosaur data not found")
	}

	// Build content from text blocks
	var contentParts []string
	for _, tb := range d.TextBlockCollection {
		if tb.TextBlock != "" {
			contentParts = append(contentParts, stripHTML(tb.TextBlock))
		}
	}
	content := strings.TrimSpace(strings.Join(contentParts, " \n "))

	// Pick picture if any (constructing full URL reliably is non-trivial; leave empty if unknown)
	picture := ""

	// Map details similar to old keys
	dinoData := map[string]string{
		"name":          utils.CleanText(d.Genus),
		"pronunciation": utils.CleanText(d.NamePronounciation),
		"meaning":       utils.CleanText(d.NameMeaning),
		"picture":       picture,
		"content":       utils.CleanText(content),
	}

	// Additional attributes akin to the previous list
	if d.BodyShape != nil && d.BodyShape.BodyShape != "" {
		dinoData["typeOfDinosaur"] = d.BodyShape.BodyShape
	}
	if d.Period != nil && d.Period.Period != "" {
		dinoData["period"] = d.Period.Period
	}
	if d.DietTypeName != "" {
		dinoData["diet"] = d.DietTypeName
	}
	if d.LengthFrom != nil {
		lf := fmt.Sprintf("%.1f", *d.LengthFrom)
		if d.LengthTo != nil {
			lt := fmt.Sprintf("%.1f", *d.LengthTo)
			dinoData["length"] = lf + "–" + lt + " m"
		} else {
			dinoData["length"] = lf + " m"
		}
	}
	if d.MassFrom != nil {
		mf := fmt.Sprintf("%.0f", *d.MassFrom)
		if d.MassTo != nil {
			mt := fmt.Sprintf("%.0f", *d.MassTo)
			dinoData["weight"] = mf + "–" + mt + " kg"
		} else {
			dinoData["weight"] = mf + " kg"
		}
	}

	return dinoData, nil
}
