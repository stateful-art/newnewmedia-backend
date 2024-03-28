package offerroutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	collections "newnew.media/db/collections"

	controller "newnew.media/microservices/offer/controller"
	repo "newnew.media/microservices/offer/repository"
	service "newnew.media/microservices/offer/service"
)

func OfferRoutes(app fiber.Router, natsClient *nats.Conn) {

	repo := repo.NewMongoOfferRepository(collections.OffersCollection.Database())
	// Create an instance of PlaylistService with the PlaylistRepository
	service := service.NewOfferService(repo)
	controller := controller.NewOfferController(service)

	app.Get("/:id", controller.GetOfferByID)
	app.Post("/", controller.CreateOffer)
	app.Patch("/:id", controller.UpdateOfferStatus)
	app.Delete("/:id", controller.DeleteOffer)

	app.Get("/counter/:id", controller.GetCounterOfferByID)
	app.Post("/counter", controller.CheckOfferValidity(), controller.CreateCounterOffer)
	app.Patch("/counter/:id", controller.CheckOfferValidity(), controller.UpdateCounterOfferStatus)

}
