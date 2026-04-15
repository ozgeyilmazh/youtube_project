package quotesapi

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	gomock "go.uber.org/mock/gomock"
)

func TestService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	Convey("Given a quotes API service is provided", t, func() {
		service := NewMockQuotesAPIClient(ctrl)
		service.EXPECT().FetchQuotes(gomock.Any(), gomock.Any()).Return([]QuotesResponse{
			{Quote: "Test Quote", Author: "Test Author"},
		}, nil)

		Convey("When the FetchQuotes is called", func() {
			quotes, err := service.FetchQuotes(context.Background(), 1)

			Convey("Then the quotes are fetched", func() {
				So(err, ShouldBeNil)
				So(quotes, ShouldNotBeNil)
				So(len(quotes), ShouldBeGreaterThan, 0)
			})
		})
	})
}
