#!/bin/bash

go build -o server cmd/server/main.go
if [ $? -ne 0 ]; then
  echo "Build failed"
  exit 1
fi

echo "Build successful, uploading to server..."

scp server livekit-1:~/arca-v3/.
if [ $? -ne 0 ]; then
  echo "Failed to copy server binary to livekit-1"
  exit 1
fi

echo "Server binary uploaded successfully"
rm -rf server