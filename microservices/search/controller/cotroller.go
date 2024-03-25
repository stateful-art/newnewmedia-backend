package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	placeDTO "newnew.media/microservices/place/dto"
	service "newnew.media/microservices/search/service"
)

type SearchIndexController struct {
	searchService *service.SearchService
	indexService  *service.IndexService
}

func NewSearchIndexController(searchService *service.SearchService, indexService *service.IndexService) *SearchIndexController {
	return &SearchIndexController{searchService: searchService, indexService: indexService}
}

// func (sic *SearchIndexController) CreateIndex(c *fiber.Ctx) error {
// 	indexName := c.Params("indexName")
// 	mapping := `
// 	{
// 	  "settings": {
// 		"number_of_shards": 1,
// 		"number_of_replicas": 1
// 	  },
// 	  "mappings": {
// 		"properties": {
// 		  "name": {
// 			"type": "text"
// 		  },
// 		  "location": {
// 			"type": "geo_point"
// 		  }
// 		}
// 	  }
// 	}`

// 	// Create the index
// 	res, err := sic.elasticClient.Indices.Create(
// 		indexName,
// 		sic.elasticClient.Indices.Create.WithBody(strings.NewReader(mapping)),
// 	)

// 	if err != nil {
// 		log.Fatalf("Error creating index: %s", err)
// 		return c.Status(500).SendString("Error creating index")
// 	}

// 	if res.IsError() {
// 		return c.Status(500).SendString("Error creating index")
// 	}

// 	return c.SendString("Index created successfully")
// }

// func (sic *SearchIndexController) DeleteIndex(c *fiber.Ctx) error {
// 	indexName := c.Params("indexName")
// 	res, err := sic.elasticClient.Indices.Delete([]string{indexName})

// 	if err != nil {
// 		log.Printf("Error deleting index: %s", err)
// 		return c.Status(500).SendString("Error deleting index")
// 	}

// 	if res.IsError() {
// 		log.Printf("Error deleting index: %s", res.String())
// 		return c.Status(500).SendString("Error deleting index")
// 	}

// 	return c.SendString("Index deleted successfully")
// }

func (sic *SearchIndexController) IndexPlace(c *fiber.Ctx) error {
	// Parse the JSON request into a Place struct
	var place placeDTO.Place
	if err := c.BodyParser(&place); err != nil {
		log.Printf("Error parsing request body: %s", err)
		return c.Status(400).SendString("Error parsing request body")
	}
	err := sic.indexService.IndexPlace(place)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendString("Place and location indexed successfully")
}

// func (sic *SearchIndexController) SearchPlaceByName(c *fiber.Ctx) error {
// 	// Extract the name query parameter from the request
// 	name := c.Query("name")
// 	if name == "" {
// 		return c.Status(400).SendString("Name query parameter is required")
// 	}

// 	// Prepare the search query
// 	var buf strings.Builder
// 	query := map[string]interface{}{
// 		"query": map[string]interface{}{
// 			"match": map[string]interface{}{
// 				"name": name,
// 			},
// 		},
// 	}
// 	if err := json.NewEncoder(&buf).Encode(query); err != nil {
// 		log.Printf("Error encoding search query: %s", err)
// 		return c.Status(500).SendString("Error processing request")
// 	}

// 	// Convert strings.Builder to strings.Reader
// 	queryReader := strings.NewReader(buf.String())

// 	// Execute the search query
// 	res, err := sic.elasticClient.Search(
// 		sic.elasticClient.Search.WithContext(context.Background()),
// 		sic.elasticClient.Search.WithIndex("places"),
// 		sic.elasticClient.Search.WithBody(queryReader),
// 		sic.elasticClient.Search.WithTrackTotalHits(true),
// 		sic.elasticClient.Search.WithPretty(),
// 	)
// 	if err != nil {
// 		log.Printf("Error executing search: %s", err)
// 		return c.Status(500).SendString("Error executing search")
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		log.Printf("Error executing search: %s", res.String())
// 		return c.Status(500).SendString("Error executing search")
// 	}

// 	// Decode the search response
// 	var r map[string]interface{}
// 	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
// 		log.Printf("Error decoding search response: %s", err)
// 		return c.Status(500).SendString("Error processing request")
// 	}

// 	// Extract the hits from the search response
// 	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})

// 	// Prepare the response
// 	var places []dto.Place
// 	for _, hit := range hits {
// 		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
// 		var place dto.Place
// 		if err := mapstructure.Decode(source, &place); err != nil {
// 			log.Printf("Error decoding place: %s", err)
// 			return c.Status(500).SendString("Error processing request")
// 		}
// 		places = append(places, place)
// 	}

// 	// Return the search results
// 	return c.JSON(places)
// }
