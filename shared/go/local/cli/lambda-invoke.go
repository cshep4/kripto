package main

import (
	"encoding/json"
	"flag"
	"log"

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

	resp, err := local.Invoke(local.Input{
		Port:    port,
		Payload: payloadData,
	})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(resp))
}
