package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/mattn/go-sqlite3"
)

// AWS S3 bucket configuration
var (
	s3Client   *s3.Client
	bucketName = "your-bucket-name" // Replace with your bucket name
	region     = "us-east-1"        // Replace with your AWS region
)

// Initialize database and S3 client
func init() {
	// Load S3 configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)

	// Initialize database
	db, err := sql.Open("sqlite3", "./filedata.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create files table if not exists
	createTable := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT NOT NULL,
		size INTEGER NOT NULL,
		upload_date DATETIME NOT NULL,
		url TEXT NOT NULL
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

// uploadHandler handles file uploads
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the uploaded file
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "File size too large", http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate a unique file key
	fileKey := fmt.Sprintf("%d-%s", time.Now().Unix(), handler.Filename)

	// Upload file to S3
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
		Body:   file,
	})
	if err != nil {
		http.Error(w, "Failed to upload file to S3", http.StatusInternalServerError)
		return
	}

	// File URL
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, fileKey)

	// Save metadata in the database
	db, err := sql.Open("sqlite3", "./filedata.db")
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO files (filename, size, upload_date, url)
		VALUES (?, ?, ?, ?);
	`, handler.Filename, handler.Size, time.Now(), fileURL)
	if err != nil {
		http.Error(w, "Failed to save metadata", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully! URL: %s", fileURL)
}

func main() {
	// File upload route
	http.HandleFunc("/upload", uploadHandler)

	// Start server
	log.Println("Server is running on http://localhost:8080/auth/login")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
