package quotesapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSchedulerConfig(t *testing.T) {
	Convey("Given config is loaded", t, func() {
		config := NewConfig(EnvironmentLocal)
		Convey("When config GetAllEnv is called", func() {
			allEnv := config.GetAllEnv()
			Convey("Then all env is returned", func() {
				So(allEnv, ShouldContainKey, "QUOTES_API_URL")
			})
		})
	})

}
