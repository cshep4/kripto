package mongo

import (
	"fmt"
	"time"

	"github.com/cshep4/kripto/services/data-storer/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type rate struct {
	Id       primitive.ObjectID `bson:"_id"`
	Rate     float64            `bson:"rate"`
	DateTime time.Time          `bson:"dateTime"`
}

func fromRate(r model.Rate) (*rate, error) {
	id := primitive.NewObjectID()

	if r.Id != "" {
		var err error
		id, err = primitive.ObjectIDFromHex(r.Id)
		if err != nil {
			return nil, fmt.Errorf("object_id_from_hex: %s", r.Id)
		}
	}

	return &rate{
		Id:       id,
		Rate:     r.Rate,
		DateTime: r.DateTime,
	}, nil
}

func toRate(r rate) model.Rate {
	return model.Rate{
		Id:       r.Id.Hex(),
		Rate:     r.Rate,
		DateTime: r.DateTime,
	}
}
