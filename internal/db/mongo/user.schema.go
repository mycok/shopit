package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollectionSchema = dbCollection{
	"users": &options.CreateCollectionOptions{
		ValidationAction: &validationAction,
		ValidationLevel:  &validationLevel,
		Validator: bson.M{
			"$jsonSchema": bson.M{
				"bsonType": "object",
				"required": []string{"username", "email", "password"},
				"additionalProperties": false,
				"properties": bson.M{
					"username": bson.M{
						"bsonType":    "string",
						"description": "user name",
					},
					"email": bson.M{
						"bsonType":    "string",
						"description": "user email address",
					},
					"password": bson.M{
						"bsonType": "binData",
						"description": "user hashed password",
					},
					"isActive": bson.M{
						"bsonType":    "bool",
						"description": "user active status",
					},
					"version": bson.M{
						"bsonType":    "string",
						"description": "user record version since last update",
					},
					"created_at": bson.M{
						"bsonType":    "date",
						"description": "user record created date and time",
					},
				},
			},
		},
	},
}
