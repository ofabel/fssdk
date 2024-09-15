.PHONY: install
update:
	git submodule update --remote ./flipperzero-protobuf && git add ./flipperzero-protobuf

.PHONY: clean
clean:
	rm -rf ./dist
	rm -rf ./flipper/rpc
	mkdir -p ./dist

.PHONY: format
format:
	gofmt -w .

.PHONY: run
run:
	go run .

.PHONY: build
build: clean protobuf
	go build -o ./dist .

.PHONY: protobuf
protobuf: clean
	protoc --proto_path=./flipperzero-protobuf \
		--go_out=./ \
		--go_opt=Mapplication.proto=github.com/ofabel/fssdk/flipper/rpc/flipper \
		--go_opt=Mdesktop.proto=github.com/ofabel/fssdk/flipper/rpc/desktop \
		--go_opt=Mflipper.proto=github.com/ofabel/fssdk/flipper/rpc/flipper \
		--go_opt=Mgpio.proto=github.com/ofabel/fssdk/flipper/rpc/gpio \
		--go_opt=Mgui.proto=github.com/ofabel/fssdk/flipper/rpc/gui \
		--go_opt=Mproperty.proto=github.com/ofabel/fssdk/flipper/rpc/property \
		--go_opt=Mstorage.proto=github.com/ofabel/fssdk/flipper/rpc/storage \
		--go_opt=Msystem.proto=github.com/ofabel/fssdk/flipper/rpc/system \
		./flipperzero-protobuf/*.proto
	mv ./github.com/ofabel/fssdk/flipper/rpc ./flipper
	rm -rf ./github.com
