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
	github.com/golang/snappy v0.0.1 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/minio/minio-go/v6 v6.0.39
	github.com/nsqio/go-nsq v1.0.7 // indirect
	github.com/pelletier/go-toml v1.5.0 // indirect
	github.com/saferwall/saferwall/pkg/crypto v0.0.1
	github.com/saferwall/saferwall/pkg/utils v0.0.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.4.0
	github.com/stoewer/go-strcase v1.0.2
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc // indirect
	golang.org/x/net v0.0.0-20191009170851-d66e71096ffb // indirect
	golang.org/x/sys v0.0.0-20191009170203-06d7bd2c5f4f // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20191009194640-548a555dbc03 // indirect
	google.golang.org/grpc v1.24.0
	gopkg.in/ini.v1 v1.48.0 // indirect
)
