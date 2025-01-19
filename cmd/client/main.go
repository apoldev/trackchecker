package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	trackingService "github.com/apoldev/trackchecker/internal/app/grpcservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr     = flag.String("addr", "localhost:7766", "the address to connect to")
	fileName = flag.String("file", "./cmd/client/tracking_numbers.txt", "the filepath to read tracking numbers from")
)

type Event struct {
	Status  string `json:"status,omitempty"`
	Date    string `json:"date,omitempty"`
	Details string `json:"details,omitempty"`
	Place   string `json:"place,omitempty"`
}
type Result struct {
	Events      []Event `json:"events,omitempty"`
	CountryTo   string  `json:"CountryTo,omitempty"`
	CountryFrom string  `json:"CountryFrom,omitempty"`
}

func main() {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	data, err := os.ReadFile(*fileName)
	if err != nil {
		log.Panicf("could not read file: %v", err)
	}

	// Read tracking numbers from file
	trackingNumbers := strings.Split(strings.TrimSpace(string(data)), "\n")

	// Create a new grpc client
	c := trackingService.NewTrackingClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Get request with tracking numbers
	r, err := c.PostTracking(ctx, &trackingService.PostTrackingRequest{
		TrackingNumbers: trackingNumbers,
	})
	if err != nil {
		log.Fatalf("could not make request: %v", err)
	}
	reqID := r.GetTrackingId()

	results := make(map[string]struct{})
	mu := sync.Mutex{}
	ch := make(chan *trackingService.TrackResult)

	// Start a goroutine to get the results
	go func() {
		for {
			result, getErr := c.GetResult(context.Background(), &trackingService.GetTrackingID{
				Id: reqID,
			})
			if getErr != nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			mu.Lock()
			for _, trackRes := range result.GetTracking() {
				for _, spiderRes := range trackRes.GetResult() {
					mapKey := trackRes.GetCode() + "." + spiderRes.GetSpider()
					if _, ok := results[mapKey]; !ok {
						ch <- spiderRes
					}
					results[mapKey] = struct{}{}
				}
			}
			mu.Unlock()
			if len(result.GetTracking()) == len(trackingNumbers) {
				close(ch)
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

	for spiderRes := range ch {
		var event Event
		var oneResult Result
		err := json.Unmarshal([]byte(spiderRes.GetResult()), &oneResult)
		if err != nil {
			continue
		}
		if spiderRes.GetError() != "" {
			fmt.Printf("%s not found at %s: %s\n", spiderRes.GetTrackingNumber(), spiderRes.GetSpider(), spiderRes.GetError())
			continue
		}
		if len(oneResult.Events) > 0 {
			event = oneResult.Events[0]
		}
		status := strings.Trim(event.Date+", "+event.Status+
			", "+event.Place+", "+event.Details, ", .;")
		fmt.Printf("%s found at %s: %s\n", spiderRes.GetTrackingNumber(), spiderRes.GetSpider(), status)
	}
}
