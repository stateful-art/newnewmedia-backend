package collections

import (
	db "newnew.media/db"
)

var (
	UsersCollection     = db.Client.Database("newnewmedia").Collection("users")
	PlacesCollection    = db.Client.Database("newnewmedia").Collection("places")
	MusicCollection     = db.Client.Database("newnewmedia").Collection("music")
	PlaylistsCollection = db.Client.Database("newnewmedia").Collection("playlists")
	ArtistsCollection   = db.Client.Database("newnewmedia").Collection("artists")
	RevenuesCollection  = db.Client.Database("newnewmedia").Collection("revenues")
	UserRolesCollection = db.Client.Database("newnewmedia").Collection("user_roles")
)
