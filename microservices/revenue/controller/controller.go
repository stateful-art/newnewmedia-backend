package controller

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"newnewmedia.com/microservices/revenue/dao"
	"newnewmedia.com/microservices/revenue/dto"
	"newnewmedia.com/microservices/revenue/service"
)

type RevenueController struct {
	revenueService *service.RevenueService
}

func NewRevenueController(revenueService *service.RevenueService) *RevenueController {
	return &RevenueController{revenueService: revenueService}
}

func (c *RevenueController) CreateRevenue(ctx *fiber.Ctx) error {
	var revenue dao.Revenue
	if err := ctx.BodyParser(&revenue); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := c.revenueService.CreateRevenue(revenue); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create revenue entry"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Revenue entry created successfully"})
}

func (c *RevenueController) GetRevenueByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid revenue ID"})
	}

	revenue, err := c.revenueService.GetRevenueByID(objectID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Revenue entry not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(revenue)
}

func (c *RevenueController) GetRevenueByArtistID(ctx *fiber.Ctx) error {
	id := ctx.Params("artist_id")
	artistID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid artist ID"})
	}

	revenues, err := c.revenueService.GetRevenueByArtistID(artistID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve revenue entries"})
	}

	return ctx.Status(fiber.StatusOK).JSON(revenues)
}

func (c *RevenueController) GetRevenueByPlaceID(ctx *fiber.Ctx) error {
	id := ctx.Params("place_id")
	placeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid place ID"})
	}

	revenues, err := c.revenueService.GetRevenueByPlaceID(placeID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve revenue entries"})
	}

	return ctx.Status(fiber.StatusOK).JSON(revenues)
}

func (c *RevenueController) GetRevenueByPlaylistID(ctx *fiber.Ctx) error {
	id := ctx.Params("playlist_id")
	playlistID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid playlist ID"})
	}

	revenues, err := c.revenueService.GetRevenueByPlaylistID(playlistID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve revenue entries"})
	}

	return ctx.Status(fiber.StatusOK).JSON(revenues)
}

func (c *RevenueController) CalculateCollectiveRevenueSplit(ctx *fiber.Ctx) error {
	var RevenueSharingRequest dto.RevenueCalculateRequest

	if err := ctx.BodyParser(&RevenueSharingRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	// Call the service function to calculate collective revenue split
	artistShares, err := c.revenueService.CalculateCollectiveRevenueSplit(RevenueSharingRequest.Playlist, RevenueSharingRequest.TotalRevenue)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.JSON(artistShares)
}

func (c *RevenueController) CalculateIndividualRevenueSplit(ctx *fiber.Ctx) error {
	var RevenueSharingRequest dto.RevenueCalculateRequest

	if err := ctx.BodyParser(&RevenueSharingRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	// Call the service function to calculate individual revenue split
	artistShares, err := c.revenueService.CalculateIndividualRevenueSplit(RevenueSharingRequest.Playlist, RevenueSharingRequest.TotalRevenue)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.JSON(artistShares)
}
