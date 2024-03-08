package repository

import (
	"context"
	"errors"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	collections "newnew.media/db/collections"
	dao "newnew.media/microservices/user/dao"
)

type UserRepository struct {
	// Any fields or dependencies needed by the repository can be added here
}

// NewUserRepository creates a new instance of the UserRepository.
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// CreateUser inserts a new user into the database.
func (ur *UserRepository) CreateUser(user dao.User) error {

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := collections.UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID retrieves a user by their ID.
func (ur *UserRepository) GetUserByID(id primitive.ObjectID) (dao.User, error) {
	var user dao.User

	filter := bson.M{"_id": id}

	err := collections.UsersCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return dao.User{}, err
	}

	return user, nil
}

// GetUsers retrieves all users from the database.
func (ur *UserRepository) GetUsers() ([]dao.User, error) {
	var users []dao.User

	cursor, err := collections.UsersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user dao.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (dao.User, error) {
	var user dao.User

	filter := bson.M{"email": email}

	err := collections.UsersCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return dao.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserBySpotifyID(spotifyID string) (dao.User, error) {
	var user dao.User

	filter := bson.M{"spotifyID": spotifyID}

	err := collections.UsersCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return dao.User{}, err
	}

	return user, nil
}

// GetUserByYouTubeID retrieves a user by their YouTube ID.
func (ur *UserRepository) GetUserByYouTubeID(youtubeID string) (dao.User, error) {
	var user dao.User

	filter := bson.M{"youtubeID": youtubeID}

	err := collections.UsersCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return dao.User{}, err
	}

	return user, nil
}

// GetUserByFavoriteGenres retrieves users by their favorite genres.
// GetUserByFavoriteGenres retrieves users by their favorite genres.
func (ur *UserRepository) GetUserByFavoriteGenres(genres []primitive.ObjectID) ([]dao.User, error) {
	var users []dao.User

	filter := bson.M{"favoriteGenres": bson.M{"$in": genres}}

	cursor, err := collections.UsersCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user dao.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserByFavoritePlaces retrieves users by their favorite places.
func (ur *UserRepository) GetUserByFavoritePlaces(places []primitive.ObjectID) ([]dao.User, error) {
	var users []dao.User

	filter := bson.M{"favoritePlaces": bson.M{"$in": places}}

	cursor, err := collections.UsersCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user dao.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateUser updates an existing user in the database.
// UpdateUser updates an existing user in the database.
func (ur *UserRepository) UpdateUser(id primitive.ObjectID, user dao.User) error {
	filter := bson.M{"_id": id}

	update := ur.generateUpdateQuery(user)

	_, err := collections.UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

//

// DeleteUser deletes a user from the database.
func (ur *UserRepository) DeleteUser(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := collections.UsersCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no document found to delete")
	}

	return nil
}

// generateUpdateQuery dynamically generates the update query based on the provided User.
func (ur *UserRepository) generateUpdateQuery(user dao.User) bson.M {
	update := bson.M{"$set": bson.M{}}

	// Iterate over the fields of the User struct and add them to the update query if they are not zero or empty.
	userValue := reflect.ValueOf(user)
	for i := 0; i < userValue.NumField(); i++ {
		field := userValue.Field(i)
		fieldName := userValue.Type().Field(i).Name
		if !field.IsZero() && field.Interface() != "" {
			update["$set"].(bson.M)[fieldName] = field.Interface()
		}
	}

	// Set updatedAt field
	update["$set"].(bson.M)["updatedAt"] = time.Now()

	return update
}

// AddUserRole adds a role to the user.
func (ur *UserRepository) AddUserRole(userRole dao.UserRole) error {
	_, err := collections.UserRolesCollection.InsertOne(context.Background(), userRole)
	if err != nil {
		return err
	}
	return nil
}

// RemoveUserRole removes a role from the user.
func (ur *UserRepository) RemoveUserRole(userID primitive.ObjectID, role dao.Role) error {
	filter := bson.M{"user_id": userID, "role": role}
	_, err := collections.UserRolesCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) RollbackUserCreation(userID primitive.ObjectID) error {
	filter := bson.M{"_id": userID}

	// Delete the user document from the collection
	_, err := collections.UsersCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
