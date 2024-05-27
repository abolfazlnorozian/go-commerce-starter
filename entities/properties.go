package entities

import (
	"encoding/json"
	"strconv"
	"time"
)

type Properties struct {
	ID        int       `json:"_id" bson:"_id"`
	Parent    *int      `json:"parent" bson:"parent"`
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	V         int       `json:"__v" bson:"__v"`
}

func (p *Properties) UnmarshalJSON(data []byte) error {
	var temp struct {
		ID        int       `json:"_id"`
		Parent    string    `json:"parent"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		V         int       `json:"__v"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	p.ID = temp.ID
	p.Name = temp.Name
	p.CreatedAt = temp.CreatedAt
	p.UpdatedAt = temp.UpdatedAt
	p.V = temp.V

	if temp.Parent != "" {
		parentID, err := strconv.Atoi(temp.Parent)
		if err != nil {
			return err
		}
		p.Parent = &parentID
	}

	return nil
}
