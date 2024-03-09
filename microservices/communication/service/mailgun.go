package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/mailgun/mailgun-go/v4/events"
)

var domain string = "stateful.art" // e.g. mg.yourcompany.com

// You can find the Private API Key in your Account Menu, under "Settings":
// (https://app.mailgun.com/app/account/security)
var key string = os.Getenv("MAILGUN_APIKEY")

// func SendEmail() {

// 	// Create an instance of the Mailgun Client
// 	mg := mailgun.NewMailgun(domain, key)

// 	//When you have an EU-domain, you must specify the endpoint:
// 	mg.SetAPIBase("https://api.eu.mailgun.net/v3")

// 	sender := "contact@stateful.art"
// 	subject := "Fancy subject!"
// 	body := "Hello from Mailgun Go!"
// 	recipient := "abbas.tolgay.yilmaz@gmail.com"

// 	// The message object allows you to add attachments and Bcc recipients
// 	message := mg.NewMessage(sender, subject, body, recipient)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()

// 	// Send the message with a 10 second timeout
// 	resp, id, err := mg.Send(ctx, message)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("ID: %s Resp: %s\n", id, resp)

// }

func Events() {
	mg := mailgun.NewMailgun("your-domain.com", "your-private-key")

	it := mg.ListEvents(&mailgun.ListEventOptions{Limit: 100})

	var page []mailgun.Event

	// The entire operation should not take longer than 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	for it.Next(ctx, &page) {
		for _, e := range page {
			// You can access some fields via the interface
			fmt.Printf("Event: '%s' TimeStamp: '%s'\n", e.GetName(), e.GetTimestamp())

			// and you can act upon each event by type
			switch event := e.(type) {
			case *events.Accepted:
				fmt.Printf("Accepted: auth: %t\n", event.Flags.IsAuthenticated)
			case *events.Delivered:
				fmt.Printf("Delivered transport: %s\n", event.Envelope.Transport)
			case *events.Failed:
				fmt.Printf("Failed reason: %s\n", event.Reason)
			case *events.Clicked:
				fmt.Printf("Clicked GeoLocation: %s\n", event.GeoLocation.Country)
			case *events.Opened:
				fmt.Printf("Opened GeoLocation: %s\n", event.GeoLocation.Country)
			case *events.Rejected:
				fmt.Printf("Rejected reason: %s\n", event.Reject.Reason)
			case *events.Stored:
				fmt.Printf("Stored URL: %s\n", event.Storage.URL)
			case *events.Unsubscribed:
				fmt.Printf("Unsubscribed client OS: %s\n", event.ClientInfo.ClientOS)
			}
		}

	}
}
