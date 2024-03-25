package searchroutes

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	controller "newnew.media/microservices/search/controller"
	service "newnew.media/microservices/search/service"
)

func SearchRoutes(app fiber.Router, elasticClient *elasticsearch.Client, natsClient *nats.Conn) {
	// Pass the PlaylistService instance to the PlaylistController

	searchService := service.NewSearchService(elasticClient)
	indexService := service.NewIndexService(elasticClient, natsClient)
	// Subscribe to user-registered subject
	if err := indexService.SubscribeToPlaceCreatedSubject(); err != nil {
		log.Fatalf("Failed to subscribe to user-registered subject: %v", err)
	}

	searchIndexController := controller.NewSearchIndexController(searchService, indexService)

	// Admin Index routes
	app.Post("/admin/create-index/:indexName", searchIndexController.CreateIndex)
	app.Delete("/admin/delete-index/:indexName", searchIndexController.DeleteIndex)

	// Place routes
	app.Post("/index-place", searchIndexController.IndexPlace)
	app.Get("/place", searchIndexController.SearchPlaceByName)
	// app.Get("/place/prefix", controller.PrefixSearchPlaceByName)

}
