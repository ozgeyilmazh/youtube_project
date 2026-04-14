package quotesdiscovery

import (
	"context"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func getTestQueue() string {
	return "quotes-discovery-test"
}

func setupWorkflowOptions() client.StartWorkflowOptions {
	return client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: getTestQueue(),
	}
}
func setupTemporalClient() client.Client {
	temporalClient, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		panic(err)
	}
	return temporalClient
}

func setupWorkerClient(temporalClient client.Client, taskQueue string, workflow HandlerWorkflow, activity WorkflowActivity) worker.Worker {
	worker := worker.New(temporalClient, taskQueue, worker.Options{})
	worker.RegisterWorkflow(workflow.FetchQuotes)
	worker.RegisterActivity(activity.FetchQuotes)
	worker.RegisterActivity(activity.BulkInsertData)
	err := worker.Start()
	if err != nil {
		panic(err)
	}
	return worker
}

func TestWorkflow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	Convey("Given a fetchQuotes activity is provided", t, func() {
		mockActivity := NewMockWorkflowActivity(ctrl)
		mockActivity.EXPECT().FetchQuotes(gomock.Any(), 1).Return([]Quotes{
			{
				Quote:  "Test Quote",
				Author: "Test Author",
			},
		}, nil).Times(1)
		mockActivity.EXPECT().BulkInsertData(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		workflowOptions := setupWorkflowOptions()
		temporalClient := setupTemporalClient()
		activityOptions := workflow.ActivityOptions{
			StartToCloseTimeout: 10 * time.Second,
		}
		config := NewConfig(EnvironmentLocal)
		workflow := NewWorkflow(mockActivity, activityOptions, config)
		worker := setupWorkerClient(temporalClient, getTestQueue(), workflow, mockActivity)
		defer worker.Stop()

		Convey("When the workflow is executed", func() {
			flowRun, err := temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, workflow.FetchQuotes, 1)
			So(err, ShouldBeNil)
			var got []Quotes
			err = flowRun.Get(context.Background(), &got)
			Convey("Then FetchQuotes and BulkInsertData each run once and the workflow returns the quotes", func() {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, []Quotes{{Quote: "Test Quote", Author: "Test Author"}})
			})
		})
	})
}
