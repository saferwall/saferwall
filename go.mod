module github.com/saferwall/saferwall

go 1.13

replace (
	github.com/saferwall/saferwall/pkg/crypto => ./pkg/crypto
	github.com/saferwall/saferwall/pkg/utils => ./pkg/utils
)

require (
	github.com/bitly/go-nsq v1.0.7
	github.com/bnagy/gapstone v0.0.0-20190828052830-ede92aaeaba7
	github.com/edsrzf/mmap-go v1.0.0
	github.com/fatih/color v1.7.0
	github.com/golang/protobuf v1.3.2
	github.com/hillu/go-yara v1.1.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/kjk/betterguid v0.0.0-20170621091430-c442874ba63a // indirect
	github.com/labstack/echo/v4 v4.1.13 // indirect
	github.com/minio/minio-go/v6 v6.0.44
	github.com/saferwall/saferwall/pkg/crypto v0.0.1
	github.com/saferwall/saferwall/pkg/utils v0.0.3
	github.com/saferwall/saferwall/web v0.0.0-20200120221542-5362c0e04bdb // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.1
	github.com/stoewer/go-strcase v1.1.0
	go.mozilla.org/pkcs7 v0.0.0-20200128120323-432b2356ecb1
	google.golang.org/grpc v1.25.1
)
