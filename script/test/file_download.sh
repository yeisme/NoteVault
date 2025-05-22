#!/usr/bin/env bash
# Download a file from the server

# Set to display error messages
set -e

# Get test file ID
TEST_FILE_1=$(curl -s "http://localhost:8888/api/v1/files/?tags=test&sortBy=name&order=asc" | jq -r .files[0].fileId)
echo "Testing with file ID: ${TEST_FILE_1}"

# Create a temporary file to store error messages
ERROR_FILE=$(mktemp)

# Download the latest version of the file
echo "Downloading latest version..."
curl -X GET "http://localhost:8888/api/v1/files/download/${TEST_FILE_1}" \
    -H "Authorization: Bearer your_token_here" \
    -o downloaded_file.md \
    -v 2>"${ERROR_FILE}" || {
    echo "Download failed with error:"
    cat "${ERROR_FILE}"
}

# Download a specific version of the file
echo "Downloading version 1..."
curl -X GET "http://localhost:8888/api/v1/files/download/${TEST_FILE_1}" \
    -H "Authorization: Bearer your_token_here" \
    -o downloaded_file_v1.md \
    -v 2>"${ERROR_FILE}" || {
    echo "Download failed with error:"
    cat "${ERROR_FILE}"
}

# Check if the files were downloaded successfully
if [ -f downloaded_file.md ] && [ -s downloaded_file.md ]; then
    echo "Latest version downloaded successfully"
else
    echo "Failed to download latest version"
fi

if [ -f downloaded_file_v1.md ] && [ -s downloaded_file_v1.md ]; then
    echo "Version 1 downloaded successfully"
else
    echo "Failed to download version 1"
fi

# Cleanup
rm -f downloaded_file.md downloaded_file_v1.md "${ERROR_FILE}"
