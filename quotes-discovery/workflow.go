package quotesdiscovery

//go:generate mockgen -destination=workflow_mock.go -package=quotesdiscovery . WorkflowActivity
import (
	"context"

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

	ctx = workflow.WithActivityOptions(ctx, w.activityOptions)
	var quotes []Quotes
	err := workflow.ExecuteActivity(ctx, w.activity.FetchQuotes, page).Get(ctx, &quotes)
	if err != nil {
		return nil, err
	}
	err = workflow.ExecuteActivity(ctx, w.activity.BulkInsertData, quotes).Get(ctx, nil)
	if err != nil {
		return nil, err
	}
	return quotes, nil
}
