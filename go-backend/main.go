package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"go-backend/models"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/stripe/stripe-go"
)

var s *Server

type Server struct {
	db             *sql.DB
	s3             *s3.Client
	testing        bool
	jwt_secret_key []byte
	stripe_key     string
	mail           *MailClient
	TestInspector  *TestInspector
}

type MailClient struct {
	Host     string
	Password string
}

type TestInspector struct {
	EmailsSent    int
	FilesUploaded int
}

func admin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("current_user").(int)
		user, err := s.QueryUser(userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		if !user.IsAdmin {
			http.Error(w, "Access denied", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")

		if tokenStr == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		tokenStr = tokenStr[len("Bearer "):]

		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return s.jwt_secret_key, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add the claims to the request context
		ctx := context.WithValue(r.Context(), "current_user", claims.Sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

type Email struct {
	Subject   string `json:"subject"`
	Recipient string `json:"recipient"`
	Body      string `json:"body"`
}

func (s *Server) SendEmail(subject, recipient, body string) error {
	if s.testing {
		s.TestInspector.EmailsSent += 1
		return nil
	}
	email := Email{
		Subject:   subject,
		Recipient: recipient,
		Body:      body,
	}

	// Convert email struct to JSON

	emailJSON, err := json.Marshal(email)
	if err != nil {
		return err
	}
	go func() {

		// Create a new request
		req, err := http.NewRequest("POST", s.mail.Host+"/api/send", bytes.NewBuffer(emailJSON))
		if err != nil {
			return
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", s.mail.Password)

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			log.Printf("failed to send email: %s", resp.Status)
			return
		}
	}()
	return nil
}

func main() {
	s = &Server{}

	dbConfig := models.DatabaseConfig{}
	dbConfig.Host = os.Getenv("DB_HOST")
	dbConfig.Port = os.Getenv("DB_PORT")
	dbConfig.User = os.Getenv("DB_USER")
	dbConfig.Password = os.Getenv("DB_PASS")
	dbConfig.DatabaseName = os.Getenv("DB_NAME")

	db, err := ConnectToDatabase(dbConfig)

	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}
	s.db = db
	s.runMigrations()
	s.s3 = s.createS3Client()

	s.stripe_key = os.Getenv("STRIPE_SECRET_KEY")
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	s.mail = &MailClient{
		Host:     os.Getenv("MAIL_HOST"),
		Password: os.Getenv("MAIL_PASSWORD"),
	}
	s.jwt_secret_key = []byte(os.Getenv("SECRET_KEY"))

	r := mux.NewRouter()
	r.HandleFunc("/api/auth", jwtMiddleware(s.CheckTokenRoute)).Methods("GET")
	r.HandleFunc("/api/login", s.LoginRoute).Methods("POST")
	r.HandleFunc("/api/reset-password", s.ResetPasswordRoute).Methods("POST")
	r.HandleFunc("/api/email-validate", jwtMiddleware(s.ResendEmailValidationRoute)).Methods("GET")
	r.HandleFunc("/api/email-validate", s.ValidateEmailRoute).Methods("POST")
	r.HandleFunc("/api/request-reset", s.RequestPasswordResetRoute).Methods("POST")

	r.HandleFunc("/api/files", jwtMiddleware(s.GetAllFilesRoute)).Methods("GET")
	r.HandleFunc("/api/files/upload", jwtMiddleware(s.UploadFileRoute)).Methods("POST")
	r.HandleFunc("/api/files/{id}", jwtMiddleware(s.GetFileMetadataRoute)).Methods("GET")
	r.HandleFunc("/api/files/{id}", jwtMiddleware(s.EditFileMetadataRoute)).Methods("PATCH")
	r.HandleFunc("/api/files/{id}", jwtMiddleware(s.DeleteFileRoute)).Methods("DELETE")
	r.HandleFunc("/api/files/download/{id}", jwtMiddleware(s.DownloadFileRoute)).Methods("GET")

	r.HandleFunc("/api/cards", jwtMiddleware(s.GetCardsRoute)).Methods("GET")
	r.HandleFunc("/api/cards", jwtMiddleware(s.CreateCardRoute)).Methods("POST")
	r.HandleFunc("/api/next", jwtMiddleware(s.NextIDRoute)).Methods("POST")
	r.HandleFunc("/api/cards/{id}", jwtMiddleware(s.GetCardRoute)).Methods("GET")
	r.HandleFunc("/api/cards/{id}", jwtMiddleware(s.UpdateCardRoute)).Methods("PUT")
	r.HandleFunc("/api/cards/{id}", jwtMiddleware(s.DeleteCardRoute)).Methods("DELETE")

	r.HandleFunc("/api/users/{id}", jwtMiddleware(admin(s.GetUserRoute))).Methods("GET")
	r.HandleFunc("/api/users/{id}", jwtMiddleware(s.UpdateUserRoute)).Methods("PUT")
	r.HandleFunc("/api/users", jwtMiddleware(admin(s.GetUsersRoute))).Methods("GET")
	r.HandleFunc("/api/users", s.CreateUserRoute).Methods("POST")
	r.HandleFunc("/api/users/{id}/subscription", jwtMiddleware(admin(s.GetUserSubscriptionRoute))).Methods("GET")
	r.HandleFunc("/api/current", jwtMiddleware(s.GetCurrentUserRoute)).Methods("GET")
	r.HandleFunc("/api/admin", jwtMiddleware(s.GetUserAdminRoute)).Methods("GET")

	r.HandleFunc("/api/tasks/{id}", jwtMiddleware(s.GetTaskRoute)).Methods("GET")
	r.HandleFunc("/api/tasks", jwtMiddleware(s.GetTasksRoute)).Methods("GET")
	r.HandleFunc("/api/tasks", jwtMiddleware(s.CreateTaskRoute)).Methods("POST")
	r.HandleFunc("/api/tasks/{id}", jwtMiddleware(s.UpdateTaskRoute)).Methods("PUT")
	r.HandleFunc("/api/tasks/{id}", jwtMiddleware(s.DeleteTaskRoute)).Methods("DELETE")

	r.HandleFunc("/api/billing/create_checkout_session", s.CreateCheckoutSession).Methods("POST")
	r.HandleFunc("/api/billing/success", s.GetSuccessfulSessionData).Methods("GET")
	r.HandleFunc("/api/webhook", s.HandleWebhook).Methods("POST")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})
	//handler := cors.Default().Handler(r)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ZETTEL_URL")},
		AllowCredentials: true,
		AllowedHeaders:   []string{"authorization", "content-type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		// Enable Debugging for testing, consider disabling in production
		//Debug: true,
	})

	handler := c.Handler(r)
	http.ListenAndServe(":8080", handler)
}
