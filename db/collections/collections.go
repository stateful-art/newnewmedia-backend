package collections

import (
	db "newnewmedia.com/db"
)

var (
	UsersCollection = db.Client.Database("newnewmedia").Collection("users")
)
