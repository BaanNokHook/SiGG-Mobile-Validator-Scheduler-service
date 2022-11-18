package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Device struct {
	ID                        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId                    string             `bson:"userId" json:"userId"`
	DeviceId                  string             `bson:"deviceId" json:"deviceId"`
	PublicKey                 string             `bson:"publicKey" json:"publicKey"`
	LatestSync                int64              `bson:"latestSync" json:"latestSync"`
	Created                   int64              `bson:"created" json:"created"`
	LatestValidateTransaction int64              `bson:"latestValidateTransaction" json:"latestValidateTransaction"`
}
