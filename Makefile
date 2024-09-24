.PHONY: install
update:
	git submodule update --remote ./flipperzero-protobuf && git add ./flipperzero-protobuf

.PHONY: clean
clean:
	rm -rf ./dist
	rm -rf ./rpc/protobuf
	mkdir -p ./dist

.PHONY: format
format:
	gofmt -w .

.PHONY: run
run:
	go run .

.PHONY: build
build: clean protobuf
	GOOS=linux go build -o ./dist .

.PHONY: protobuf
protobuf: clean
	protoc --proto_path=./flipperzero-protobuf \
		--go_out=./ \
		--go_opt=Mapplication.proto=github.com/ofabel/fssdk/rpc/protobuf/flipper \
		--go_opt=Mdesktop.proto=github.com/ofabel/fssdk/rpc/protobuf/desktop \
		--go_opt=Mflipper.proto=github.com/ofabel/fssdk/rpc/protobuf/flipper \
		--go_opt=Mgpio.proto=github.com/ofabel/fssdk/rpc/protobuf/gpio \
		--go_opt=Mgui.proto=github.com/ofabel/fssdk/rpc/protobuf/gui \
		--go_opt=Mproperty.proto=github.com/ofabel/fssdk/rpc/protobuf/property \
		--go_opt=Mstorage.proto=github.com/ofabel/fssdk/rpc/protobuf/storage \
		--go_opt=Msystem.proto=github.com/ofabel/fssdk/rpc/protobuf/system \
		./flipperzero-protobuf/*.proto
	mv ./github.com/ofabel/fssdk/rpc/protobuf ./rpc
	rm -rf ./github.com
