echo "List of files"
curl "http://localhost:8888/api/v1/files/"

echo "List of files with fileName=test2.md"
curl "http://localhost:8888/api/v1/files/?fileName=test2.md"

echo "List of files with fileName=test2.md and fileType=text/markdown"
curl "http://localhost:8888/api/v1/files/?fileName=test2.md&fileType=text/markdown"

echo "List of files with tags=test"
curl "http://localhost:8888/api/v1/files/?tags=test&sortBy=name&order=asc"
