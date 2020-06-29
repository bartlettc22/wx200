docker run \
  -it \
  -v /home/chris/github.com/bartlettc22/wx200:/app \
  -w /app \
  --privileged \
  golang \
  sh -c "go get github.com/davecgh/go-spew/spew && go get github.com/jacobsa/go-serial/serial && go run main.go"
