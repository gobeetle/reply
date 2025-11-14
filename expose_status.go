package reply

import (
	statuspkg "github.com/gobeetle/reply/internal/status"
)

type (
	Status = statuspkg.Status
)

var (
	NewStatus = statuspkg.New
)
