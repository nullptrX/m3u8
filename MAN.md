
# intel
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dump github.com/nullptrx/v2

# m1
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dump github.com/nullptrx/v2

darwin/amd64

darwin/arm64



lipo -archs dump


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dump github.com/nullptrx/v2
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dump github.com/nullptrx/v2
