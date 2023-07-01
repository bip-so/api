package utils

import gonanoid "github.com/matoous/go-nanoid/v2"

func NewNanoid() string {
	id, err := gonanoid.New(12)
	if err != nil {
		return "" // sending empty
	}
	return id
}

func NewShortNanoid() string {
	id, err := gonanoid.New(6)
	if err != nil {
		return "" // sending empty
	}
	return id
}

func HandleExtender(length int) string {
	id, err := gonanoid.New(length)
	if err != nil {
		return "" // sending empty
	}
	return id
}
