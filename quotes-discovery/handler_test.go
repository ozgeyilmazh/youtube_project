package quotesdiscovery_test

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	quotesdiscovery "github.com/ozgeyilmazh/youtube-project/quotes-discovery"
	"github.com/pact-foundation/pact-go/v2/provider"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	Convey("Given a HandlerWorkflow is provided", t, func() {
		mockWorkflow := quotesdiscovery.NewMockHandlerWorkflow(ctrl)
		mockWorkflow.EXPECT().FetchQuotes(gomock.Any(), gomock.Any()).Return([]quotesdiscovery.Quotes{
			{
				Quote:  "Test Quote",
				Author: "Test Author",
			},
		}, nil).AnyTimes()

		mockActivity := quotesdiscovery.NewMockWorkflowActivity(ctrl)
		mockActivity.EXPECT().FetchQuotes(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		mockActivity.EXPECT().BulkInsertData(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		verifier := provider.NewVerifier()

		app, w := quotesdiscovery.StartServer(mockWorkflow, mockActivity)
		defer w.Stop()

		go func() {
			if err := app.Listen(":4321"); err != nil {
				fmt.Println("Failed to start server", err)
			}
		}()

		if err := quotesdiscovery.WaitForServer("http://localhost:4321"); err != nil {
			t.Fatalf("Failed to wait for completion: %v", err)
		}

		Convey("When the Pact contract is verified", func() {

			err := verifier.VerifyProvider(t, provider.VerifyRequest{
				Provider:        "QuotesDiscovery",
				ProviderBaseURL: "http://localhost:4321",
				PactFiles: []string{
					"pacts/QuotesScraperGateway-QuotesDiscovery.json",
				},
			})
			Convey("Then the response should be 202 Accepted", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
