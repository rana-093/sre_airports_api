package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	region = "ap-south-1"
)

type Airport struct {
	Name     string `json:"name"`
	City     string `json:"city"`
	IATA     string `json:"iata"`
	ImageURL string `json:"image_url"`
}

type AirportV2 struct {
	Airport
	RunwayLength int `json:"runway_length"`
}

// Mock data for airports in Bangladesh
var airports = []Airport{
	{"Hazrat Shahjalal International Airport", "Dhaka", "DAC", "https://storage.googleapis.com/bd-airport-data/dac.jpg"},
	{"Shah Amanat International Airport", "Chittagong", "CGP", "https://storage.googleapis.com/bd-airport-data/cgp.jpg"},
	{"Osmani International Airport", "Sylhet", "ZYL", "https://storage.googleapis.com/bd-airport-data/zyl.jpg"},
}

// Mock data for airports in Bangladesh (with runway length for V2)
var airportsV2 = []AirportV2{
	{Airport{"Hazrat Shahjalal International Airport", "Dhaka", "DAC", "https://storage.googleapis.com/bd-airport-data/dac.jpg"}, 3200},
	{Airport{"Shah Amanat International Airport", "Chittagong", "CGP", "https://storage.googleapis.com/bd-airport-data/cgp.jpg"}, 2900},
	{Airport{"Osmani International Airport", "Sylhet", "ZYL", "https://storage.googleapis.com/bd-airport-data/zyl.jpg"}, 2500},
}

// HomePage handler
func HomePage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status: OK"))
}

// Airports handler for the first endpoint
func Airports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(airports)
}

// AirportsV2 handler for the second version endpoint
func AirportsV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(airportsV2)
}

func createAndGetNewAWSSession() (*session.Session, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretAccessKey, ""),
	})
	return sess, err
}

func UpdateAirportImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	sess, err := createAndGetNewAWSSession()
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	svc := s3.New(sess)

	airportName := r.FormValue("airportName")
	file, handler, err := r.FormFile("airportImage")
	fileExtension := filepath.Ext(handler.Filename)

	slog.Info("Got request for airport with name %s%s", airportName, fileExtension)

	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}

	defer file.Close()

	indexOfTheAirport := -1
	for idx, airport := range airportsV2 {
		if strings.EqualFold(airport.Name, airportName) {
			indexOfTheAirport = idx
			break
		}
	}
	if indexOfTheAirport == -1 {
		slog.Info("Airport not found")
		http.Error(w, "Airport not found", http.StatusNotFound)
		return
	}

	bucketName := os.Getenv("AWS_BUCKET_NAME")

	fileName := fmt.Sprintf("airports/%s%s", uuid.New(), fileExtension)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		slog.Info("err is: ", err)
		http.Error(w, "Failed to upload image to S3", http.StatusInternalServerError)
		return
	}

	imageURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, fileName)

	airports[indexOfTheAirport].ImageURL = imageURL

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Image successfully updated",
		"image_url": imageURL,
	})
}

func main() {

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/airports", Airports)
	http.HandleFunc("/airports_v2", AirportsV2)
	http.HandleFunc("/update_airport_image", UpdateAirportImage)

	slog.Info("Gonna listen to 8080 port!")

	http.ListenAndServe(":8080", nil)
}
