package shared

type ComputeInput struct {
	InitialValue int `json:"initial_value"`
}

type StepStatus string

const (
	StepPending   StepStatus = "pending"
	StepRunning   StepStatus = "running"
	StepCompleted StepStatus = "completed"
	StepFailed    StepStatus = "failed"
	StepRollback  StepStatus = "rollback"
)

type Compensation struct {
	Name string      `json:"name"`
	Do   interface{} `json:"fn"`
}

type WorkflowStep struct {
	Key         string     `json:"key"`
	Label       string     `json:"label"`
	Status      StepStatus `json:"status"`
	Message     string     `json:"message,omitempty"`
	StartedAt   string     `json:"started_at,omitempty"`
	CompletedAt string     `json:"completed_at,omitempty"`
}

type ProgressState struct {
	WorkflowID  string         `json:"workflow_id,omitempty"`
	Status      string         `json:"status"`
	CurrentStep string         `json:"current_step,omitempty"`
	Error       string         `json:"error,omitempty"`
	Steps       []WorkflowStep `json:"steps"`
}

type ActivityResult struct {
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

type MockError struct {
	Message string `json:"message"`
}

func (e *MockError) Error() string {
	return e.Message
}
