build:
	GOOS=linux GOARCH=amd64 go build  -a -ldflags '-extldflags "-static"' -o storejson .
	docker build -t inge4pres/storejson:0.2	 .

push: build
	docker push inge4pres/storejson:0.2

.PHONY: build push
