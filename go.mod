module github.com/saferwall/saferwall

go 1.13

replace (
	github.com/saferwall/saferwall/pkg/crypto => ./pkg/crypto
	github.com/saferwall/saferwall/pkg/utils => ./pkg/utils
)

require (
	github.com/bitly/go-nsq v1.0.7
	github.com/bnagy/gapstone v0.0.0-20190828052830-ede92aaeaba7
	github.com/fatih/color v1.7.0
	github.com/golang/protobuf v1.3.2
	github.com/jessevdk/go-flags v1.4.0
	github.com/minio/minio-go/v6 v6.0.34
	github.com/saferwall/saferwall/pkg/crypto v0.0.1
	github.com/saferwall/saferwall/pkg/utils v0.0.1
	github.com/saferwall/saferwall/web v0.0.0-20190930033257-cd60aab01c85 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.4.0
	github.com/stoewer/go-strcase v1.0.2
	google.golang.org/grpc v1.21.0
)
