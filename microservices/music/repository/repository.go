package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	collections "newnewmedia.com/db/collections"
	dao "newnewmedia.com/microservices/music/dao"
)

func CreateMusic(music dao.Song) error {

	_, err := collections.MusicCollection.InsertOne(context.Background(), music)
	if err != nil {
		return err
	}
	return nil
}

func GetMusicByPlace(id *primitive.ObjectID) ([]dao.Song, error) {
	var music []dao.Song
	cursor, err := collections.MusicCollection.Find(context.Background(), bson.M{"place": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var m dao.Song
		cursor.Decode(&m)
		music = append(music, m)
	}
	return music, nil
}

func GetMusicById(id primitive.ObjectID) (dao.Song, error) {
	var music dao.Song
	log.Println(id)
	// Define the filter to find music by ID
	filter := bson.M{"_id": id}

	// Execute the query to find music documents
	err := collections.MusicCollection.FindOne(context.Background(), filter).Decode(&music)
	if err != nil {
		return music, err
	}

	log.Println(music)

	return music, nil
}
