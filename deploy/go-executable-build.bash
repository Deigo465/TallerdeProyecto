#!/usr/bin/env bash

package_name=$1
if [[ -z "$package_name" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64")
# Override to build for only one platform
platforms=("linux/amd64")

for platform in "${platforms[@]}"
do
    echo "Building for" + $platform
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    extra_env=''
    if [ $GOOS = "linux" ]; then
        extra_env='CC=x86_64-unknown-linux-gnu-gcc CXX=x86_64-unknown-linux-gnu-g++' # go build -ldflags "-linkmode external -extldflags -static"
    fi
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
        extra_env="CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++"
    fi

    env $extra_env CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o ./$output_name ./cmd/web-server
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
    # env $extra_env CGO_ENABLED=1 GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name'-seed' ./cmd/seed
    # if [ $? -ne 0 ]; then
        # echo 'An error has occurred! Aborting the script execution...'
        # # exit 1
    # fi
done



# taken from https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04
