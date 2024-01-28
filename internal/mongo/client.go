package mongo

import (
	"book_keeper/internal/logger"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Client interface {
	InsertFile(ctx context.Context, fileBytes []byte, name string) (int, error)
	DeleteFileByID(ctx context.Context, ID string) error
}

type client struct {
	db *mongo.Database
}

func (d *client) InsertFile(ctx context.Context, fileBytes []byte, name string) (int, error) {
	bucket, err := gridfs.NewBucket(d.db)
	if err != nil {
		logger.ErrorWithCtx(ctx, "error creating bucket for inserting file", "event", "INSERT_FILE", "error", err)
	}

	uploadStream, err := bucket.OpenUploadStream(name)
	if err != nil {
		logger.ErrorWithCtx(ctx, "error creating upload stream for inserting file", "event", "INSERT_FILE", "error", err)
	}
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(fileBytes)
	if err != nil {
		logger.ErrorWithCtx(ctx, "error writing file to db", "event", "INSERT_FILE", "error", err)

	}
	logger.DebugWithCtx(ctx, "Uploaded a file of size: %d", fileSize, "event", "INSERT_FILE")
	return fileSize, err
}

func (d *client) DeleteFileByID(ctx context.Context, ID string) error {
	bucket, err := gridfs.NewBucket(d.db)
	if err != nil {
		logger.ErrorWithCtx(ctx, "error creating bucket for inserting file", "event", "INSERT_FILE", "error", err)
	}
	fileId, err := primitive.ObjectIDFromHex(ID)
	err = bucket.Delete(fileId)
	if err != nil {
		return err
	}
	return nil
}

func NewClient(config *Config) Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mongoUri := fmt.Sprintf("mongodb+srv://%s:%s@%s/", config.rwUsername, config.rwPassword, config.host)
	clientOptions := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)

	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	//defer mongoClient.Disconnect(context.Background())
	// Create a new MongoDB database and collection
	db := mongoClient.Database(config.databaseName)
	return &client{
		db: db,
	}
}
