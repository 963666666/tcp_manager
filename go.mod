module tcp_manager

go 1.12

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.46.3
	cloud.google.com/go/bigquery => github.com/kikinteractive/go-bqstreamer v2.0.1+incompatible
	cloud.google.com/go/datastore => github.com/gomods/athens v0.6.1
	cloud.google.com/go/pubsub => github.com/cskr/pubsub v1.0.2
	contrib.go.opencensus.io/exporter/stackdriver => github.com/frodenas/stackdriver_exporter v0.6.0
	go.opencensus.io => github.com/egymgmbh/opencensus-go-exporter-influxdb v0.0.0-20190125120553-91f0f40078a8
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20191002192127-34f69633bfdc
	golang.org/x/exp => github.com/golang/exp v0.0.0-20191002040644-a1355ae1e2c3
	golang.org/x/image => github.com/golang/image v0.0.0-20191009234506-e7c1f5e7dbb8
	golang.org/x/lint => github.com/golang/lint v0.0.0-20190930215403-16217165b5de
	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20191002175909-6d0d39b2ca82
	golang.org/x/mod => github.com/golang/mod v0.1.0
	golang.org/x/net => github.com/golang/net v0.0.0-20191009170851-d66e71096ffb
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys => github.com/golang/sys v0.0.0-20191009170203-06d7bd2c5f4f
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/time => github.com/golang/time v0.0.0-20190921001708-c4c64cad1fd0
	golang.org/x/tools => github.com/golang/tools v0.0.0-20191010075000-0337d82405ff
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20190717185122-a985d3407aa7
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.11.0
	google.golang.org/appengine => github.com/golang/appengine v1.6.5
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20191009194640-548a555dbc03
	google.golang.org/grpc => github.com/grpc/grpc-go v1.24.0
	gopkg.in/DataDog/dd-trace-go.v1 => github.com/DataDog/dd-trace-go v1.18.0
)

require (
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.9
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.6.1 // indirect
	golang.org/x/tools v0.0.0-20190927191325-030b2cf1153e
)
