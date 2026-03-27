package scraper_gateway

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceSpec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	Convey("Given start and end provided the quotes api client", t, func() {

		start := 1
		end := 10
		mockClient := NewMockServiceQuotesDiscoveryClient(ctrl)
		mockClient.EXPECT().
			TriggerScrape(gomock.Any(), start, end).
			Return(nil)
		service := NewService(mockClient)

		Convey("When the TriggerScrape is called", func() {
			err := service.TriggerScrape(context.Background(), start, end)
			Convey("Then it should trigger the call to the quots discovery client", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
