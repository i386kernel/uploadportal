build:
	GOOS=linux GOARCH=amd64 go build -o refarchapp .
	chmod 777 refarchapp

imgbuild:
	docker build -t i386kernel/ra-portal:"${IMG_VER}" .

imgpush:
	docker push i386kernel/ra-portal:${IMG_VER}

