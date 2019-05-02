docker run --rm -it -v "$PWD":/go/src/github.com/ekara-platform/cli -w /go/src/github.com/ekara-platform/cli dockercore/golang-cross:1.11.0 sh -c '
    for GOOS in darwin linux windows; do
      for GOARCH in amd64; do
        echo "Building $GOOS-$GOARCH"
        export GOOS=$GOOS
        export GOARCH=$GOARCH
        go build -o ekara-$GOOS
      done
    done
    '
