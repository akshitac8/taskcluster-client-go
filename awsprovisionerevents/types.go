// This source code file is AUTO-GENERATED by github.com/taskcluster/jsonschema2go

package awsprovisionerevents

type (
	// Message reporting that an action occured to a worker type
	WorkerTypeMessage struct {
		Version float64 `json:"version"`

		// Name of the worker type which was created
		WorkerType string `json:"workerType"`
	}
)