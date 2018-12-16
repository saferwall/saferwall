package app

import (
	"os"

	"github.com/spf13/viper"
)

var (
	// StoragePath is where we save the samples
	StoragePath string

	// MaxFileSize allowed
	MaxFileSize int64
)

// Init will create some directories
func Init() {

	StoragePath = viper.GetString("storage.tmp_samples")
	MaxFileSize = int64(viper.GetInt("storage.max_file_size"))
	os.MkdirAll(StoragePath, os.ModePerm)
}
