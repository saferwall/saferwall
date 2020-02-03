// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"os"
	"path"

	nsq "github.com/bitly/go-nsq"
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

	// Debug mode
	Debug bool

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

	// FileActionSchema represent a change password.
	FileActionSchema *gojsonschema.Schema

	// UserUpdateSchema represent an update user.
	UserUpdateSchema *gojsonschema.Schema

	// PasswordUpdateSchema represent an update password.
	PasswordUpdateSchema *gojsonschema.Schema

	// EmailUpdateSchema represent an update email.
	EmailUpdateSchema *gojsonschema.Schema

	// SamplesSpaceBucket contains the space name of bucket to save samples.
	SamplesSpaceBucket string

	// AvatarSpaceBucket contains the space name of bucket to save user avatars.
	AvatarSpaceBucket string

	// Hermes represents an instance email generator.
	Hermes hermes.Hermes

	// SMTPConfig holds smtp config.
	SMTPConfig SMTPAuthentication

	// AvatarFileDesc holds a descriptor to the default image avatar file path.
	AvatarFileDesc *os.File

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

	jsonPath := path.Join(dir, "app", "schema", "user.json")
	source := fmt.Sprintf("file:///%s", jsonPath)
	jsonLoader := gojsonschema.NewReferenceLoader(source)
	UserSchema, err = gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		log.Fatalln("Error while loading user schema: ", err)
	}

	jsonPath = path.Join(dir, "app", "schema", "file.json")
	source = fmt.Sprintf("file:///%s", jsonPath)
	jsonLoader = gojsonschema.NewReferenceLoader(source)
	FileSchema, err = gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		log.Fatalln("Error while loading file schema: ", err)
	}

	jsonPath = path.Join(dir, "app", "schema", "login.json")
	source = fmt.Sprintf("file:///%s", jsonPath)
	jsonLoader = gojsonschema.NewReferenceLoader(source)
	LoginSchema, err = gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		log.Fatalln("Error while loading file schema: ", err)
	}

	jsonPath = path.Join(dir, "app", "schema", "email.json")
	source = fmt.Sprintf("file:///%s", jsonPath)
	jsonLoader = gojsonschema.NewReferenceLoader(source)
	EmailSchema, err = gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		log.Fatalln("Error while loading email schema: ", err)
	}

	jsonPath = path.Join(dir, "app", "schema", "password.json")
	source = fmt.Sprintf("file:///%s", jsonPath)
	jsonLoader = gojsonschema.NewReferenceLoader(source)
	ResetPasswordSchema, err = gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		log.Fatalln("Error while loading email schema: ", err)
	}

	jsonPath = path.Join(dir, "app", "schema", "action.json")
	source = fmt.Sprintf("file:///%s", jsonPath)
	jsonLoader = gojsonschema.NewReferenceLoader(source)
	FileActionSchema, err = gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		log.Fatalln("Error while loading file actions schema: ", err)
	}

	jsonPath = path.Join(dir, "app", "schema", "update-user.json")
	source = fmt.Sprintf("file:///%s", jsonPath)
	jsonLoader = gojsonschema.NewReferenceLoader(source)
	UserUpdateSchema, err = gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		log.Fatalln("Error while loading user update schema: ", err)
	}

	jsonPath = path.Join(dir, "app", "schema", "update-password.json")
	source = fmt.Sprintf("file:///%s", jsonPath)
	jsonLoader = gojsonschema.NewReferenceLoader(source)
	PasswordUpdateSchema, err = gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		log.Fatalln("Error while loading password update schema: ", err)
	}

	jsonPath = path.Join(dir, "app", "schema", "update-email.json")
	source = fmt.Sprintf("file:///%s", jsonPath)
	jsonLoader = gojsonschema.NewReferenceLoader(source)
	EmailUpdateSchema, err = gojsonschema.NewSchema(jsonLoader)
	if err != nil {
		log.Fatalln("Error while loading email update schema: ", err)
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

	bucketName := SamplesSpaceBucket
	err = client.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := client.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s", bucketName)
	}
	
	bucketName = AvatarSpaceBucket
	err = client.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := client.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s", bucketName)
	}
	log.Infoln("Got Object Storage client instance")
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
	AvatarFileDesc, err = os.Open(defaultAvatarPath)
	if err != nil {
		log.Fatalf("Failed to open default avatar from %s, reason: %s",
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

	Debug = viper.GetBool("app.debug")
	StoragePath = viper.GetString("app.tmp_samples")
	MaxFileSize = int64(viper.GetInt("app.max_file_size"))
	MaxAvatarFileSize = int64(viper.GetInt("app.max_avatar_file_size"))
	os.MkdirAll(StoragePath, os.ModePerm)

}
