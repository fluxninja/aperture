package utils

import (
	"context"
	"net/url"

	"go.opentelemetry.io/otel/baggage"
)

// LabelsFromCtx extracts baggage labels from context.
func LabelsFromCtx(ctx context.Context) map[string]string {
	labels := make(map[string]string)
	b := baggage.FromContext(ctx)
	for _, member := range b.Members() {
		value, err := url.QueryUnescape(member.Value())
		if err != nil {
			continue
		}
		labels[member.Key()] = value
	}
	return labels
}
