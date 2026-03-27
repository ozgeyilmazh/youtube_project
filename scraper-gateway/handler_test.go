package scraper_gateway

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	gomock "github.com/golang/mock/gomock"
	"github.com/pact-foundation/pact-go/v2/provider"
	. "github.com/smartystreets/goconvey/convey"
)

func StartServer(handler Handler) {

	app := fiber.New()
	handler.RegisterRoutes(app)
	if err := app.Listen(":1234"); err != nil {
		panic(err)
	}
}

func TestHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	Convey("Given a request to trigger scraping", t, func() {
		mockService := NewMockHandlerService(ctrl)
		mockService.EXPECT().TriggerScraping(gomock.Any(), 1, 10).Return(nil)
		verifier := provider.NewVerifier()
		handler := NewHandler(mockService)
		go StartServer(*handler)

		Convey("When the handler is called", func() {
			err := verifier.VerifyProvider(t, provider.VerifyRequest{
				Provider:        "QuotesScrapingGateway",
				ProviderBaseURL: "http://localhost:1234",
				PactFiles: []string{
					"pacts/QuotesScrapingScheduler-QuotesScrapingGateway.json",
				},
			})
			Convey("Then the response is successful", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
