package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type SignupResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

type UserToolsRequest struct {
	Username string   `json:"username"`
	Tools    []string `json:"tools"`
}

type UserToolsResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Tools   []string `json:"tools,omitempty"`
}

var db *sql.DB

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s - From: %s", r.Method, r.RequestURI, r.RemoteAddr, r.Header.Get("User-Agent"))
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s - Completed in %v", r.Method, r.RequestURI, time.Since(start))
	})
}

func initDB() error {
	var err error
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://loginapp:loginapp123@localhost:5432/loginapp?sslmode=disable"
	}

	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	// Create users table
	createUsersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createUsersTableQuery)
	if err != nil {
		return err
	}

	// Create user_tools table
	createUserToolsTableQuery := `
	CREATE TABLE IF NOT EXISTS user_tools (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		tool_name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(username, tool_name)
	);
	`

	_, err = db.Exec(createUserToolsTableQuery)
	if err != nil {
		return err
	}

	log.Println("✅ Database initialized successfully")
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var signupReq SignupRequest
	err := json.NewDecoder(r.Body).Decode(&signupReq)
	if err != nil {
		log.Printf("❌ Signup failed - Invalid request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SignupResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if signupReq.Username == "" || signupReq.Password == "" {
		log.Printf("❌ Signup failed - Missing username or password")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SignupResponse{
			Success: false,
			Message: "Username and password are required",
		})
		return
	}

	// Check if user already exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", signupReq.Username).Scan(&exists)
	if err != nil {
		log.Printf("❌ Signup failed - Database error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SignupResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	if exists {
		log.Printf("❌ Signup failed - Username already exists: %s", signupReq.Username)
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(SignupResponse{
			Success: false,
			Message: "Username already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := hashPassword(signupReq.Password)
	if err != nil {
		log.Printf("❌ Signup failed - Password hashing error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SignupResponse{
			Success: false,
			Message: "Error processing password",
		})
		return
	}

	// Insert user
	var user User
	err = db.QueryRow(
		"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username, created_at",
		signupReq.Username, hashedPassword,
	).Scan(&user.ID, &user.Username, &user.CreatedAt)

	if err != nil {
		log.Printf("❌ Signup failed - Error creating user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SignupResponse{
			Success: false,
			Message: "Error creating user",
		})
		return
	}

	log.Printf("✅ Signup successful - Username: %s (ID: %d)", user.Username, user.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SignupResponse{
		Success: true,
		Message: "User created successfully",
		User:    &user,
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		log.Printf("❌ Login failed - Invalid request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if loginReq.Username == "" {
		log.Printf("❌ Login failed - Empty username")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Username is required",
		})
		return
	}

	// If password is provided, check against database
	if loginReq.Password != "" {
		var user User
		err = db.QueryRow("SELECT id, username, password FROM users WHERE username=$1", loginReq.Username).
			Scan(&user.ID, &user.Username, &user.Password)

		if err == sql.ErrNoRows {
			log.Printf("❌ Login failed - User not found: %s", loginReq.Username)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(LoginResponse{
				Success: false,
				Message: "Invalid credentials",
			})
			return
		} else if err != nil {
			log.Printf("❌ Login failed - Database error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(LoginResponse{
				Success: false,
				Message: "Database error",
			})
			return
		}

		if !checkPasswordHash(loginReq.Password, user.Password) {
			log.Printf("❌ Login failed - Invalid password for user: %s", loginReq.Username)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(LoginResponse{
				Success: false,
				Message: "Invalid credentials",
			})
			return
		}

		log.Printf("✅ Login successful - Username: %s (ID: %d)", user.Username, user.ID)
	} else {
		// Legacy: accept any username without password
		log.Printf("✅ Login successful (legacy) - Username: %s", loginReq.Username)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   "demo-token-" + loginReq.Username,
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

func getUserToolsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Printf("❌ Get tools failed - Missing username")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UserToolsResponse{
			Success: false,
			Message: "Username is required",
		})
		return
	}

	rows, err := db.Query("SELECT tool_name FROM user_tools WHERE username=$1", username)
	if err != nil {
		log.Printf("❌ Get tools failed - Database error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserToolsResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}
	defer rows.Close()

	var tools []string
	for rows.Next() {
		var tool string
		if err := rows.Scan(&tool); err != nil {
			log.Printf("❌ Get tools failed - Scan error: %v", err)
			continue
		}
		tools = append(tools, tool)
	}

	// If no tools found, return empty array
	if tools == nil {
		tools = []string{}
	}

	log.Printf("✅ Get tools successful - Username: %s, Tools: %v", username, tools)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserToolsResponse{
		Success: true,
		Message: "Tools retrieved successfully",
		Tools:   tools,
	})
}

func saveUserToolsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var toolsReq UserToolsRequest
	err := json.NewDecoder(r.Body).Decode(&toolsReq)
	if err != nil {
		log.Printf("❌ Save tools failed - Invalid request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UserToolsResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if toolsReq.Username == "" {
		log.Printf("❌ Save tools failed - Missing username")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UserToolsResponse{
			Success: false,
			Message: "Username is required",
		})
		return
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("❌ Save tools failed - Transaction error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserToolsResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	// Delete existing tools for user
	_, err = tx.Exec("DELETE FROM user_tools WHERE username=$1", toolsReq.Username)
	if err != nil {
		tx.Rollback()
		log.Printf("❌ Save tools failed - Delete error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserToolsResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	// Insert new tools
	for _, tool := range toolsReq.Tools {
		_, err = tx.Exec("INSERT INTO user_tools (username, tool_name) VALUES ($1, $2)", toolsReq.Username, tool)
		if err != nil {
			tx.Rollback()
			log.Printf("❌ Save tools failed - Insert error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(UserToolsResponse{
				Success: false,
				Message: "Database error",
			})
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("❌ Save tools failed - Commit error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserToolsResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	log.Printf("✅ Save tools successful - Username: %s, Tools: %v", toolsReq.Username, toolsReq.Tools)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserToolsResponse{
		Success: true,
		Message: "Tools saved successfully",
		Tools:   toolsReq.Tools,
	})
}

func main() {
	// Initialize database
	if err := initDB(); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/api/login", loginHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/signup", signupHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/user/tools", getUserToolsHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/user/tools", saveUserToolsHandler).Methods("POST", "OPTIONS")

	// Add logging middleware
	r.Use(loggingMiddleware)

	log.Println("🚀 Server starting on port 8080...")
	log.Println("📝 Logging enabled for all requests")
	log.Println("💾 Database connection established")
	log.Fatal(http.ListenAndServe(":8080", r))
}
