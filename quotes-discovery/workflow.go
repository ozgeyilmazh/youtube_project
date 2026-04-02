package quotesdiscovery

//go:generate mockgen -destination=workflow_mock.go -package=quotesdiscovery . WorkflowActivity
import (
	"context"
	"fmt"

	"go.temporal.io/sdk/workflow"
)

type WorkflowActivity interface {
	FetchQuotes(ctx context.Context, page int) ([]Quotes, error)
	BulkInsertData(ctx context.Context, quotes []Quotes) error
}

type Workflow struct {
	activity        WorkflowActivity
	activityOptions workflow.ActivityOptions
	config          *Config
}

func NewWorkflow(activity WorkflowActivity, activityOptions workflow.ActivityOptions, config *Config) *Workflow {
	return &Workflow{activity: activity, activityOptions: activityOptions, config: config}
}

func (w *Workflow) FetchQuotes(ctx workflow.Context, page int) ([]Quotes, error) {
	return nil, fmt.Errorf("not implemented")
}

func (w *Workflow) BulkInsertData(ctx workflow.Context, quotes []Quotes) error {
	return fmt.Errorf("not implemented")
}
