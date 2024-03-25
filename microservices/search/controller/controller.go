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

func (sic *SearchIndexController) CreateIndex(c *fiber.Ctx) error {
	indexName := c.Params("indexName")
	mapping := c.Body()

	err := sic.indexService.CreateIndex(indexName, mapping)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Index created successfully"})
}

func (sic *SearchIndexController) DeleteIndex(c *fiber.Ctx) error {
	indexName := c.Params("indexName")

	err := sic.indexService.DeleteIndex(indexName)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Index deleted successfully"})
}

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
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Place indexed successfully"})
}

func (sic *SearchIndexController) SearchPlaceByName(c *fiber.Ctx) error {
	// Extract the name query parameter from the request
	name := c.Query("name")
	if name == "" {
		return c.Status(400).SendString("Name query parameter is required")
	}
	var places []placeDTO.Place
	places, err := sic.searchService.SearchPlaceByName(name)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No place found."})
	}
	// Return the search results
	return c.Status(fiber.StatusAccepted).JSON(places)
}
