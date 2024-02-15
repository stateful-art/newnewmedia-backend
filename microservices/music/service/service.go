package service

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	dao "newnewmedia.com/microservices/music/dao"
	dto "newnewmedia.com/microservices/music/dto"
	repository "newnewmedia.com/microservices/music/repository"
	utils "newnewmedia.com/microservices/music/utils"
)

func CreateMusic(c *fiber.Ctx, music dto.Music) error {

	saveMusic, err := utils.SaveMusicFile(c)
	if err != nil {
		return err
	}
	placeObjID, err := primitive.ObjectIDFromHex(music.Place)
	if err != nil {
		return err
	}
	daoMusic := dao.Music{
		Name:   music.Name,
		Artist: music.Artist,
		Path:   saveMusic,
		Place:  placeObjID,
	}
	repoErr := repository.CreateMusic(c, daoMusic)
	if repoErr != nil {
		return repoErr
	}
	return nil
}

func GetMusicByPlace(c *fiber.Ctx, id string) ([]dao.Music, error) {
	placeObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	music, err := repository.GetMusicByPlace(c, &placeObjID)
	if err != nil {
		return nil, err
	}
	return music, nil
}
