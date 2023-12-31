package main

import (
	"book_keeper/internal/config"
	"bytes"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
)

func main() {
	// MongoDB's connection string
	// Uses the SetServerAPIOptions() method to set the Stable API version to 1
	configStore := config.InitConfig("production")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mongoconfigs := configStore.MongoConfig
	mongoUri := fmt.Sprintf("mongodb+srv://%s:%s@%s/", mongoconfigs.RWUser, mongoconfigs.RWPassword, mongoconfigs.Host)
	clientOptions := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Create a new MongoDB database and collection
	db := client.Database("book_keeper")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		log.Fatal(err)
	}

	// Read a PDF file
	pdfFile, err := ioutil.ReadFile(getAssetsPathName("example.pdf"))
	if err != nil {
		log.Fatal(err)
	}

	// Store the PDF file in MongoDB
	uploadStream, err := bucket.OpenUploadStream("example.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(pdfFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Uploaded a file of size: %d\n", fileSize)

	// Retrieve the PDF file from MongoD

	var buf bytes.Buffer
	fileName := "example.pdf"
	dStream, err := bucket.DownloadToStreamByName(fileName, &buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File size to download: %v\n", dStream)
	downloadFileName := "example-download.pdf"
	err = ioutil.WriteFile(getAssetsPathName(downloadFileName), buf.Bytes(), 0600)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Downloaded the file from MongoDB")
}

func getAssetsPathName(name string) string {
	return fmt.Sprintf("assets/%s", name)
}
