package shared

type ProvisionInput struct {
	AppName   string `json:"app_name"`
	Namespace string `json:"namespace"`
	Cluster   string `json:"cluster"`
}

type StepStatus string

const (
	StepPending   StepStatus = "pending"
	StepRunning   StepStatus = "running"
	StepCompleted StepStatus = "completed"
	StepFailed    StepStatus = "failed"
	StepRollback  StepStatus = "rollback"
)

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
