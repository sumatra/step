package machine

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"time"

	"github.com/coinbase/step/utils/to"
)

type HistoryEvent struct {
	types.HistoryEvent
}

type Execution struct {
	Output     map[string]interface{}
	OutputJSON string
	Error      error

	LastOutput     map[string]interface{} // interim output
	LastOutputJSON string
	LastError      error // interim error

	ExecutionHistory []HistoryEvent
}

func (sm *Execution) SetOutput(output interface{}, err error) {
	switch output.(type) {
	case map[string]interface{}:
		sm.Output = output.(map[string]interface{})
		sm.OutputJSON, _ = to.PrettyJSON(output)
	}

	if err != nil {
		sm.Error = err
	}
}

func (sm *Execution) SetLastOutput(output interface{}, err error) {
	switch output.(type) {
	case map[string]interface{}:
		sm.LastOutput = output.(map[string]interface{})
		sm.LastOutputJSON, _ = to.PrettyJSON(output)
	}

	if err != nil {
		sm.LastError = err
	}
}

func (sm *Execution) EnteredEvent(s State, input interface{}) {
	sm.ExecutionHistory = append(sm.ExecutionHistory, createEnteredEvent(s, input))
}

func (sm *Execution) ExitedEvent(s State, output interface{}) {
	sm.ExecutionHistory = append(sm.ExecutionHistory, createExitedEvent(s, output))
}

func (sm *Execution) Start() {
	sm.ExecutionHistory = []HistoryEvent{createEvent("ExecutionStarted")}
}

func (sm *Execution) Failed() {
	sm.ExecutionHistory = append(sm.ExecutionHistory, createEvent("ExecutionFailed"))
}

func (sm *Execution) Succeeded() {
	sm.ExecutionHistory = append(sm.ExecutionHistory, createEvent("ExecutionSucceeded"))
}

// Path returns the Path of States, ignoreing TaskFn states
func (sm *Execution) Path() []string {
	path := []string{}
	for _, er := range sm.ExecutionHistory {
		if er.StateEnteredEventDetails != nil {
			name := *er.StateEnteredEventDetails.Name
			path = append(path, name)
		}
	}
	return path
}

func createEvent(name string) HistoryEvent {
	t := time.Now()
	return HistoryEvent{
		types.HistoryEvent{
			Type:      types.HistoryEventType(name),
			Timestamp: &t,
		},
	}
}

func createEnteredEvent(state State, input interface{}) HistoryEvent {
	event := createEvent(fmt.Sprintf("%vStateEntered", *state.GetType()))
	json_raw, err := json.Marshal(input)

	if err != nil {
		json_raw = []byte{}
	}

	event.StateEnteredEventDetails = &types.StateEnteredEventDetails{
		Name:  state.Name(),
		Input: to.Strp(string(json_raw)),
	}

	return event
}

func createExitedEvent(state State, output interface{}) HistoryEvent {
	event := createEvent(fmt.Sprintf("%vStateExited", *state.GetType()))
	json_raw, err := json.Marshal(output)

	if err != nil {
		json_raw = []byte{}
	}

	event.StateExitedEventDetails = &types.StateExitedEventDetails{
		Name:   state.Name(),
		Output: to.Strp(string(json_raw)),
	}

	return event
}
