package quotesdiscovery

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestActivity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	Convey("Given a quotes discovery activity is provided", t, func() {
		mockClient := NewMockActivityQuotesAPIClient(ctrl)
		mockRepository := NewMockRepository(ctrl)
		mockClient.EXPECT().FetchQuotes(gomock.Any(), gomock.Any()).Return([]Quotes{
			{
				Quote:  "Test Quote",
				Author: "Test Author",
			},
		}, nil).AnyTimes()
		mockRepository.EXPECT().BulkInsertData(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		activity := NewActivity(mockClient, mockRepository)
		Convey("When the activity is executed", func() {
			quotes, err := activity.FetchQuotes(context.Background(), 1)
			if err != nil {
				t.Fatalf("Failed to fetch quotes: %v", err)
			}
			Convey("Then the quotes are fetched", func() {
				So(err, ShouldBeNil)
				So(quotes, ShouldNotBeNil)
				So(len(quotes), ShouldBeGreaterThan, 0)
			})
		})
	})
}
