package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	routeURLKey = "RouteURL"
)

type Message struct {
	Message string `json:"message"`
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)

		routeURL := os.Getenv(routeURLKey)

		m := Message{}
		err := json.Unmarshal([]byte(message.Body), &m)
		if err != nil {
			return errors.New("error unmarshalling message body")
		}
		err = routeMessage(m, routeURL)
		if err != nil {
			return errors.New("Error routing message")
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}

func routeMessage(message Message, url string) error {

	out, err := json.Marshal(message)

	if err != nil {
		return err
	}
	jsonString := string(out)
	log.Println("about to send message ", jsonString, " to url ", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonString)))

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Println("ERROR: error making http call ", err)
		return errors.New("error making http call")
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)

	if is2xx(&resp.StatusCode) {
		return nil
	}

	log.Println("ERROR: returned response ", resp)
	return errors.New("Client returned http status code")
}

func is2xx(status *int) bool {
	switch *status {
	case 200:
		return true
	case 201:
		return true
	case 202:
		return true
	default:
		return false
	}
}
