package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	options := ntp.QueryOptions{Timeout: 10 * time.Second}

	response, err := ntp.QueryWithOptions("0.beevik-ntp.pool.ntp.org", options)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	curTime := response.Time.Add(response.ClockOffset)
	fmt.Printf("Текущее время: %s\n", curTime)
}
