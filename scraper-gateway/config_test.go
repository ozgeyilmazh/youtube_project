package scraper_gateway

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSchedulerConfig(t *testing.T) {
	Convey("Given dev config file is loaded", t, func() {
		config := NewConfig(EnvironmentDev)
		Convey("Then temporal host port is returned", func() {
			schedulerCron := config.GetQuotesDiscoveryHost()
			Convey("Then the cron scraping is set to every 10 minutes", func() {
				So(schedulerCron, ShouldNotBeEmpty)
			})
		})
	})

	Convey("Given local config file is loaded", t, func() {
		config := NewConfig(EnvironmentLocal)
		Convey("Then temporal host port is returned", func() {
			schedulerCron := config.GetQuotesDiscoveryHost()
			Convey("Then the cron scraping is set to every 10 minutes", func() {
				So(schedulerCron, ShouldNotBeEmpty)
			})
		})
	})

	Convey("Given prod config file is loaded", t, func() {
		config := NewConfig(EnvironmentProd)
		Convey("Then temporal host port is returned", func() {
			schedulerCron := config.GetQuotesDiscoveryHost()
			Convey("Then the cron scraping is set to every 10 minutes", func() {
				So(schedulerCron, ShouldNotBeEmpty)
			})
		})
	})

	Convey("Given config is loaded", t, func() {
		config := NewConfig(EnvironmentLocal)
		Convey("When config GetAllEnv is called", func() {
			allEnv := config.GetAllEnv()
			Convey("Then all env is returned", func() {
				So(allEnv, ShouldContainKey, "QUOTES_DISCOVERY_HOST")
			})
		})
	})

}
