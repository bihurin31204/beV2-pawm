// handlers/userHandler.go
package handlers

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "vlab-backend/models"

    jwt "github.com/golang-jwt/jwt/v4"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func UserHandler(collection *mongo.Collection) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
        defer cancel()

        // Ambil token dari context
        tokenRaw := r.Context().Value("user")
        if tokenRaw == nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        token, ok := tokenRaw.(*jwt.Token)
        if !ok {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Invalid claims", http.StatusUnauthorized)
            return
        }

        auth0ID, ok := claims["sub"].(string)
        if !ok {
            http.Error(w, "Invalid subject", http.StatusUnauthorized)
            return
        }

        switch r.Method {
        case "GET":
            var user models.User
            err := collection.FindOne(ctx, bson.M{"auth0Id": auth0ID}).Decode(&user)
            if err != nil {
                if err == mongo.ErrNoDocuments {
                    http.Error(w, "User not found", http.StatusNotFound)
                } else {
                    log.Printf("Error finding user with auth0Id %s: %v", auth0ID, err)
                    http.Error(w, "Internal server error", http.StatusInternalServerError)
                }
                return
            }
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(user)
        case "POST":
            var user models.User
            err := json.NewDecoder(r.Body).Decode(&user)
            if err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
            }

            // Pastikan auth0ID konsisten
            user.Auth0ID = auth0ID

            filter := bson.M{"auth0Id": user.Auth0ID}
            update := bson.M{"$set": user}
            opts := options.Update().SetUpsert(true)
            _, err = collection.UpdateOne(ctx, filter, update, opts)
            if err != nil {
                log.Printf("Error saving user with auth0Id %s: %v", user.Auth0ID, err)
                http.Error(w, "Error saving user", http.StatusInternalServerError)
                return
            }
            w.WriteHeader(http.StatusOK)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }
}
