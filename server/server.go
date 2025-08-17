package server

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Server struct {
	DB *sql.DB
	// Configuration fields
	SpotifyClientID     string
	SpotifyClientSecret string
	SpotifyRedirectURI  string
}

// SpotifyTokenResponse represents the response from Spotify's token endpoint
type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// SpotifyUserInfo represents the user info from Spotify's API
type SpotifyUserInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

// User model
type User struct {
	ID             int       `json:"id"`
	SpotifyID      string    `json:"spotify_id"`
	DisplayName    string    `json:"display_name"`
	Email          string    `json:"email"`
	AccessToken    string    `json:"-"`
	RefreshToken   string    `json:"-"`
	TokenExpiresAt time.Time `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (s *Server) InitDB() {
	db, err := sql.Open("sqlite", "./statify.db")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("DB connection created")

	s.DB = db
	s.createTables()
}

func (s *Server) createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	spotify_id TEXT UNIQUE NOT NULL,
	display_name TEXT,
	email TEXT,
	access_token TEXT,
	refresh_token TEXT,
	token_expires_at DATETIME,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createStatsTable := `
	CREATE TABLE IF NOT EXISTS user_stats (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER,
	artist_name TEXT,
	play_count INTEGER,
	last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users (id)
	);`

	var err error

	_, err = s.DB.Exec(createUsersTable)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = s.DB.Exec(createStatsTable)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("DB tables created")
}

func (s *Server) CloseDB() {
	if s.DB != nil {
		s.DB.Close()
	}
}

// LoadConfig loads configuration from environment variables
func (s *Server) LoadConfig() error {
	s.SpotifyClientID = os.Getenv("SPOTIFY_CLIENT_ID")
	s.SpotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	s.SpotifyRedirectURI = os.Getenv("SPOTIFY_REDIRECT_URI")

	// Validate required environment variables
	if s.SpotifyClientID == "" {
		return fmt.Errorf("SPOTIFY_CLIENT_ID environment variable is required")
	}
	if s.SpotifyClientSecret == "" {
		return fmt.Errorf("SPOTIFY_CLIENT_SECRET environment variable is required")
	}
	if s.SpotifyRedirectURI == "" {
		return fmt.Errorf("SPOTIFY_REDIRECT_URI environment variable is required")
	}

	log.Println("Configuration loaded successfully")
	return nil
}

// LoginUser handles the Spotify OAuth callback
func (s *Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received callback request: %s", r.URL.String())

	// Get the authorization code from URL parameters
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Println("No authorization code received")
		http.Error(w, "No authorization code", http.StatusBadRequest)
		return
	}

	// Exchange the code for access tokens
	tokenResponse, err := s.exchangeCodeForToken(code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	// Get user information from Spotify
	userInfo, err := s.getSpotifyUserInfo(tokenResponse.AccessToken)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Save or update user in database
	err = s.saveOrUpdateUser(userInfo, tokenResponse)
	if err != nil {
		log.Printf("Error saving user: %v", err)
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	log.Printf("User authenticated successfully: %s (%s)", userInfo.DisplayName, userInfo.ID)

	// Set a session cookie
	cookie := &http.Cookie{
		Name:     "spotify_user_id",
		Value:    userInfo.ID,
		Path:     "/",
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600,
	}

	http.SetCookie(w, cookie)
	log.Printf("Cookie set: %s=%s", cookie.Name, cookie.Value)

	// Redirect to the main app
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Helper methods
func (s *Server) exchangeCodeForToken(code string) (*SpotifyTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.SpotifyRedirectURI)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(s.SpotifyClientID+":"+s.SpotifyClientSecret)))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse SpotifyTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token response: %v", err)
	}

	return &tokenResponse, nil
}

func (s *Server) getSpotifyUserInfo(accessToken string) (*SpotifyUserInfo, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get user info with status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo SpotifyUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}

	return &userInfo, nil
}

func (s *Server) saveOrUpdateUser(userInfo *SpotifyUserInfo, tokenResponse *SpotifyTokenResponse) error {
	expiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	// Check if user exists
	var userID int
	err := s.DB.QueryRow("SELECT id FROM users WHERE spotify_id = ?", userInfo.ID).Scan(&userID)

	if err == sql.ErrNoRows {
		// Create new user
		_, err = s.DB.Exec(`
			INSERT INTO users (spotify_id, display_name, email, access_token, refresh_token, token_expires_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, userInfo.ID, userInfo.DisplayName, userInfo.Email, tokenResponse.AccessToken, tokenResponse.RefreshToken, expiresAt)
		return err
	} else if err != nil {
		return err
	}

	// Update existing user
	_, err = s.DB.Exec(`
		UPDATE users
		SET display_name = ?, email = ?, access_token = ?, refresh_token = ?, token_expires_at = ?, updated_at = CURRENT_TIMESTAMP
		WHERE spotify_id = ?
	`, userInfo.DisplayName, userInfo.Email, tokenResponse.AccessToken, tokenResponse.RefreshToken, expiresAt, userInfo.ID)

	return err
}

// GetUserBySpotifyID retrieves a user by their Spotify ID
func (s *Server) GetUserBySpotifyID(spotifyID string) (*User, error) {
	user := &User{}
	err := s.DB.QueryRow(`
		SELECT id, spotify_id, display_name, email, access_token, refresh_token, token_expires_at, created_at, updated_at
		FROM users WHERE spotify_id = ?
	`, spotifyID).Scan(
		&user.ID, &user.SpotifyID, &user.DisplayName, &user.Email,
		&user.AccessToken, &user.RefreshToken, &user.TokenExpiresAt,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
