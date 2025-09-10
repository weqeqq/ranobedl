
mkdir build

for GOOS in darwin linux windows; do
  for GOARCH in amd64 arm64; do
    go build -o build/ranobedl-$GOOS-$GOARCH
  done
done
