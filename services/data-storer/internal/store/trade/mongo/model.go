package mongo

import (
	"fmt"
	"time"

	"github.com/cshep4/kripto/services/data-storer/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type trade struct {
	Id       primitive.ObjectID `bson:"_id"`
	GBP      float64            `bson:"gbp"`
	BTC      float64            `bson:"btc"`
	Rate     float64            `bson:"rate"`
	Type     string             `bson:"type"`
	DateTime time.Time          `bson:"dateTime"`
}

func fromTrade(t model.Trade) (*trade, error) {
	id := primitive.NewObjectID()

	if t.Id != "" {
		var err error
		id, err = primitive.ObjectIDFromHex(t.Id)
		if err != nil {
			return nil, fmt.Errorf("object_id_from_hex: %s", t.Id)
		}
	}

	return &trade{
		Id:       id,
		GBP:      t.GBP,
		BTC:      t.BTC,
		Rate:     t.Rate,
		Type:     string(t.Type),
		DateTime: t.DateTime,
	}, nil
}

func toTrade(t trade) model.Trade {
	return model.Trade{
		Id:       t.Id.Hex(),
		GBP:      t.GBP,
		BTC:      t.BTC,
		Rate:     t.Rate,
		Type:     model.TradeType(t.Type),
		DateTime: t.DateTime,
	}
}
