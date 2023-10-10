package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GenericTimeStampedDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at"`
}

type GenericTimeStampedModel struct {
	ID        primitive.ObjectID `json:"id"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}
