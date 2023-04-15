# heatmap
Stock market price activity heatmap. But really just an excuse to learn Golang for no reason

## Mixed Day
![Mixed Day](images/myport_lg.JPG)

## Good Day
![Good Day](images/myportgreen_lg.JPG)

## Bad Day
![Bad Day](images/myportred_lg.JPG)

### Prerequisites
1. Go
2. Polygon.io API key stored as `POLYGON_API_KEY` in `.env`
3. The imports in `getportfolio.go`

### Usage
1. go run `getportfolio.go`
2. go run `server.go`
* Defaults to loading `data.csv`
* Optional command line arguments include `green` and `red`, ie, `go run server.go red` to simulate a red day
3. Navigate to `http://localhost:8000/` in a browser window

### Notes
`getportfolio.go` leverages the free version of the Polygon.io API and thus data is limited to previous day.  The color of each square in the Heatmap depends on the previous day difference between the close price and the open price.  The rectangle will approach bright red the more negative it closed; conversely, the rectangle will approach bright green the more positive it closed.  
I did very little testing on this but I imagine it works best with 15-30 stock tickers.  You can modify `myport.csv` to reflect stock tickers followed by the number of shares you own like this:
* msft,7  
* NVDA,13  

Due to Polygon.io API limitations, the script will sleep 12 seconds in between API calls if you have more than 5 items in `myport.csv`  

And lastly, the stocks mentioned in this project were randomly chosen and some of the data falsified for demonstrative purposes.