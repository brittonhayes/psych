package api

type Therapist struct {
	ID          int    `bun:"id,pk,autoincrement" json:"id"`
	Title       string `json:"title"`
	Credentials string `json:"credentials"`
	Verified    string `json:"verified"`
	Statement   string `json:"statement"`
	Phone       string `json:"phone"`
	Location    string `json:"location"`
	Link        string `json:"link"`
}

type GetTherapistParams struct {
	Title       *string `json:"title"`
	Credentials *string `json:"credentials"`
	Verified    *string `json:"verified"`
	Statement   *string `json:"statement"`
	Phone       *string `json:"phone"`
	Location    *string `json:"location"`
	Link        *string `json:"link"`
	Limit       *int    `json:"limit"`
	Offset      *int    `json:"offset"`
}
