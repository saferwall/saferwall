module github.com/saferwall/saferwall/web

go 1.13

replace github.com/saferwall/saferwall/pkg/utils => ../pkg/utils

require (
	github.com/Pallinder/go-randomdata v1.2.0
	github.com/bitly/go-nsq v1.0.7
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/golang/snappy v0.0.1 // indirect
	github.com/labstack/echo/v4 v4.1.10
	github.com/matcornic/hermes/v2 v2.0.2
	github.com/minio/minio-go/v6 v6.0.34
	github.com/nsqio/go-nsq v1.0.7 // indirect
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/saferwall/saferwall/pkg/crypto v0.0.1
	github.com/saferwall/saferwall/pkg/utils v0.0.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.4.0
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.1.0
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/couchbase/gocb.v1 v1.6.2
	gopkg.in/couchbase/gocbcore.v7 v7.1.13 // indirect
	gopkg.in/couchbaselabs/gocbconnstr.v1 v1.0.4 // indirect
	gopkg.in/couchbaselabs/gojcbmock.v1 v1.0.3 // indirect
	gopkg.in/couchbaselabs/jsonx.v1 v1.0.0 // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
)
