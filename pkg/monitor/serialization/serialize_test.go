package monitorserialization

import (
	"testing"
	"time"

	"github.com/openshift/origin/pkg/monitor/monitorapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntervalsJSONPreservesKubeEventUIDAnnotation(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	original := monitorapi.Intervals{
		monitorapi.NewInterval(monitorapi.SourceKubeEvent, monitorapi.Info).
			Locator(monitorapi.Locator{
				Type: monitorapi.LocatorTypeKind,
				Keys: map[monitorapi.LocatorKey]string{
					monitorapi.LocatorNamespaceKey:  "openshift-authentication",
					monitorapi.LocatorDeploymentKey: "oauth-openshift",
					monitorapi.LocatorHmsgKey:       "0123456789",
				},
			}).
			Message(monitorapi.NewMessage().
				Reason("ScalingReplicaSet").
				HumanMessage("sample message").
				WithAnnotation(monitorapi.AnnotationKubeEventUID, "uid-a")).
			Build(now, now.Add(time.Second)),
	}

	serialized, err := IntervalsToJSON(original)
	require.NoError(t, err)
	assert.Contains(t, string(serialized), `"event-uid": "uid-a"`)

	roundTripped, err := IntervalsFromJSON(serialized)
	require.NoError(t, err)
	require.Len(t, roundTripped, 1)
	assert.Equal(t, "uid-a", roundTripped[0].Message.Annotations[monitorapi.AnnotationKubeEventUID])
	assert.Equal(t, original[0].Locator, roundTripped[0].Locator, "event UID metadata must not change locator/hmsg identity")
}
