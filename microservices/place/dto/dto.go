package dto

// PlaceDto is a struct that represents the data transfer object for the place microservice
// type Location struct {
// 	Address   string  `json:"address"`
// 	City      string  `json:"city"`
// 	Country   string  `json:"country"`
// 	Latitude  float64 `json:"latitude"`
// 	Longitude float64 `json:"longitude"`
// }

type GeometryType string

const (
	Point              GeometryType = "Point"
	LineString         GeometryType = "LineString"
	Polygon            GeometryType = "Polygon"
	MultiPoint         GeometryType = "MultiPoint"
	MultiLineString    GeometryType = "MultiLineString"
	MultiPolygon       GeometryType = "MultiPolygon"
	GeometryCollection GeometryType = "GeometryCollection"
)

// Location struct with a GeometryType field
type Location struct {
	Type        GeometryType `json:"type"`
	Coordinates []float64    `json:"coordinates"`
}
type Link struct {
	Platform string `json:"platform"`
	Url      string `json:"url"`
}

type Place struct {
	ID        string   `json:"id,omitempty" bson:"_id,omitempty"`
	Owner     string   `json:"owner,omitempty" bson:"owner,omitempty"`
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	SpotifyID string   `json:"spotify_id" bson:"spotify_id"`
	Name      string   `json:"name"`
	Location  Location `json:"location"`
	City      string   `json:"city"`
	Country   string   `json:"country"`

	Description string `json:"description"`
	Image       string `json:"image"`
	Links       []Link `json:"links"`
}

type RecentPlayItem struct {
	SongID     string `bson:"song_id"`
	SongName   string `bson:"song_name"`
	ArtistName string `bson:"artist_name"`
	PlaceName  string `bson:"place_name"`
	Duration   int    `bson:"duration"`
	PlaceID    string `bson:"place_id"`
}

type RecentPlays struct {
	ID    string           `bson:"_id,omitempty"`
	Place string           `bson:"place_id"` // For place's recent plays
	Items []RecentPlayItem `bson:"items"`
}

//  GeometryType Examples:

// Point
// A Point represents a single point in space. It's the simplest form of geometry and is used to represent a single location.

// var pointLocation Location
// pointLocation.Type = Point
// pointLocation.Coordinates = []float64{102.0, 0.5} // Longitude, Latitude

// LineString
// A LineString represents a series of points connected by line segments. It's used to represent a path or a line.

// var lineStringLocation Location
// lineStringLocation.Type = LineString
// lineStringLocation.Coordinates = []float64{
// 	102.0, 0.0, // First point
// 	103.0, 1.0, // Second point
// 	104.0, 0.0, // Third point
// 	105.0, 1.0, // Fourth point
// }

// Polygon
// A Polygon represents a closed shape defined by a series of points. It's used to represent areas or regions.

// var polygonLocation Location
// polygonLocation.Type = Polygon
// polygonLocation.Coordinates = []float64{
// 	30.0, 10.0, // First point
// 	40.0, 40.0, // Second point
// 	20.0, 40.0, // Third point
// 	10.0, 20.0, // Fourth point
// 	30.0, 10.0, // Closing point (same as first point)
// }

// var kadikoyLocation Location
// kadikoyLocation.Type = Polygon
// kadikoyLocation.Coordinates = []float64{
// 	28.9833, 41.0000, // First point (approximate)
// 	28.9833, 41.0000, // Second point (approximate)
// 	28.9833, 41.0000, // Third point (approximate)
// 	28.9833, 41.0000, // Fourth point (approximate)
// 	28.9833, 41.0000, // Closing point (same as first point)
// }

// fmt.Printf("Kadikoy District: %+v\n", kadikoyLocation)
// }

// MultiPoint
// A MultiPoint represents a collection of points.

// var multiPointLocation Location
// multiPointLocation.Type = MultiPoint
// multiPointLocation.Coordinates = []float64{
// 	10.0, 10.0, // First point
// 	20.0, 20.0, // Second point
// 	30.0, 30.0, // Third point
// }

// MultiLineString
// A MultiLineString represents a collection of line strings.

// var multiLineStringLocation Location
// multiLineStringLocation.Type = MultiLineString
// multiLineStringLocation.Coordinates = []float64{
// 	10.0, 10.0, 20.0, 20.0, // First line string
// 	30.0, 30.0, 40.0, 40.0, // Second line string
// }

// MultiPolygon
// A MultiPolygon represents a collection of polygons.

// var multiPolygonLocation Location
// multiPolygonLocation.Type = MultiPolygon
// multiPolygonLocation.Coordinates = []float64{
// 	30.0, 20.0, 10.0, 20.0, 40.0, 40.0, 20.0, 40.0, 30.0, 20.0, // First polygon
// 	30.0, 10.0, 10.0, 10.0, 20.0, 20.0, 10.0, 20.0, 30.0, 10.0, // Second polygon
// }

// var kadikoyLocation Location
// 	kadikoyLocation.Type = MultiPolygon
// 	kadikoyLocation.Coordinates = [][][]float64{
// 		{
// 			{28.9833, 41.0000}, // First point of the first polygon
// 			{28.9833, 41.0000}, // Second point of the first polygon
// 			{28.9833, 41.0000}, // Third point of the first polygon
// 			{28.9833, 41.0000}, // Fourth point of the first polygon
// 			{28.9833, 41.0000}, // Closing point (same as first point)
// 		},
// 		{
// 			{28.9833, 41.0000}, // First point of the second polygon
// 			{28.9833, 41.0000}, // Second point of the second polygon
// 			{28.9833, 41.0000}, // Third point of the second polygon
// 			{28.9833, 41.0000}, // Fourth point of the second polygon
// 			{28.9833, 41.0000}, // Closing point (same as first point)
// 		},
// 		// Add more polygons as needed
// 	}

// 	fmt.Printf("Kadikoy District: %+v\n", kadikoyLocation)

// GeometryCollection
// A GeometryCollection represents a collection of geometries of any type.

// var geometryCollectionLocation Location
// geometryCollectionLocation.Type = GeometryCollection
// geometryCollectionLocation.Coordinates = []float64{
// 	// Point
// 	10.0, 10.0,
// 	// LineString
// 	20.0, 20.0, 30.0, 30.0,
// 	// Polygon
// 	40.0, 40.0, 50.0, 50.0, 60.0, 60.0, 40.0, 40.0,
// }
