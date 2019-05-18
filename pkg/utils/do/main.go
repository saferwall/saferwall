package do

import (
	minio "github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// loadConfig loads our configration.
func loadConfig() {
	viper.SetConfigName("saferwall") // no need to include file extension
	viper.AddConfigPath(".")         // set the path of your config file

	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
		panic(err)
	}
}

// GetClient returns a client instance for DigitalOcean Spaces.
func GetClient() *minio.Client {
	loadConfig()
	accessKey := viper.GetString("do.accesskey")
	secKey := viper.GetString("do.seckey")
	endpoint := viper.GetString("do.endpoint")
	ssl := viper.GetBool("do.ssl")

	// Initiate a client using DigitalOcean Spaces.
	client, err := minio.New(endpoint, accessKey, secKey, ssl)
	if err != nil {
		log.Errorln(err)
	}

	log.Infoln("GetClient success")
	return client
}
