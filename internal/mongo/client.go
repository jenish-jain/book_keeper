package mongo

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jenish-jain/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	InsertFile(name string, fileBytes []byte, fileType string) (string, error)
	DeleteFileByID(ID string) error
	GetAllFiles(ctx context.Context) (interface{}, error)
}

type client struct {
	db         *mongo.Database
	fileBucket *gridfs.Bucket
}

type gridfsFile struct {
	ID         primitive.ObjectID     `bson:"_id"`
	Name       string                 `bson:"filename"`
	Length     int64                  `bson:"length"`
	UploadDate time.Time              `bson:"uploadDate"`
	Metadata   map[string]interface{} `bson:"metadata"`
}

func (d *client) GetAllFiles(ctx context.Context) (interface{}, error) {
	filter := bson.D{}
	cursor, err := d.fileBucket.GetFilesCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var foundFiles []gridfsFile
	if err = cursor.All(context.TODO(), &foundFiles); err != nil {
		log.Fatal(err)
	}

	for _, file := range foundFiles {
		logger.DebugWithCtx(ctx, "found files with data", "fileData", file, "event", "FETCH_FILES")
	}
	return foundFiles, nil
}

func (d *client) InsertFile(name string, fileBytes []byte, fileType string) (string, error) {
	uploadOpts := options.GridFSUpload().
		SetMetadata(bson.D{{"type", fileType}})

	fileID, err := d.fileBucket.UploadFromStream(
		name,
		bytes.NewBuffer(fileBytes),
		uploadOpts)
	if err != nil {
		return "", err
	}

	return fileID.Hex(), nil
}

func (d *client) DeleteFileByID(ID string) error {
	fileId, err := primitive.ObjectIDFromHex(ID)
	err = d.fileBucket.Delete(fileId)
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
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		logger.Error("error creating bucket for fetching file", "event", "NEW_CLIENT", "error", err)
	}
	return &client{
		db:         db,
		fileBucket: bucket,
	}
}
