package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

// Round function from https://gosamples.dev/round-float/
func roundFloat(val float64, precision uint) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val * ratio) / ratio
}

func serve(arg string) {
	if arg == "green" {
			arg = "datagreen.csv"
	} else if arg == "red" {
		arg = "datared.csv"
	} else {
		arg = "data.csv"
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

func main() {
	if len(os.Args) == 1 {
		// Struct to store ticker information
		type TickerInfo struct {
			TickerName string
			PercentageChange string
			ClosePrice float64
			CloseTotalValue float64
			PercentOfPort *float64
		}

		TickerList := make([]TickerInfo, 0)

		// Create Polygon API connection with API key in .env
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		apiKey := os.Getenv("POLYGON_API_KEY")

		c := polygon.New(apiKey)

		// Open the portfolio CSV file; for best results, fill myport.csv with 15-30 stock tickers and number of shares like: msft,17
		inFile, err := os.Open("myport.csv")
		if err != nil {
			panic(err)
		}
		defer inFile.Close()

		// Create a new CSV reader
		reader := csv.NewReader(inFile)

		// Count lines
		lineCount := 0
		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
		
			// Increment the line count for each line
			lineCount++
		}
		_, err = inFile.Seek(0, 0)
		if err != nil {
			log.Fatal(err)
		}

		var totalPortValue float64 = 0.00

		// Read the CSV file line by line
		for {
			record, err := reader.Read()
			if err != nil {
				// If we reach the end of the file, break out of the loop
				if err.Error() == "EOF" {
					break
				}
				panic(err)
			}
			
			// Set ticker name and add to parameters for API call
			ticker := strings.ToUpper(record[0])
			numShares, err := strconv.ParseFloat(record[1], 64)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			params := models.GetPreviousCloseAggParams{
				Ticker: ticker,
			}.WithAdjusted(true)

			// Get opening/close price of ticker
			res, err := c.GetPreviousCloseAgg(context.Background(), params)
			if err != nil {
				log.Fatal(err)
			}
			open := res.Results[0].Open
			close := res.Results[0].Close
			closeTotal := (close * numShares)
			percentageChange := ((close - open) / open) * 100
			formattedPercentageChange := fmt.Sprintf("%.2f", percentageChange)
			totalPortValue += closeTotal

			// Use info to create new struct and append to list
			TickerInstance := TickerInfo{
				TickerName: ticker,
				PercentageChange: formattedPercentageChange,
				ClosePrice: close,
				CloseTotalValue: roundFloat(closeTotal, 2),
				PercentOfPort: nil,
			}

			TickerList = append(TickerList, TickerInstance)

			if lineCount > 5 {
				time.Sleep(12 * time.Second)
			}
		}

		// Iterate through TickerList and calculate/set PercentOfPort
		for i := range TickerList {
			tickerPercentage := TickerList[i].CloseTotalValue / totalPortValue
			tickerString := fmt.Sprintf("%.2f", tickerPercentage)
			tickerFloat, err := strconv.ParseFloat(tickerString, 64)
			if err == nil {
				TickerList[i].PercentOfPort = &tickerFloat
			}

		}

		// Create/overwrite data.csv file and instantiate csv writer
		outFile, err := os.Create("data.csv")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer outFile.Close()

		// CSV Setup
		writer := csv.NewWriter(outFile)
		defer writer.Flush()
		writer.Comma = ','

		// Columns
		firstRow := []string{"ticker", "closeprice", "daychange","portweight"}
		err = writer.Write(firstRow)

		for _, tickerInfo := range TickerList {
			percentOfPortStr := ""
			if tickerInfo.PercentOfPort != nil {
				percentOfPortStr = strconv.FormatFloat(*tickerInfo.PercentOfPort, 'f', -1, 64)
			}

			csvLine := []string{tickerInfo.TickerName, strconv.FormatFloat(tickerInfo.ClosePrice, 'f', -1, 64), tickerInfo.PercentageChange, percentOfPortStr}
			err := writer.Write(csvLine)
			if err != nil {
				fmt.Println("Error writing record to CSV:", err)
				return
			}
		}

		fmt.Println("CSV file written successfully")
	} else if len(os.Args) == 2 {
		serve(os.Args[1])
	} 
}