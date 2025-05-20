#!/usr/bin/env bash
# Download a file from the server

# 设置显示错误信息
set -e

# 获取测试文件ID
TEST_FILE_1=$(curl -s "http://localhost:8888/api/v1/files/?tags=test&sortBy=name&order=asc" | jq -r .files[0].fileId)
echo "Testing with file ID: ${TEST_FILE_1}"

# 创建临时文件来存储错误信息
ERROR_FILE=$(mktemp)

# 下载最新版本的文件
echo "Downloading latest version..."
curl -X GET "http://localhost:8888/api/v1/files/download/${TEST_FILE_1}" \
    -H "Authorization: Bearer your_token_here" \
    -o downloaded_file.md \
    -v 2>"${ERROR_FILE}" || {
    echo "Download failed with error:"
    cat "${ERROR_FILE}"
}

# 下载指定版本的文件
echo "Downloading version 1..."
curl -X GET "http://localhost:8888/api/v1/files/download/${TEST_FILE_1}" \
    -H "Authorization: Bearer your_token_here" \
    -o downloaded_file_v1.md \
    -v 2>"${ERROR_FILE}" || {
    echo "Download failed with error:"
    cat "${ERROR_FILE}"
}

# 检查文件是否下载成功
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

# 清理
rm -f downloaded_file.md downloaded_file_v1.md "${ERROR_FILE}"
