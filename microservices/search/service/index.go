package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/nats-io/nats.go"
	placeDTO "newnew.media/microservices/place/dto"
)

type IndexService struct {
	elasticClient *elasticsearch.Client
	natsClient    *nats.Conn
}

func NewIndexService(elasticClient *elasticsearch.Client, natsClient *nats.Conn) *IndexService {
	return &IndexService{elasticClient: elasticClient, natsClient: natsClient}
}

// func (is *IndexService) CreateIndex(indexName string, mapping []byte) error {
// 	// indexName := c.Params("indexName")
// 	// mapping := `
// 	// {
// 	//   "settings": {
// 	// 	"number_of_shards": 1,
// 	// 	"number_of_replicas": 1
// 	//   },
// 	//   "mappings": {
// 	// 	"properties": {
// 	// 	  "name": {
// 	// 		"type": "text"
// 	// 	  },
// 	// 	  "location": {
// 	// 		"type": "geo_point"
// 	// 	  }
// 	// 	}
// 	//   }
// 	// }`

// 	// Create the index
// 	res, err := is.elasticClient.Indices.Create(
// 		indexName,
// 		is.elasticClient.Indices.Create.WithBody(bytes.NewReader(mapping)),
// 	)

// 	if err != nil {
// 		log.Fatalf("Error creating index: %s", err)
// 		return err
// 	}

// 	if res.IsError() {
// 		return err
// 	}

// 	return nil
// }

func (is *IndexService) CreateIndex(indexName string, mapping []byte) error {

	// Create the index with the updated mapping
	res, err := is.elasticClient.Indices.Create(
		indexName,
		is.elasticClient.Indices.Create.WithBody(bytes.NewReader(mapping)),
	)

	if err != nil {
		log.Fatalf("Error creating index: %s", err)
		return err
	}

	if res.IsError() {
		return err
	}

	return nil
}

func (is *IndexService) DeleteIndex(indexName string) error {
	res, err := is.elasticClient.Indices.Delete([]string{indexName})

	if err != nil {
		log.Printf("Error deleting index: %s", err)
		return err
	}

	if res.IsError() {
		log.Printf("Error deleting index: %s", res.String())
		return err
	}

	return nil
}

func (is *IndexService) IndexPlace(place placeDTO.Place) error {
	// Index the place document into the "places" index
	placeJSON, err := json.Marshal(place)
	if err != nil {
		log.Printf("Error marshaling place to JSON: %s", err)
		// return c.Status(500).SendString("Error processing request")
		return err
	}

	res, err := is.elasticClient.Index(
		"places",
		strings.NewReader(string(placeJSON)),
	)
	if err != nil {
		log.Printf("Error indexing place: %s", err)
		return err
	}
	if res.IsError() {
		log.Printf("Error indexing place: %s", res.String())
		return err
	}

	// Index the location document into the "locations" index
	locationJSON, err := json.Marshal(place.Location)
	if err != nil {
		log.Printf("Error marshaling location to JSON: %s", err)
		return err
	}

	res, err = is.elasticClient.Index(
		"locations",
		strings.NewReader(string(locationJSON)),
	)
	if err != nil {
		log.Printf("Error indexing location: %s", err)
		return err
	}
	if res.IsError() {
		log.Printf("Error indexing location: %s", res.String())
		return err
	}

	return nil
}

func (is *IndexService) SubscribeToPlaceCreatedSubject() error {
	// Subscribe to place-created subject
	_, err := is.natsClient.Subscribe("place-created", func(msg *nats.Msg) {
		// Deserialize the JSON message to a placeDTO.Place struct
		var place placeDTO.Place
		err := json.Unmarshal(msg.Data, &place)
		if err != nil {
			log.Printf("Failed to deserialize place from JSON: %v", err)
			return
		}

		// Now you can use the place struct
		err = is.IndexPlace(place)
		if err != nil {
			log.Printf("Failed to index new place to elastic: %v", err)
			return
		}
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to place-created subject: %v", err)
	}
	// defer sub.Unsubscribe()
	return nil
}
