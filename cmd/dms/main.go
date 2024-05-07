package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/joshua468/document-management-system/internal/auth"
	"github.com/joshua468/document-management-system/internal/document"
	"github.com/joshua468/document-management-system/internal/user"
)

func main() {
	db, err := gorm.Open(sqlite.Open("dms.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	err = db.AutoMigrate(&user.User{}, &document.Document{})
	if err != nil {
		log.Fatal("Error migrating database schema:", err)
	}
	userRepo := user.NewRepository(db)
	docRepo := document.NewRepository(db)
	authService := auth.NewService(userRepo)

	router := mux.NewRouter()
	router.HandleFunc("/login", auth.HandleLogin(authService)).Methods("POST")
	router.HandleFunc("/documents", authMiddleware(authService, document.HandleDocument(docRepo))).Methods("GET")
	router.HandleFunc("/documents", authMiddleware(authService, document.CreateDocument(docRepo))).Methods("POST")
	router.HandleFunc("/documents/{id}", authMiddleware(authService, document.GetDocument(docRepo))).Methods("GET")
	router.HandleFunc("/documents/{id}", authMiddleware(authService, document.UpdateDocument(docRepo))).Methods("PUT")
	router.HandleFunc("/documents/{id}", authMiddleware(authService, document.DeleteDocument(docRepo))).Methods("DELETE")

	fmt.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func authMiddleware(authService *auth.Service, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := authService.ValidateToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := auth.SetUserIDInContext(r.Context(), userID)
		next(w, r.WithContext(ctx))
	}
}
