// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	nsq "github.com/nsqio/go-nsq"
	"github.com/matcornic/hermes/v2"
	minio "github.com/minio/minio-go/v6"
	"github.com/saferwall/saferwall/pkg/utils"
	"github.com/saferwall/saferwall/web/app/common/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
)

type SMTPAuthentication struct {
	Server         string
	Port           int
	SenderEmail    string
	SenderIdentity string
	SMTPUser       string
	SMTPPassword   string
}


var (
	// RootDir points to the root dir
	RootDir string

	// StoragePath is where we save the samples
	StoragePath string

	// MaxFileSize allowed
	MaxFileSize int64

	// MaxAvatarFileSize allowed
	MaxAvatarFileSize int64

	// Domain name
	Domain string

	// NsqProducer holds an instance of NSQ producer.
	NsqProducer *nsq.Producer

	// MinioClient represents an instance of Object Space API client.
	MinioClient *minio.Client

	// UserSchema represent a user
	UserSchema *gojsonschema.Schema

	// FileSchema represent a user
	FileSchema *gojsonschema.Schema

	// LoginSchema represent a user login
	LoginSchema *gojsonschema.Schema

	// EmailSchema represent a request to change password
	EmailSchema *gojsonschema.Schema

	// ResetPasswordSchema represent a change password.
	ResetPasswordSchema *gojsonschema.Schema

	// FileActionSchema represent an action over a file.
	FileActionSchema *gojsonschema.Schema

	// UserUpdateSchema represent an update user.
	UserUpdateSchema *gojsonschema.Schema

	// PasswordUpdateSchema represent an update password.
	PasswordUpdateSchema *gojsonschema.Schema

	// EmailUpdateSchema represent an update email.
	EmailUpdateSchema *gojsonschema.Schema

	// UserActionSchema represent an action over a user.
	UserActionSchema *gojsonschema.Schema

	// CommentSchema represent an comment creation.
	CommentSchema *gojsonschema.Schema
	
	// SamplesSpaceBucket contains the space name of bucket to save samples.
	SamplesSpaceBucket string

	// AvatarSpaceBucket contains the space name of bucket to save user avatars.
	AvatarSpaceBucket string

	// Hermes represents an instance email generator.
	Hermes hermes.Hermes

	// SMTPConfig holds smtp config.
	SMTPConfig SMTPAuthentication

	// AvatarFileBuff contains the content of the default image avatar file path.
	AvatarFileBuff []byte

	// SfwAvatarFileDesc holds a descriptor to saferwall's image avatar file path.
	SfwAvatarFileDesc *os.File
)

// loadConfig loads our configration.
func loadConfig() {
	viper.AddConfigPath("config") // set the path of your config file

	// Load the config type depending on env variable.
	var name string
	env := os.Getenv("ENVIRONMENT")
	switch env {
	case "dev":
		name = "app.dev"
	case "prod":
		name = "app.prod"
	case "test":
		name = "app.test"
	default:
		log.Fatal("ENVIRONMENT is not set")
	}

	viper.SetConfigName(name)    // no need to include file extension
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Config %s was loaded", name)
}

// createNSQProducer creates a new NSQ producer.
func createNSQProducer() *nsq.Producer {

	// The default config settings provide a pretty good starting point for
	// our new consumer.
	config := nsq.NewConfig()

	// Create a NewProducer with the name of our topic, the channel, and our config
	addr := viper.GetString("nsq.address")
	p, err := nsq.NewProducer(addr, config)
	if err != nil {
		log.Fatal(err)
	}
	if p.Ping() != nil {
		log.Fatal(err)
	}

	log.Infoln("NSQ publisher was created")
	return p
}

// initLogging initialize our logging system.
func initLogging() {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	log.Infoln("Logger created")
}

// loadSchemas will load schemas at server startup.
func loadSchemas() {

	dir, err := utils.Getwd()
	if err != nil {
		log.Fatalln("Failed to Get current directory, err: ", err)
	}
	schemasPath := path.Join(dir, "app", "schema")

	// Get all schemas definitions
	var files []string
	err = filepath.Walk(schemasPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
        return nil
    })
    if err != nil {
		log.Fatalln("Failed to walk schemas directory, err: ", err)
	}
	
	// Iterate though json schemas and load them.
    for _, file := range files {
		jsonSchema := filepath.Base(file)
		source := fmt.Sprintf("file:///%s", file)
		jsonLoader := gojsonschema.NewReferenceLoader(source)

		switch jsonSchema{
		case "user.json":
			UserSchema, err = gojsonschema.NewSchema(jsonLoader)
		case "file.json":
			FileSchema, err = gojsonschema.NewSchema(jsonLoader)
		case "login.json":
			LoginSchema, err = gojsonschema.NewSchema(jsonLoader)
		case "email.json":
			EmailSchema, err = gojsonschema.NewSchema(jsonLoader)
		case "password.json":
			ResetPasswordSchema, err = gojsonschema.NewSchema(jsonLoader)
		case "file-action.json":
			FileActionSchema, err = gojsonschema.NewSchema(jsonLoader)
		case "update-user.json":
			UserUpdateSchema, err = gojsonschema.NewSchema(jsonLoader)
		case "update-password.json":
			PasswordUpdateSchema, err = gojsonschema.NewSchema(jsonLoader)
		case "update-email.json":
			EmailUpdateSchema, err = gojsonschema.NewSchema(jsonLoader)
		case "user-action.json":
			UserActionSchema, err = gojsonschema.NewSchema(jsonLoader)
		case "comment.json":
			CommentSchema, err = gojsonschema.NewSchema(jsonLoader)
		}

		if err != nil {
			log.Fatalf("Error while loading %s schema: ", jsonSchema, err)
		}
    }

	log.Infoln("Schemas were loaded")
}

// initOSClient returns a client for our Object Storage interface.
func initOSClient() *minio.Client {
	accessKey := viper.GetString("minio.accesskey")
	secKey := viper.GetString("minio.seckey")
	endpoint := viper.GetString("minio.endpoint")
	SamplesSpaceBucket = viper.GetString("minio.spacename")
	AvatarSpaceBucket = viper.GetString("minio.avatarspace")
	ssl := viper.GetBool("minio.ssl")
	location := "us-east-1"

	// Initiate a client using DigitalOcean Spaces.
	client, err := minio.New(endpoint, accessKey, secKey, ssl)
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Got Object Storage client instance")

	found, err := client.BucketExists(SamplesSpaceBucket)
	if err != nil {
		log.Fatalln(err)
	}
	if !found {
		err = client.MakeBucket(SamplesSpaceBucket, location)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Object storage Bucket %s exists already", SamplesSpaceBucket)
	}

	found, err = client.BucketExists(AvatarSpaceBucket)
	if err != nil {
		log.Fatalln(err)
	}
	if !found {
		err = client.MakeBucket(AvatarSpaceBucket, location)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Object storage bucket %s exists already", AvatarSpaceBucket)
	}

	return client
}

// initEmail initialize SMTP email config and setup hermes theme and product info.
func initHermes() {
	SMTPConfig = SMTPAuthentication{
		Server:         viper.GetString("smtp.server"),
		Port:           viper.GetInt("smtp.port"),
		SenderEmail:    viper.GetString("smtp.sender"),
		SenderIdentity: viper.GetString("smtp.identity"),
		SMTPPassword:   viper.GetString("smtp.password"),
		SMTPUser:       viper.GetString("smtp.user"),
	}
	Hermes = hermes.Hermes{
		Product: hermes.Product{
			Name:      "Saferwall",
			Link:      "https://saferwall.com/",
			Logo:      "https://i.imgur.com/zjCOKPo.png",
			Copyright: "Copyright © 2019 Saferwall. All rights reserved.",
		},
	}
	log.Println("Successfully created hermes instance")
}

// loadAvatars load default and saferwall images from file system.
func loadAvatars() {
	dir, err := utils.Getwd()
	if err != nil {
		log.Fatalln("Failed to Get current directory, err: ", err)
	}

	defaultAvatarPath := path.Join(dir, "data", "default-avatar.png")
	AvatarFileBuff, err = utils.ReadAll(defaultAvatarPath)
	if err != nil {
		log.Fatalf("Failed to open read avatar from %s, reason: %s",
		 defaultAvatarPath, err.Error())
	}

	sfwAvatarPath := path.Join(dir, "data", "saferwall-avatar.png")
	SfwAvatarFileDesc, err = os.Open(sfwAvatarPath)
	if err != nil {
		log.Fatalf("Failed to open saferwall avatar from %s, reason: %s",
		sfwAvatarPath, err.Error())
	}

	log.Println("Load Avatars success")
}

// Init will initiate required objects
func Init() {

	// Init our logger
	initLogging()

	// Load the configuration file
	loadConfig()

	// Load schemas
	loadSchemas()

	// Init email generator
	initHermes()

	// Get an instance of NSQ
	NsqProducer = createNSQProducer()

	// Connect to the database
	db.Connect()

	// Get an Object Storage client instance
	MinioClient = initOSClient()

	// Load default & saferwall image avatar.
	loadAvatars()

	Domain = viper.GetString("app.domain")
	StoragePath = viper.GetString("app.tmp_samples")
	MaxFileSize = int64(viper.GetInt("app.max_file_size"))
	MaxAvatarFileSize = int64(viper.GetInt("app.max_avatar_file_size"))
	os.MkdirAll(StoragePath, os.ModePerm)

}
