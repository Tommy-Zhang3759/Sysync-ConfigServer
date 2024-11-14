#!/bin/bash

BUILD_DIR="./build"
DATA_DIR="./data"
WEBUI_DIR="./webUI/static"
OUTPUT_BIN="$BUILD_DIR/ConfigServer"
VERSION=$(git describe --tags --abbrev=0)

rm -rf $BUILD_DIR

echo "Fetching dependencies..."
go get ConfigServer/utils/database

echo "Building the project with version $VERSION..."
go build -x -ldflags "-X 'main.Version=$VERSION'" -o $OUTPUT_BIN

echo "Creating build directories..."
mkdir -p $BUILD_DIR/data
mkdir -p $BUILD_DIR/webUI

echo "Copying database files..."
cp $DATA_DIR/clientInfoTemp.db $BUILD_DIR/data/clientInfo.db

echo "Copying static files..."
cp -r $WEBUI_DIR $BUILD_DIR/webUI/

echo "Build completed successfully!"
