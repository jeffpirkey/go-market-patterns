package main

import "market-patterns/model"

/*
Requirements:
1) Lookup periods by ticker symbol
2) Sort periods by date
3) Build patterns for a ticker symbol using Up, NoChange, and Down
4) Count the number of Up, NoChange, and Down results for a pattern over time for a ticker symbol
5) Find a pattern for a ticker symbol
6) Find the most current period by date

*/

var Tickers = model.NewTickers() // Maps ticker symbol to a map of Result patterns and Pattern

func main() {

}
