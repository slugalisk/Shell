default: clean go

go:
	mkdir -p ./go;
	protoc -I ./ *.proto --go_out=plugins=grpc:go;

clean:
	rm -rf ./go;
