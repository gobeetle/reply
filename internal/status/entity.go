package status

// Status represents the status of the API
type Status string

var EnumStatus = struct {
	StatusSuccess    Status
	StatusError      Status
	StatusProcessing Status
	StatusPending    Status
	StatusFailed     Status
	StatusOther      Status
}{
	"success",
	"error",
	"processing",
	"pending",
	"failed",
	"other",
}

// NewStatus returns the Status based on the status code
func New(statusCode int) Status {
	switch statusCode / 100 {
	case 2: // 2xx Success
		return EnumStatus.StatusSuccess
	case 4: // 4xx Client Error
		return EnumStatus.StatusError
	case 5: // 5xx Server Error
		return EnumStatus.StatusFailed
	default: // Other status codes
		return EnumStatus.StatusOther
	}
}

func (s Status) String() string {
	return string(s)
}
