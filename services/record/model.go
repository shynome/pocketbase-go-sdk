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
	Username        string `json:"username"`
	Email           string `json:"email"`
	EmailVisibility bool   `json:"emailVisibility"`
	Verified        bool   `json:"verified"`
}
