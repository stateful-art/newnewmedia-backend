package controller

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	dto "newnew.media/microservices/offer/dto"
	service "newnew.media/microservices/offer/service"
)

type OfferController struct {
	offerService service.OfferService
}

var layoutString string = os.Getenv("TIME_LAYOUT_STRING")

// NewOfferController creates a new instance of OfferController with the provided OfferService.
func NewOfferController(offerService service.OfferService) *OfferController {
	return &OfferController{offerService: offerService}
}

// CreateOffer handles the creation of a new offer.
func (oc *OfferController) CreateOffer(ctx *fiber.Ctx) error {
	var dto dto.CreateOffer
	if err := ctx.BodyParser(&dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	offer, err := oc.offerService.CreateOffer(context.Background(), &dto)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create offer"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(offer)
}

// GetOfferByID handles getting an offer by its ID.
func (oc *OfferController) GetOfferByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	offer, err := oc.offerService.GetOfferByID(context.Background(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Offer not found"})
	}

	return ctx.JSON(offer)
}

// GetOffersByArtist handles getting offers by artist ID.
func (oc *OfferController) GetOffersByArtist(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	offers, err := oc.offerService.GetOffersByArtist(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(offers)
}

// GetOffersByPlace handles getting offers by place ID.
func (oc *OfferController) GetOffersByPlace(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	offers, err := oc.offerService.GetOffersByPlace(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(offers)
}

// UpdateOfferStatus handles updating an offer's status.
func (oc *OfferController) UpdateOfferStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var dto dto.UpdateOfferStatus
	if err := ctx.BodyParser(&dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	err := oc.offerService.UpdateOfferStatus(context.Background(), id, "")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update offer status"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// DeleteOffer handles deleting an offer by its ID.
func (oc *OfferController) DeleteOffer(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	err := oc.offerService.DeleteOffer(context.Background(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete offer"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (oc *OfferController) CheckOfferValidity() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Try to extract the offer ID from the URL parameters
		offerID := ctx.Params("id")

		// If the offer ID is not in the URL parameters, try to extract it from the request body
		// if offerID == "" {
		// 	var dto dto.CreateCounterOffer
		// 	if err := ctx.BodyParser(&dto); err != nil {
		// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		// 	}
		// 	offerID = dto.OfferID
		// }

		// Fetch the offer from the database
		offer, err := oc.offerService.GetOfferByID(context.Background(), offerID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch offer"})
		}

		validUntil, err := time.Parse(layoutString, offer.ValidUntil)
		if err != nil {
			fmt.Println("Error parsing time:", err)
		}

		// Check if the offer is still valid
		if time.Now().After(validUntil) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Offer is no longer valid"})
		}

		// If the offer is valid, proceed with the request
		return ctx.Next()
	}
}
