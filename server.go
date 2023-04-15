package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	// Default to sample stock data found in data.csv
	arg := "data.csv"

	// Check if a command-line argument was provided
	if len(os.Args) == 2 {
		if os.Args[1] == "green" {
			arg = "datagreen.csv"
		} else if os.Args[1] == "red" {
			arg = "datared.csv"
		}
	}
	// Parse the HTML template
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		panic(err)
	}

	// Create index.html file from template
	output, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	defer output.Close()

	// Execute the template with the command-line argument
	err = tmpl.Execute(output, struct{ Arg string }{Arg: arg})
	if err != nil {
		panic(err)
	}

    http.Handle("/", http.FileServer(http.Dir(".")))
    fmt.Println("Serving files from the current directory on http://localhost:8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
