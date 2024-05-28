package entities

import (
	"encoding/json"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MixProducts struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Images    UrlImage           `json:"image" bson:"image"`
	Price     int                `json:"price" bson:"price"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	V         int                `json:"__v" bson:"__v"`
}
type UrlImage struct {
	Url string `json:"url" bson:"url"`
}

func (mp *MixProducts) UnmarshalJSON(data []byte) error {
	type Alias MixProducts
	aux := &struct {
		Price string `json:"price"`
		*Alias
	}{
		Alias: (*Alias)(mp),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	price, err := strconv.Atoi(aux.Price)
	if err != nil {
		return err
	}
	mp.Price = price
	return nil
}
