package valueobject

import "errors"

type Rating int

func NewRating(value int) (Rating, error) {
	if value < 1 || value > 5 {
		return 0, errors.New("rating must be between 1 and 5")
	}
	return Rating(value), nil
}

func (r Rating) Int() int {
	return int(r)
}
