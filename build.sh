#!/bin/bash
BINARY_NAME="moddownloader-go"
OUTPUT_DIR="./builds"
mkdir -p "$OUTPUT_DIR"
PLATFORMS=(
    "windows/amd64"
    "linux/amd64"
)
for PLATFORM in "${PLATFORMS[@]}"; do
    IFS="/" read -r GOOS GOARCH <<< "$PLATFORM"
    OUTPUT_NAME="$OUTPUT_DIR/${BINARY_NAME}-$GOOS-$GOARCH"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME+=".exe"
    fi

    echo "Building for $GOOS/$GOARCH..."
    env GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$OUTPUT_NAME" .
    if [ $? -ne 0 ]; then
        echo "An error occurred while building for $GOOS/$GOARCH"
        exit 1
    fi
done

echo "Builds completed successfully."