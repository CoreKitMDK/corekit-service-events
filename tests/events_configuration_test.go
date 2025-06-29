package tests

import (
	"github.com/CoreKitMDK/corekit-service-events/v2/pkg/events"
	"testing"
	"time"
)

func TestMetricsConfiguration(t *testing.T) {

	config := events.NewConfiguration()

	config.UseConsole = true

	config.UseNATS = true
	config.NatsURL = "internal-events-broker-nats-client"

	config.NatsPassword = "internal-events-broker"
	config.NatsUsername = "internal-events-broker"

	ogger := config.Init()
	defer ogger.Stop()

	_ = ogger.Emit("REGISTER_USER", "{\"uuid\":\"1234\"}")
	_ = ogger.Emit("UPDATE_USER_TOKEN", "{\"uuid\":\"1234\"}")
	_ = ogger.Emit("DELETE_USER", "{\"uuid\":\"1234\"}")

	time.Sleep(2 * time.Second)
}
