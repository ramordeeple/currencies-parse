package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shopspring/decimal"
)

type envelope struct {
	ResultCode string
	Payload    payload
}

type payload struct {
	LastUpdate timestamp
	Rates      []rate
}

type timestamp struct {
	Milliseconds int64
}

type rate struct {
	Category     string
	FromCurrency currency
	ToCurrency   currency
	Buy          decimal.Decimal
	Sell         decimal.Decimal
}

type currency struct {
	Name string
}

func main() {
	httpResp, httpErr := http.Get("https://api.tinkoff.ru/v1/currency_rates")
	checkErr(httpErr)

	httpBody := httpResp.Body
	defer httpBody.Close()

	envelop := new(envelope)

	decoder := json.NewDecoder(httpBody)
	decodeErr := decoder.Decode(envelop)
	checkErr(decodeErr)

	if envelop.ResultCode != "OK" {
		log.Fatal("Could not load: ", envelop.ResultCode)
	}

	payload := envelop.Payload

	actualAt := time.Unix(payload.LastUpdate.Milliseconds/1000, 0)

	for _, rate := range payload.Rates {
		if rate.Category != "DebitCardsTransfers" {
			continue
		}

		fmt.Printf("%s | %-8s %8s |%8s\n",
			actualAt.Format("15:04:05"),
			rate.FromCurrency.Name+rate.ToCurrency.Name,
			rate.Buy.StringFixed(2),
			rate.Sell.StringFixed(2))
	}
	fmt.Println("\nPress any key to exit")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
