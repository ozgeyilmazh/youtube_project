package quotesapi

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/pact-foundation/pact-go/v2/provider"
	. "github.com/smartystreets/goconvey/convey"
	gomock "go.uber.org/mock/gomock"
)

func startServer(handler *Handler) {
	app := fiber.New()
	handler.RegisterRoutes(app)
	if err := app.Listen(":1234"); err != nil {
		panic(err)
	}
}
func TestHandler(t *testing.T) {

	pactPath := "pacts/QuotesDiscovery-QuotesAPI.json"
	if _, err := os.Stat(pactPath); os.IsNotExist(err) {
		allPath := filepath.Join("quotes-api", pactPath)
		if _, err := os.Stat(allPath); os.IsNotExist(err) {
			t.Fatalf("Pact file not found: %s", allPath)
		}
		t.Fatalf("Pact file not found: %s", pactPath)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	Convey("Given a pact file and a verifier is provided", t, func() {
		mockService := NewMockHandlerService(ctrl)
		mockService.EXPECT().FetchQuotes(gomock.Any(), gomock.Any()).Return([]QuotesResponse{
			{Quote: "Test Quote", Author: "Test Author"},
		}, nil)

		verifier := provider.NewVerifier()
		handler := NewHandler(mockService)
		go startServer(handler)
		Convey("When the VerifiyProvider is executed", func() {
			err := verifier.VerifyProvider(t, provider.VerifyRequest{
				Provider:        "QuotesAPI",
				ProviderBaseURL: "http://localhost:1234",
				PactFiles: []string{
					pactPath,
				},
			})
			Convey("Then the error is nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
