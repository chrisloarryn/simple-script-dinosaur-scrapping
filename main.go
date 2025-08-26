package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"simple-script-dino/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/getAllDinoList", handlers.GetAllDinoList)
	http.HandleFunc("/getDinoDataByName", handlers.GetDinoDataByName)
	http.HandleFunc("/getAllDinoListWithDetails", handlers.GetAllDinoListWithDetails)

	fmt.Printf("Server listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
