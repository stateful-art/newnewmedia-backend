package service

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mitchellh/mapstructure"
	placeDTO "newnew.media/microservices/place/dto"
)

type SearchService struct {
	elasticClient *elasticsearch.Client
}

func NewSearchService(elasticClient *elasticsearch.Client) *SearchService {
	return &SearchService{elasticClient: elasticClient}
}

// func (ss *SearchService) SearchPlaceByName(name string) ([]placeDTO.Place, error) {
// 	// Extract the name query parameter from the request
// 	var places []placeDTO.Place

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
// 		return places, err
// 	}

// 	// Convert strings.Builder to strings.Reader
// 	queryReader := strings.NewReader(buf.String())

// 	// Execute the search query
// 	res, err := ss.elasticClient.Search(
// 		ss.elasticClient.Search.WithContext(context.Background()),
// 		ss.elasticClient.Search.WithIndex("places"),
// 		ss.elasticClient.Search.WithBody(queryReader),
// 		ss.elasticClient.Search.WithTrackTotalHits(true),
// 		ss.elasticClient.Search.WithPretty(),
// 	)
// 	if err != nil {
// 		log.Printf("Error executing search: %s", err)
// 		return places, err
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		log.Printf("Error executing search: %s", res.String())
// 		return places, err
// 	}

// 	// Decode the search response
// 	var r map[string]interface{}
// 	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
// 		log.Printf("Error decoding search response: %s", err)
// 		return places, err
// 	}

// 	// Extract the hits from the search response
// 	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})

// 	// Prepare the response
// 	for _, hit := range hits {
// 		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
// 		var place placeDTO.Place
// 		if err := mapstructure.Decode(source, &place); err != nil {
// 			log.Printf("Error decoding place: %s", err)
// 			return places, err
// 		}
// 		places = append(places, place)
// 	}

// 	// Return the search results
// 	return places, nil
// }

// func (ss *SearchService) SearchPlaceByName(name string) ([]placeDTO.Place, error) {
// 	var places []placeDTO.Place

// 	// Prepare the search query using a prefix query
// 	var buf strings.Builder
// 	query := map[string]interface{}{
// 		"query": map[string]interface{}{
// 			"prefix": map[string]interface{}{
// 				"name": map[string]interface{}{
// 					"value": name,
// 				},
// 			},
// 		},
// 	}
// 	if err := json.NewEncoder(&buf).Encode(query); err != nil {
// 		return places, err
// 	}

// 	// Convert strings.Builder to strings.Reader
// 	queryReader := strings.NewReader(buf.String())

// 	// Execute the search query
// 	res, err := ss.elasticClient.Search(
// 		ss.elasticClient.Search.WithContext(context.Background()),
// 		ss.elasticClient.Search.WithIndex("places"),
// 		ss.elasticClient.Search.WithBody(queryReader),
// 		ss.elasticClient.Search.WithTrackTotalHits(true),
// 		ss.elasticClient.Search.WithPretty(),
// 	)
// 	if err != nil {
// 		log.Printf("Error executing search: %s", err)
// 		return places, err
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		log.Printf("Error executing search: %s", res.String())
// 		return places, err
// 	}

// 	// Decode the search response
// 	var r map[string]interface{}
// 	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
// 		log.Printf("Error decoding search response: %s", err)
// 		return places, err
// 	}

// 	// Extract the hits from the search response
// 	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})

// 	// Prepare the response
// 	for _, hit := range hits {
// 		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
// 		var place placeDTO.Place
// 		if err := mapstructure.Decode(source, &place); err != nil {
// 			log.Printf("Error decoding place: %s", err)
// 			return places, err
// 		}
// 		places = append(places, place)
// 	}

// 	// Return the search results
// 	return places, nil
// }

// FUZZY SEARCH

// func (ss *SearchService) SearchPlaceByName(name string) ([]placeDTO.Place, error) {
// 	var places []placeDTO.Place

// 	// Prepare the search query using a fuzzy query
// 	var buf strings.Builder
// 	query := map[string]interface{}{
// 		"query": map[string]interface{}{
// 			"fuzzy": map[string]interface{}{
// 				"name": map[string]interface{}{
// 					"value":     name,
// 					"fuzziness": 2, // Adjust fuzziness as needed
// 				},
// 			},
// 		},
// 	}
// 	if err := json.NewEncoder(&buf).Encode(query); err != nil {
// 		return places, err
// 	}

// 	// Convert strings.Builder to strings.Reader
// 	queryReader := strings.NewReader(buf.String())

// 	// Execute the search query
// 	res, err := ss.elasticClient.Search(
// 		ss.elasticClient.Search.WithContext(context.Background()),
// 		ss.elasticClient.Search.WithIndex("places"),
// 		ss.elasticClient.Search.WithBody(queryReader),
// 		ss.elasticClient.Search.WithTrackTotalHits(true),
// 		ss.elasticClient.Search.WithPretty(),
// 	)
// 	if err != nil {
// 		log.Printf("Error executing search: %s", err)
// 		return places, err
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		log.Printf("Error executing search: %s", res.String())
// 		return places, err
// 	}

// 	// Decode the search response
// 	var r map[string]interface{}
// 	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
// 		log.Printf("Error decoding search response: %s", err)
// 		return places, err
// 	}

// 	// Extract the hits from the search response
// 	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})

// 	// Prepare the response
// 	for _, hit := range hits {
// 		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
// 		var place placeDTO.Place
// 		if err := mapstructure.Decode(source, &place); err != nil {
// 			log.Printf("Error decoding place: %s", err)
// 			return places, err
// 		}
// 		places = append(places, place)
// 	}

// 	// Return the search results
// 	return places, nil
// }

// MATCH QUERY WITH FUZZINESS

func (ss *SearchService) SearchPlaceByName(name string) ([]placeDTO.Place, error) {
	var places []placeDTO.Place

	// Prepare the search query using a match query with auto fuzziness
	var buf strings.Builder
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": map[string]interface{}{
					"query":     name,
					"fuzziness": "AUTO",
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return places, err
	}

	// Convert strings.Builder to strings.Reader
	queryReader := strings.NewReader(buf.String())

	// Execute the search query
	res, err := ss.elasticClient.Search(
		ss.elasticClient.Search.WithContext(context.Background()),
		ss.elasticClient.Search.WithIndex("places"),
		ss.elasticClient.Search.WithBody(queryReader),
		ss.elasticClient.Search.WithTrackTotalHits(true),
		ss.elasticClient.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("Error executing search: %s", err)
		return places, err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error executing search: %s", res.String())
		return places, err
	}

	// Decode the search response
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error decoding search response: %s", err)
		return places, err
	}

	// Extract the hits from the search response
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})

	// Prepare the response
	// for _, hit := range hits {
	// 	source := hit.(map[string]interface{})["_source"].(map[string]interface{})
	// 	var place placeDTO.Place
	// 	if err := mapstructure.Decode(source, &place); err != nil {
	// 		log.Printf("Error decoding place: %s", err)
	// 		return places, err
	// 	}

	// 	places = append(places, place)
	// }

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		var place placeDTO.Place
		if err := mapstructure.Decode(source, &place); err != nil {
			log.Printf("Error decoding place: %s", err)
			return places, err
		}

		// Check if the place's name includes the search term
		if strings.Contains(strings.ToLower(place.Name), strings.ToLower(name)) {
			places = append(places, place)
		}
	}

	// Return the search results
	return places, nil
}
