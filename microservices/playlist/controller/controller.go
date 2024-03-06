package controller

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	dao "newnewmedia.com/microservices/playlist/dao"
	dto "newnewmedia.com/microservices/playlist/dto"
	"newnewmedia.com/microservices/playlist/service"
)

type PlaylistController struct {
	playlistService *service.PlaylistService
}

func NewPlaylistController(playlistService *service.PlaylistService) *PlaylistController {
	return &PlaylistController{playlistService: playlistService}
}

func (c *PlaylistController) CreatePlaylist(ctx *fiber.Ctx) error {
	var playlist dao.Playlist
	if err := ctx.BodyParser(&playlist); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	// Call the service function with the additional parameters
	if err := c.playlistService.CreatePlaylist(playlist); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create playlist"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Playlist created successfully"})
}

func (c *PlaylistController) GetPlaylistByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid playlist ID"})
	}

	playlist, err := c.playlistService.GetPlaylistByID(objectID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Playlist not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(playlist)
}

// GetPlaylists retrieves all playlists
func (c *PlaylistController) GetPlaylists(ctx *fiber.Ctx) error {
	playlists, err := c.playlistService.GetPlaylists()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return ctx.JSON(playlists)
}

func (c *PlaylistController) UpdatePlaylist(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid playlist ID"})
	}

	var playlist dao.Playlist
	if err := ctx.BodyParser(&playlist); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := c.playlistService.UpdatePlaylist(objectID, playlist); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update playlist"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Playlist updated successfully"})
}

func (c *PlaylistController) DeletePlaylist(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid playlist ID"})
	}

	if err := c.playlistService.DeletePlaylist(objectID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to delete playlist"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Playlist deleted successfully"})
}

func (c *PlaylistController) AddSongsToPlaylist(ctx *fiber.Ctx) error {
	var PlaylistSongsAddRequest dto.PlaylistSongsUpdateRequest

	if err := ctx.BodyParser(&PlaylistSongsAddRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	// Call the service function to add songs to the playlist
	if err := c.playlistService.AddSongsToPlaylist(PlaylistSongsAddRequest.PlaylistID, PlaylistSongsAddRequest.SongIDs); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to add songs to playlist"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Songs added to playlist successfully"})
}

func (c *PlaylistController) RemoveSongsFromPlaylist(ctx *fiber.Ctx) error {
	var PlaylistSongsRemoveRequest dto.PlaylistSongsUpdateRequest

	if err := ctx.BodyParser(&PlaylistSongsRemoveRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	// Call the service function to remove songs from the playlist
	if err := c.playlistService.RemoveSongsFromPlaylist(PlaylistSongsRemoveRequest.PlaylistID, PlaylistSongsRemoveRequest.SongIDs); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to remove songs from playlist"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Songs removed from playlist successfully"})
}
