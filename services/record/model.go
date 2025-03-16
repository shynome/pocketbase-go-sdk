package record

type Base struct {
	CollectionID   string   `json:"collectionId"`
	CollectionName string   `json:"collectionName"`
	ID             string   `json:"id"`
	Created        DateTime `json:"created"`
	Updated        DateTime `json:"updated"`
}

type UserBase struct {
	Base
	Email           string `json:"email"`
	EmailVisibility bool   `json:"emailVisibility"`
	Verified        bool   `json:"verified"`

	// Deprecated: after v0.23.0 username is remove, email is the default unique feild
	Username string `json:"username"`
}
