package prometheus

import (
	"aw-sync-agent/activitywatch"
	"aw-sync-agent/util"
	"fmt"
	"regexp"
	"strconv"
)

func AttachTimeSeriesPayload(event activitywatch.Event, includeHostName bool, watcher string, userID string) TimeSeries {
	var labels []Label

	AddMetricLabel(&labels, "__name__", SanitizeLabelName(watcher))
	AddMetricLabel(&labels, "unique_id", fmt.Sprintf("aw-%d", event.ID)) // déterministe, plus de doublons visibles
	AddMetricLabel(&labels, "aw_id", strconv.Itoa(event.ID))
	AddMetricLabel(&labels, "user", userID)

	hostValue := "Unknown"
	if includeHostName {
		hostValue = util.GetHostname()
	}
	AddMetricLabel(&labels, "host", hostValue)

	if event.Data["category"] == nil {
		AddMetricLabel(&labels, "category", "Other")
	}
	for key, value := range event.Data {
		AddMetricLabel(&labels, key, fmt.Sprintf("%v", value))
	}

	sample := Sample{
		Value: event.Duration,
		Time:  event.Timestamp,
	}

	return TimeSeries{
		Labels: labels,
		Sample: sample,
		ID:     event.ID, // NOUVEAU
	}
}

func AddMetricLabel(labels *[]Label, key string, value string) {
	if value != "" {
		*labels = append(*labels, Label{
			Name:  key,
			Value: value,
		})
	}
}

// SanitizeLabelName ensures the label name conforms to Prometheus naming conventions
func SanitizeLabelName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	return re.ReplaceAllString(name, "_")
}
