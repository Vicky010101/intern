package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func init() {
	var err error
	minioClient, err = minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minio-access-key", "minio-secret-key", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %s", err)
	}
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Upload the file to MinIO (or S3)
	objectName := "uploaded-file.jpg"
	bucketName := "user-files"
	contentType := "image/jpeg"

	_, err = minioClient.PutObject(r.Context(), bucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading file: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully!")
}
