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

	searchService := service.NewSearchService(elasticClient, natsClient)
	indexService := service.NewIndexService(elasticClient, natsClient)
	// Subscribe to user-registered subject
	if err := indexService.SubscribeToPlaceCreatedSubject(); err != nil {
		log.Fatalf("Failed to subscribe to user-registered subject: %v", err)
	}

	sic := controller.NewSearchIndexController(searchService, indexService)

	// Admin Index routes
	// app.Post("/admin/create-index/:indexName", controller.CreateIndex)
	// app.Delete("/admin/delete-index/:indexName", controller.DeleteIndex)

	// Place routes
	app.Post("/index-place", sic.IndexPlace)
	// app.Get("/place", controller.SearchPlaceByName)
	// app.Get("/place/prefix", controller.PrefixSearchPlaceByName)

}
