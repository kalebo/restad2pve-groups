
build:
	CGO_ENABLED=0 go build -v -a -tags netgo -installsuffix netgo
