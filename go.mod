module ahtapot

go 1.24.3

require github.com/gocql/gocql v1.7.0

require github.com/inconshreveable/mousetrap v1.1.0 // indirect

require (
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/spf13/cobra v1.9.1
	github.com/spf13/pflag v1.0.6 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

replace github.com/gocql/gocql => github.com/scylladb/gocql v1.15.0
