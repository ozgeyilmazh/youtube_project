package scrape_scheduler

import (
	"context"
	"log/slog"
	"os"
	"testing"

	gomock "github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestQuotesScrapeActivity(t *testing.T) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("Given a quotes scrape start at 1 end at 10 is scheduled", t, func() {
		start := 1
		end := 10
		mockClient := NewMockClient(ctrl)
		mockClient.EXPECT().TriggerScraping(gomock.Any(), start, end).Return(nil)
		activity := QuotesScrapeActivity{
			logger: logger,
			client: mockClient,
		}
		Convey("When the activity is executed", func() {
			err := activity.BeginScrape(context.Background(), start, end)
			Convey("Then the activity is scheduled", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
