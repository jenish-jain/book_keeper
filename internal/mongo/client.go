package mongo

import (
	"book_keeper/internal/logger"
	"bytes"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Client interface {
	InsertFile(ctx context.Context, name string, fileBytes []byte, fileType string) (string, error)
	DeleteFileByID(ctx context.Context, ID string) error
	GetAll(ctx context.Context) (interface{}, error)
}

type client struct {
	db *mongo.Database
}

type gridfsFile struct {
	ID         primitive.ObjectID     `bson:"_id"`
	Name       string                 `bson:"filename"`
	Length     int64                  `bson:"length"`
	UploadDate time.Time              `bson:"uploadDate"`
	Metadata   map[string]interface{} `bson:"metadata"`
}

func (d *client) GetAll(ctx context.Context) (interface{}, error) {
	bucket, err := gridfs.NewBucket(d.db)
	if err != nil {
		logger.ErrorWithCtx(ctx, "error creating bucket for fetching file", "event", "FETCH_FILES", "error", err)
	}

	filter := bson.D{}
	cursor, err := bucket.GetFilesCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var foundFiles []gridfsFile
	if err = cursor.All(context.TODO(), &foundFiles); err != nil {
		log.Fatal(err)
	}

	for _, file := range foundFiles {
		logger.Info("file data %+v", file)
	}
	return foundFiles, nil
}

func (d *client) InsertFile(ctx context.Context, name string, fileBytes []byte, fileType string) (string, error) {
	bucket, err := gridfs.NewBucket(d.db)
	if err != nil {
		logger.ErrorWithCtx(ctx, "error creating bucket for inserting file", "event", "INSERT_FILE", "error", err)
	}
	uploadOpts := options.GridFSUpload().
		SetMetadata(bson.D{{"type", fileType}})

	fileID, err := bucket.UploadFromStream(
		name,
		bytes.NewBuffer(fileBytes),
		uploadOpts)
	if err != nil {
		logger.ErrorWithCtx(ctx, "error uploading file to db", "event", "INSERT_FILE", "error", err)
		return "", err
	}

	return fileID.String(), nil
}

func (d *client) DeleteFileByID(ctx context.Context, ID string) error {
	bucket, err := gridfs.NewBucket(d.db)
	if err != nil {
		logger.ErrorWithCtx(ctx, "error creating bucket for deleting file", "event", "DELETE_FILE", "error", err)
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
