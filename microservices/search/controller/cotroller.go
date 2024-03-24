package controller

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"newnew.media/microservices/place/dto"
)

type ElasticIndexController struct {
	elasticClient *elasticsearch.Client
}

func NewElasticIndexController(elasticClient *elasticsearch.Client) *ElasticIndexController {
	return &ElasticIndexController{elasticClient: elasticClient}
}

func (eic *ElasticIndexController) CreateIndex(c *fiber.Ctx) error {
	indexName := c.Params("indexName")
	mapping := `
	{
	  "settings": {
		"number_of_shards": 1,
		"number_of_replicas": 1
	  },
	  "mappings": {
		"properties": {
		  "name": {
			"type": "text"
		  },
		  "location": {
			"type": "geo_point"
		  }
		}
	  }
	}`

	// Create the index
	res, err := eic.elasticClient.Indices.Create(
		indexName,
		eic.elasticClient.Indices.Create.WithBody(strings.NewReader(mapping)),
	)

	if err != nil {
		log.Fatalf("Error creating index: %s", err)
		return c.Status(500).SendString("Error creating index")
	}

	if res.IsError() {
		return c.Status(500).SendString("Error creating index")
	}

	return c.SendString("Index created successfully")
}

func (eic *ElasticIndexController) DeleteIndex(c *fiber.Ctx) error {
	indexName := c.Params("indexName")
	res, err := eic.elasticClient.Indices.Delete([]string{indexName})

	if err != nil {
		log.Printf("Error deleting index: %s", err)
		return c.Status(500).SendString("Error deleting index")
	}

	if res.IsError() {
		log.Printf("Error deleting index: %s", res.String())
		return c.Status(500).SendString("Error deleting index")
	}

	return c.SendString("Index deleted successfully")
}

func (eic *ElasticIndexController) IndexPlace(c *fiber.Ctx) error {
	// Parse the JSON request into a Place struct
	var place dto.Place
	if err := c.BodyParser(&place); err != nil {
		log.Printf("Error parsing request body: %s", err)
		return c.Status(400).SendString("Error parsing request body")
	}

	// Index the place document into the "places" index
	placeJSON, err := json.Marshal(place)
	if err != nil {
		log.Printf("Error marshaling place to JSON: %s", err)
		return c.Status(500).SendString("Error processing request")
	}

	res, err := eic.elasticClient.Index(
		"places",
		strings.NewReader(string(placeJSON)),
	)
	if err != nil {
		log.Printf("Error indexing place: %s", err)
		return c.Status(500).SendString("Error indexing place")
	}
	if res.IsError() {
		log.Printf("Error indexing place: %s", res.String())
		return c.Status(500).SendString("Error indexing place")
	}

	// Index the location document into the "locations" index
	locationJSON, err := json.Marshal(place.Location)
	if err != nil {
		log.Printf("Error marshaling location to JSON: %s", err)
		return c.Status(500).SendString("Error processing request")
	}

	res, err = eic.elasticClient.Index(
		"locations",
		strings.NewReader(string(locationJSON)),
	)
	if err != nil {
		log.Printf("Error indexing location: %s", err)
		return c.Status(500).SendString("Error indexing location")
	}
	if res.IsError() {
		log.Printf("Error indexing location: %s", res.String())
		return c.Status(500).SendString("Error indexing location")
	}

	return c.SendString("Place and location indexed successfully")
}

func (eic *ElasticIndexController) SearchPlaceByName(c *fiber.Ctx) error {
	// Extract the name query parameter from the request
	name := c.Query("name")
	if name == "" {
		return c.Status(400).SendString("Name query parameter is required")
	}

	// Prepare the search query
	var buf strings.Builder
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": name,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("Error encoding search query: %s", err)
		return c.Status(500).SendString("Error processing request")
	}

	// Convert strings.Builder to strings.Reader
	queryReader := strings.NewReader(buf.String())

	// Execute the search query
	res, err := eic.elasticClient.Search(
		eic.elasticClient.Search.WithContext(context.Background()),
		eic.elasticClient.Search.WithIndex("places"),
		eic.elasticClient.Search.WithBody(queryReader),
		eic.elasticClient.Search.WithTrackTotalHits(true),
		eic.elasticClient.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("Error executing search: %s", err)
		return c.Status(500).SendString("Error executing search")
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error executing search: %s", res.String())
		return c.Status(500).SendString("Error executing search")
	}

	// Decode the search response
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error decoding search response: %s", err)
		return c.Status(500).SendString("Error processing request")
	}

	// Extract the hits from the search response
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})

	// Prepare the response
	var places []dto.Place
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		var place dto.Place
		if err := mapstructure.Decode(source, &place); err != nil {
			log.Printf("Error decoding place: %s", err)
			return c.Status(500).SendString("Error processing request")
		}
		places = append(places, place)
	}

	// Return the search results
	return c.JSON(places)
}
