// models/userModel.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User mewakili struktur data pengguna dalam MongoDB
type User struct {
    ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Auth0ID              string             `bson:"auth0Id" json:"auth0Id"`
    Name                 string             `bson:"name" json:"name"`
    Email                string             `bson:"email" json:"email"`
    ProfilePicture       string             `bson:"profilePicture" json:"profilePicture"`
    LastSimulation       string             `bson:"lastSimulation" json:"lastSimulation"`
    CompletedSimulations []string           `bson:"completedSimulations" json:"completedSimulations"`
}
