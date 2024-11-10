package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "github.com/gorilla/mux"
    gorillaHandlers "github.com/gorilla/handlers" // Berikan alias untuk menghindari konflik
    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "vlab-backend/handlers"      // Mengacu pada package handlers untuk route
    "vlab-backend/middleware"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Tidak dapat memuat file .env, menggunakan variabel lingkungan sistem")
    }

    mongoURI := os.Getenv("MONGODB_URI")
    if mongoURI == "" {
        log.Fatal("MONGODB_URI tidak ditemukan dalam variabel lingkungan")
    }

    clientOptions := options.Client().ApplyURI(mongoURI)
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err = client.Disconnect(context.Background()); err != nil {
            log.Fatal(err)
        }
    }()

    db := client.Database("salsaDB")
    userCollection := db.Collection("stateRecord")

    r := mux.NewRouter()
    r.Handle("/api/userstate", middleware.AuthMiddleware(handlers.UserHandler(userCollection))).Methods("GET", "POST")

    // Konfigurasi CORS dengan gorillaHandlers
    corsHandler := gorillaHandlers.CORS(
        gorillaHandlers.AllowedOrigins([]string{"http://localhost:3000"}),      // Izinkan asal frontend
        gorillaHandlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),      // Izinkan metode tertentu
        gorillaHandlers.AllowedHeaders([]string{"Authorization", "Content-Type"}), // Izinkan header tertentu
    )

    log.Println("Server berjalan pada port 8000")
    if err := http.ListenAndServe(":8000", corsHandler(r)); err != nil {
        log.Fatal("Gagal menjalankan server:", err)
    }
}
