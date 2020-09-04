module github.com/ekara-platform/cli

go 1.13

require (
	docker.io/go-docker v1.0.0
	github.com/GroupePSA/componentizer v0.0.0-20200904074711-6001dbb137ca
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/coreos/go-etcd v2.0.0+incompatible // indirect
	github.com/cpuguy83/go-md2man v1.0.10 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/ekara-platform/engine v1.0.1-0.20200227174114-77fea9e5b9bd
	github.com/ekara-platform/model v1.0.1-0.20200227174022-a451d2e5d22b
	github.com/fatih/color v1.9.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/ugorji/go/codec v0.0.0-20181204163529-d75b2dcb6bc8 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	gopkg.in/src-d/go-git.v4 v4.13.1 // indirect
)

replace github.com/ekara-platform/engine => ../engine
