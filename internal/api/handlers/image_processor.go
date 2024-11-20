
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"net/url"
	"os"
	"zocket/internal/db"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"github.com/nfnt/resize"
)


func ProcessImage(imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		log.Printf("Failed to download image %s: %v", imageURL, err)
		return "", err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Printf("Failed to decode image %s: %v", imageURL, err)
		return "", err
	}

	resizedImg := resize.Resize(800, 0, img, resize.Lanczos3)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resizedImg, nil)
	if err != nil {
		log.Printf("Failed to compress image %s: %v", imageURL, err)
		return "", err
	}
	s3URL, err := uploadToS3(buf.Bytes(), imageURL)
	if err != nil {
		log.Printf("Failed to upload image %s to S3: %v", imageURL, err)
		return "", err
	}

	return s3URL, nil
}

func uploadToS3(imageBytes []byte, imageURL string) (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	if awsAccessKeyID == "" || awsSecretAccessKey == "" || awsRegion == "" {
		log.Fatalf("Missing AWS credentials in environment variables")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKeyID, awsSecretAccessKey, "",
		),
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		return "", err
	}

	svc := s3.New(sess)

	firstEncodedURL := url.QueryEscape(imageURL)
	doubleEncodedURL := url.QueryEscape(firstEncodedURL)
	fileName := fmt.Sprintf("compressed-%s", doubleEncodedURL)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String("golangbucket-yolo"),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(imageBytes),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		log.Printf("Failed to upload file to S3: %v", err)
		return "", err
	}
	filee := url.QueryEscape(fileName)
	s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", "golangbucket-yolo", awsRegion, filee)
	log.Printf("Successfully uploaded compressed image to S3: %s", s3URL)

	return s3URL, nil
}

func StartImageProcessingWorker() {
	msgs, err := RabbitMQChannel.Consume(
		"imageProcessingQueue", 
		"",                    
		true,                   
		false,                  
		false,                  
		false,                 
		nil,                    
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Waiting for messages...")

	for d := range msgs {
		var payload map[string]string
		if err := json.Unmarshal(d.Body, &payload); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		imageURL, ok := payload["image_url"]
		if !ok {
			log.Printf("Invalid payload, missing 'image_url': %v", payload)
			continue
		}

		compressedURL, err := ProcessImage(imageURL)
		if err != nil {
			log.Printf("Error processing image %s: %v", imageURL, err)
			continue
		}
		err = updateCompressedImageURL(imageURL, compressedURL)
		if err != nil {
			log.Printf("Failed to update compressed image URL in database: %v", err)
		} else {
			log.Printf("Successfully updated compressed image URL for %s", imageURL)
		}
		log.Printf("Successfully processed image %s", imageURL)
	}
}

func updateCompressedImageURL(originalURL, compressedURL string) error {
	query := `
        UPDATE products
        SET compressed_product_images = array_append(compressed_product_images, $1)
        WHERE $2 = ANY(product_images);
    `
	_, err := db.DB.Exec(context.Background(), query, compressedURL, originalURL)
	return err
}
