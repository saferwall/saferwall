module github.com/saferwall/saferwall

go 1.14

require (
	github.com/aws/aws-sdk-go v1.35.1
	github.com/bnagy/gapstone v0.0.0-20190828052830-ede92aaeaba7
	github.com/golang/protobuf v1.4.1
	github.com/hillu/go-yara v1.2.2
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/minio/minio-go/v6 v6.0.57
	github.com/nsqio/go-nsq v1.0.8
	github.com/pelletier/go-toml v1.5.0 // indirect
	github.com/saferwall/saferwall/pkg/crypto v0.0.1
	github.com/saferwall/saferwall/pkg/peparser v0.0.19
	github.com/saferwall/saferwall/pkg/utils v0.0.4
	github.com/sirupsen/logrus v1.5.0
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.6.3
	github.com/stoewer/go-strcase v1.2.0
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79 // indirect
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b // indirect
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
)
