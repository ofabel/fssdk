#!/bin/bash

set -e

cat << EOF > version.go
package main

const version = "$(git describe --tags --abbrev=0 master)"
EOF

for osArch in `cat ./build.txt`; do
    mkdir -p ./dist/temp

    IFS='/' read -ra parts <<< "$osArch"

    os=${parts[0]}
    arch=${parts[1]}

    echo ${osArch}

    GOOS=${os} GOARCH=${arch} go build -o ./dist/temp .

    cd ./dist/temp && zip "fssdk-${os}-${arch}.zip" ./* && cd ../..
    mv ./dist/temp/*.zip ./dist

    rm -rf ./dist/temp
done
