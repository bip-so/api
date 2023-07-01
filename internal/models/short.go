package models

type Short struct {
	BaseModel
	OriginalURL string
	Method      string
	ShortCode   string
	Count       int
}

func (m Short) NewShort(originalURL string, method string, shortCode string, count int) *Short {
	return &Short{OriginalURL: originalURL, Method: method, ShortCode: shortCode, Count: count}
}
