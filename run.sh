docker run \
  -it \
  -v $(pwd):/app \
  -w /app \
  -p 9041:9041 \
  --privileged \
  golang \
  bash
