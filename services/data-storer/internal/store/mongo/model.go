package mongo

import (
	"fmt"
	"time"

	"github.com/cshep4/kripto/services/data-storer/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	Id              primitive.ObjectID `bson:"_id"`
	FirstName       string             `bson:"firstName"`
	LastName        string             `bson:"lastName"`
	Email           string             `bson:"email"`
	Type            model.UserType     `bson:"type"`
	Verified        bool               `bson:"verified"`
	Password        string             `bson:"password"`
	Joined          time.Time          `bson:"joined"`
	RentReminders   bool               `bson:"rentReminders"`
	VerifySignature string             `bson:"verifySignature"`
	ResetSignature  string             `bson:"resetSignature"`
}

func fromUser(u model.User) (user, error) {
	id, err := primitive.ObjectIDFromHex(u.Id)
	if err != nil {
		return user{}, fmt.Errorf("cannot_create_object_id_from_hex: %s", u.Id)
	}

	return user{
		Id:              id,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Email:           u.Email,
		Type:            u.Type,
		Verified:        u.Verified,
		Password:        u.Password,
		Joined:          u.Joined,
		RentReminders:   u.RentReminders,
		VerifySignature: u.VerifySignature,
		ResetSignature:  u.ResetSignature,
	}, nil
}

func toUser(u user) model.User {
	return model.User{
		Id:              u.Id.Hex(),
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Email:           u.Email,
		Type:            u.Type,
		Verified:        u.Verified,
		Password:        u.Password,
		Joined:          u.Joined,
		RentReminders:   u.RentReminders,
		VerifySignature: u.VerifySignature,
		ResetSignature:  u.ResetSignature,
	}
}
