package status

import "time"

type Status struct {
	Version string        `json:"version"`                      // Version of the service
	Uptime  time.Duration `json:"uptime" swaggertype:"integer"` // Indicates the duration of the service uptime
}
