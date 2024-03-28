package dto

type Status string

const (
	Pending   Status = "pending"
	Accepted  Status = "accepted"
	Rejected  Status = "rejected"
	Countered Status = "countered"
)

type Preference string

const (
	Public     Preference = "public"
	Private    Preference = "private"
	Collective Preference = "collective"
	Individual Preference = "individual"
)

type Offer struct {
	ID            string         `json:"id"`
	Songs         []string       `json:"songs"`
	ArtistID      string         `json:"artist_id"`
	PlaceID       string         `json:"place_id"`
	ValidUntil    string         `json:"valid_until"`
	OfferedAt     string         `json:"offered_at"`
	Preferences   []string       `json:"preferences"`
	CounterOffers []CounterOffer `json:"counter_offers"`
}

type Counter struct {
	ID          string   `json:"id"`
	OfferID     string   `json:"offer_id"`
	OfferedAt   string   `json:"offered_at"`
	ValidUntil  string   `json:"valid_until"`
	Status      string   `json:"status"`
	Preferences []string `json:"preferences"`
	ParentOffer string   `json:"parent_offer"`
}

type CounterOffer struct {
	CounterID string `json:"counter_id"`
	Status    string `json:"status"`
}

type CreateOffer struct {
	Songs       []string `json:"songs"`
	ArtistID    string   `json:"artist_id"`
	PlaceID     string   `json:"place_id"`
	ValidUntil  int      `json:"valid_until"` // number of days from now.
	Preferences []string `json:"preferences"`
}

type CreateOfferResponse struct {
	ID         string   `json:"id"`
	Songs      []string `json:"songs"`
	PlaceID    string   `json:"place_id"`
	OfferedAt  string   `json:"offered_at"`
	ValidUntil string   `json:"valid_until"`
	Status     string   `json:"status"`
}

type UpdateOfferStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type UpdateOfferPreferences struct {
	ID          string       `json:"id"`
	Preferences []Preference `json:"status"`
}

type CreateCounterOffer struct {
	OfferID     string   `json:"offer_id"`
	ValidUntil  int      `json:"valid_until"` // number of days from now.
	Preferences []string `json:"preferences"`
}

type NegotiationHistory struct {
	OriginalOfferID string         `json:"original_offer_id"`
	CounterOffers   []CounterOffer `json:"counter_offers"`
}
