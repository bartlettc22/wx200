module github.com/bartlettc22/wx200

go 1.14

replace github.com/bartlettc22/pkg => ./pkg

require (
	github.com/jacobsa/go-serial v0.0.0-20180131005756-15cf729a72d4
	github.com/prometheus/client_golang v1.7.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
)
