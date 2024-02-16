package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	collections "newnewmedia.com/db/Collections"
	dao "newnewmedia.com/microservices/music/dao"
)

func CreateMusic(music dao.Music) error {
	_, err := collections.MusicCollection.InsertOne(context.Background(), music)
	if err != nil {
		return err
	}
	return nil
}

func GetMusicByPlace(id *primitive.ObjectID) ([]dao.Music, error) {
	var music []dao.Music
	cursor, err := collections.MusicCollection.Find(context.Background(), bson.M{"place": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var m dao.Music
		cursor.Decode(&m)
		music = append(music, m)
	}
	return music, nil
}

func GetMusicById(id *primitive.ObjectID) (dao.Music, error) {
	var music dao.Music

	// Define the filter to find music by ID
	filter := bson.M{"_id": id}

	// Execute the query to find music documents
	cursor, err := collections.MusicCollection.Find(context.Background(), filter)
	if err != nil {
		return music, err
	}
	defer cursor.Close(context.Background())

	// Decode the first document from the cursor into the music object
	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&music); err != nil {
			return music, err
		}
	}

	return music, nil
}
