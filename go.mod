module github.com/saferwall/saferwall

go 1.13

replace (
	github.com/saferwall/saferwall/pkg/crypto => ./pkg/crypto
	github.com/saferwall/saferwall/pkg/utils => ./pkg/utils
)

require (
	github.com/bitly/go-nsq v1.0.7
	github.com/bnagy/gapstone v0.0.0-20190828052830-ede92aaeaba7
	github.com/golang/protobuf v1.3.2
	github.com/golang/snappy v0.0.1 // indirect
	github.com/minio/minio-go/v6 v6.0.34
	github.com/nsqio/go-nsq v1.0.7 // indirect
	github.com/saferwall/saferwall/pkg/crypto v0.0.1
	github.com/saferwall/saferwall/pkg/utils v0.0.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.4.0
	github.com/stoewer/go-strcase v1.0.2
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a // indirect
	google.golang.org/grpc v1.21.0
)
