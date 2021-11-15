package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var usersIndexModels = collectionIndexModel{
	"users": []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{bson.E{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	},
}

var tokensIndexModels = collectionIndexModel{
	"tokens": []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "expiry", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
		{
			Keys:    bson.D{bson.E{Key: "hash", Value: 1}},
			Options: options.Index().SetName("hash_index"),
		},
		{
			Keys:    bson.D{bson.E{Key: "scope", Value: 1}},
			Options: options.Index().SetName("scope_index"),
		},
	},
}

var usersCollectionSchema = dbCollection{
	"users": &options.CreateCollectionOptions{
		ValidationAction: &validationAction,
		ValidationLevel:  &validationLevel,
		Validator: bson.M{
			"$jsonSchema": bson.M{
				"bsonType":             "object",
				"required":             []string{"username", "email", "password"},
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
						"bsonType":    "binData",
						"description": "user hashed password",
					},
					"isActive": bson.M{
						"bsonType":    "bool",
						"description": "user active status",
					},
					"isSeller": bson.M{
						"bsonType":    "bool",
						"description": "user seller status",
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
