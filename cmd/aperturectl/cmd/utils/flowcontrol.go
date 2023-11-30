package utils

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	flowcontrolv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/check/v1"
	previewv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/preview/v1"
)

// ParseControlPoints parses the control points.
func ParseControlPoints(client IntrospectionClient) error {
	resp, err := client.ListFlowControlPoints(
		context.Background(),
		&cmdv1.ListFlowControlPointsRequest{},
	)
	if err != nil {
		return err
	}

	if resp.ErrorsCount != 0 {
		fmt.Fprintf(os.Stderr, "Could not get answer from %d agents", resp.ErrorsCount)
	}

	slices.SortFunc(resp.GlobalFlowControlPoints, func(a, b *cmdv1.GlobalFlowControlPoint) int {
		if a.AgentGroup != b.AgentGroup {
			return strings.Compare(a.AgentGroup, b.AgentGroup)
		}
		if a.FlowControlPoint.Service != b.FlowControlPoint.Service {
			return strings.Compare(a.FlowControlPoint.Service, b.FlowControlPoint.Service)
		}
		return strings.Compare(a.FlowControlPoint.ControlPoint, b.FlowControlPoint.ControlPoint)
	})

	tabwriter := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', 0)
	fmt.Fprintln(tabwriter, "AGENT GROUP\tSERVICE\tNAME\tTYPE")
	for _, cp := range resp.GlobalFlowControlPoints {
		fmt.Fprintf(tabwriter, "%s\t%s\t%s\t%s\n",
			cp.AgentGroup,
			cp.FlowControlPoint.Service,
			cp.FlowControlPoint.ControlPoint,
			cp.FlowControlPoint.Type)
	}
	tabwriter.Flush()

	return nil
}

type PreviewInput struct {
	AgentGroup    string
	Service       string
	ControlPoint  string
	NumSamples    int
	IsHTTPPreview bool
}

// ParsePreview parses the preview.
func ParsePreview(client IntrospectionClient, input PreviewInput) error {
	previewReq := &previewv1.PreviewRequest{
		Samples:      int64(input.NumSamples),
		Service:      input.Service,
		ControlPoint: input.ControlPoint,
		// FIXME LabelMatcher: Figure out how to represent label matcher in CLI.
	}

	if input.IsHTTPPreview {
		resp, err := client.PreviewHTTPRequests(
			context.Background(),
			&cmdv1.PreviewHTTPRequestsRequest{
				AgentGroup: input.AgentGroup,
				Request:    previewReq,
			},
		)
		if err != nil {
			return err
		}
		samplesJSON, err := protojson.MarshalOptions{Multiline: true}.Marshal(resp.Response)
		if err != nil {
			return err
		}
		os.Stdout.Write(samplesJSON)
	} else {
		resp, err := client.PreviewFlowLabels(
			context.Background(),
			&cmdv1.PreviewFlowLabelsRequest{
				AgentGroup: input.AgentGroup,
				Request:    previewReq,
			},
		)
		if err != nil {
			return err
		}
		samplesJSON, err := protojson.MarshalOptions{Multiline: true}.Marshal(resp.Response)
		if err != nil {
			return err
		}
		os.Stdout.Write(samplesJSON)
	}

	return nil
}

type CacheLookupInput struct {
	AgentGroup   string
	ControlPoint string
	Key          string
}

func ParseResultCacheLookup(client IntrospectionClient, input CacheLookupInput) error {
	resp, err := client.CacheLookup(
		context.Background(),
		&cmdv1.GlobalCacheLookupRequest{
			AgentGroup: input.AgentGroup,
			Request: &flowcontrolv1.CacheLookupRequest{
				ControlPoint:   input.ControlPoint,
				ResultCacheKey: input.Key,
			},
		},
	)
	if err != nil {
		return err
	}

	if resp.ResultCacheResponse == nil {
		fmt.Fprintf(os.Stderr, "Could not get answer")
		return nil
	}
	if resp.ResultCacheResponse.Error != "" {
		fmt.Fprintf(os.Stderr, "Error: %s", resp.ResultCacheResponse.Error)
		return nil
	}
	if resp.ResultCacheResponse.LookupStatus == flowcontrolv1.CacheLookupStatus_MISS {
		fmt.Fprintf(os.Stderr, "Cache miss")
		return nil
	}

	val := string(resp.ResultCacheResponse.Value)
	fmt.Fprintf(os.Stdout, "%s\n", val)

	return nil
}

func ParseGlobalCacheLookup(client IntrospectionClient, input CacheLookupInput) error {
	resp, err := client.CacheLookup(
		context.Background(),
		&cmdv1.GlobalCacheLookupRequest{
			AgentGroup: input.AgentGroup,
			Request: &flowcontrolv1.CacheLookupRequest{
				GlobalCacheKeys: []string{input.Key},
			},
		},
	)
	if err != nil {
		return err
	}

	if resp.GlobalCacheResponses == nil || resp.GlobalCacheResponses[input.Key] == nil {
		fmt.Fprintf(os.Stderr, "Could not get answer")
		return nil
	}
	lookupResponse := resp.GlobalCacheResponses[input.Key]
	if lookupResponse.Error != "" {
		fmt.Fprintf(os.Stderr, "Error: %s", lookupResponse.Error)
		return nil
	}
	if lookupResponse.LookupStatus == flowcontrolv1.CacheLookupStatus_MISS {
		fmt.Fprintf(os.Stderr, "Cache miss")
		return nil
	}

	val := string(lookupResponse.Value)
	fmt.Fprintf(os.Stdout, "%s\n", val)

	return nil
}

type CacheUpsertInput struct {
	AgentGroup   string
	ControlPoint string
	Key          string
	Value        string
	TTL          time.Duration
}

func ParseResultCacheUpsert(client IntrospectionClient, input CacheUpsertInput) error {
	resp, err := client.CacheUpsert(
		context.Background(),
		&cmdv1.GlobalCacheUpsertRequest{
			AgentGroup: input.AgentGroup,
			Request: &flowcontrolv1.CacheUpsertRequest{
				ControlPoint: input.ControlPoint,
				ResultCacheEntry: &flowcontrolv1.CacheEntry{
					Key:   input.Key,
					Value: []byte(input.Value),
					Ttl:   durationpb.New(input.TTL),
				},
			},
		},
	)
	if err != nil {
		return err
	}

	samplesJSON, err := protojson.MarshalOptions{Multiline: true}.Marshal(resp)
	if err != nil {
		return err
	}
	os.Stdout.Write(samplesJSON)

	return nil
}

func ParseGlobalCacheUpsert(client IntrospectionClient, input CacheUpsertInput) error {
	resp, err := client.CacheUpsert(
		context.Background(),
		&cmdv1.GlobalCacheUpsertRequest{
			AgentGroup: input.AgentGroup,
			Request: &flowcontrolv1.CacheUpsertRequest{
				GlobalCacheEntries: map[string]*flowcontrolv1.CacheEntry{
					input.Key: {
						Key:   input.Key,
						Value: []byte(input.Value),
						Ttl:   durationpb.New(input.TTL),
					},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	samplesJSON, err := protojson.MarshalOptions{Multiline: true}.Marshal(resp)
	if err != nil {
		return err
	}
	os.Stdout.Write(samplesJSON)

	return nil
}

type CacheDeleteInput struct {
	AgentGroup   string
	ControlPoint string
	Key          string
}

func ParseResultCacheDelete(client IntrospectionClient, input CacheDeleteInput) error {
	resp, err := client.CacheDelete(
		context.Background(),
		&cmdv1.GlobalCacheDeleteRequest{
			AgentGroup: input.AgentGroup,
			Request: &flowcontrolv1.CacheDeleteRequest{
				ControlPoint:   input.ControlPoint,
				ResultCacheKey: input.Key,
			},
		},
	)
	if err != nil {
		return err
	}

	samplesJSON, err := protojson.MarshalOptions{Multiline: true}.Marshal(resp)
	if err != nil {
		return err
	}
	os.Stdout.Write(samplesJSON)

	return nil
}

func ParseGlobalCacheDelete(client IntrospectionClient, input CacheDeleteInput) error {
	resp, err := client.CacheDelete(
		context.Background(),
		&cmdv1.GlobalCacheDeleteRequest{
			AgentGroup: input.AgentGroup,
			Request: &flowcontrolv1.CacheDeleteRequest{
				GlobalCacheKeys: []string{input.Key},
			},
		},
	)
	if err != nil {
		return err
	}

	samplesJSON, err := protojson.MarshalOptions{Multiline: true}.Marshal(resp)
	if err != nil {
		return err
	}
	os.Stdout.Write(samplesJSON)

	return nil
}
