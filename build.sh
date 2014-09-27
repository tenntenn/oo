#!/bin/sh

pushd `dirname $0` > /dev/null

APP_NAME='oo'
APP_VERSION=`cat version.txt`
APP_OS="linux darwin windows"
APP_ARCH="386 amd64"
LDFLAGS="-X oo.Version $APP_VERSION"

# goxが必要
# https://github.com/mitchellh/gox
GOPATH=$GOPATH:`pwd` gox -os="$APP_OS" -arch="$APP_ARCH" -output bin/"v"$APP_VERSION"_{{.OS}}_{{.Arch}}/"$APP_NAME -ldflags "$LDFLAGS" main

popd > /dev/null
