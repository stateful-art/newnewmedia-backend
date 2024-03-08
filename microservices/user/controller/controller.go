package controller

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"newnew.media/microservices/user/dto"
	"newnew.media/microservices/user/service"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var user dto.CreateUserRequest
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := c.userService.CreateUser(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
}

func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	user, err := c.userService.GetUserByID(objectID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (c *UserController) GetUsers(ctx *fiber.Ctx) error {
	users, err := c.userService.GetUsers()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return ctx.JSON(users)
}

func (c *UserController) GetUserBySpotifyID(ctx *fiber.Ctx) error {
	spotifyID := ctx.Params("spotify_id")

	user, err := c.userService.GetUserBySpotifyID(spotifyID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (c *UserController) GetUserByYouTubeID(ctx *fiber.Ctx) error {
	youtubeID := ctx.Params("youtube_id")

	user, err := c.userService.GetUserByYouTubeID(youtubeID)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to get user"})

	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (c *UserController) AddRole(ctx *fiber.Ctx) error {
	var addRoleRequest dto.UpdateRoleRequest

	if err := ctx.BodyParser(&addRoleRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	userObjID, err := primitive.ObjectIDFromHex(addRoleRequest.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid userID format"})
	}

	if err := c.userService.AddRole(userObjID, addRoleRequest.Role); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to add role to user"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Role added to user successfully"})
}

func (c *UserController) RemoveRole(ctx *fiber.Ctx) error {
	var removeRoleRequest dto.UpdateRoleRequest
	if err := ctx.BodyParser(&removeRoleRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	userObjID, err := primitive.ObjectIDFromHex(removeRoleRequest.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid userID format"})
	}

	if err := c.userService.RemoveRole(userObjID, removeRoleRequest.Role); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to remove role from user"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Role removed from user successfully"})
}
