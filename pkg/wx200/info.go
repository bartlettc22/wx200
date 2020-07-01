package wx200

type Info struct {
	// SamplesRecieved is a counter of the total number of samples received
	SamplesRecieved uint64

	// ChecksumFailures is a counter of the total number of checksum failures
	ChecksumFailures uint64
}
