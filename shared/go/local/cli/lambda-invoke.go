package main

import (
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda/messages"

	"github.com/cshep4/kripto/shared/go/local"
)

func main() {
	p := flag.Int("p", 8001, "Port of running lambda.")
	d := flag.String("d", "", "request payload.")

	flag.Parse()

	port := *p
	payload := *d

	var payloadData map[string]interface{}
	if payload != "" {
		if err := json.Unmarshal([]byte(payload), &payloadData); err != nil {
			log.Println(err)
			return
		}
	}

	t := time.Now().AddDate(1, 0, 0)

	resp, err := local.Invoke(local.Input{
		Port:    port,
		Payload: payloadData,
		Deadline: &messages.InvokeRequest_Timestamp{
			Seconds: t.Unix(),
			Nanos:   int64(t.Nanosecond()),
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(resp))
}
