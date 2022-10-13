package jobs

import (
	"context"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/fluxninja/aperture/pkg/log"
)

// swagger:operation POST /liveness common-configuration Liveness
// ---
// x-fn-config-env: true
// parameters:
// - name: service
//   in: body
//   schema:
//     "$ref": "#/definitions/JobConfig"
// - name: scheduler
//   in: body
//   schema:
//     "$ref": "#/definitions/JobGroupConfig"

// swagger:operation POST /readiness common-configuration Readiness
// ---
// x-fn-config-env: true
// parameters:
// - name: service
//   in: body
//   schema:
//     "$ref": "#/definitions/JobConfig"
// - name: scheduler
//   in: body
//   schema:
//     "$ref": "#/definitions/JobGroupConfig"

// Default jobs groups and names.
const (
	livenessGroup         = "liveness"
	livenessMultiJobName  = "service"
	readinessGroup        = "readiness"
	readinessMultiJobName = "service"
)

// Module is a fx module that provides default jobs for internal service components.
func Module() fx.Option {
	return fx.Options(
		provideServiceLivenessCheck(),
		provideServiceReadinessCheck(),
		fx.Invoke(registerSelfChecks),
	)
}

// Liveness job for internal service components (MultiJob).
func provideServiceLivenessCheck() fx.Option {
	mjc := MultiJobConstructor{
		Name:         livenessMultiJobName,
		JobGroupName: livenessGroup,
	}

	return fx.Options(
		JobGroupConstructor{Name: livenessGroup}.Annotate(),
		mjc.Annotate(),
	)
}

// Readiness job for internal service components (MultiJob).
func provideServiceReadinessCheck() fx.Option {
	mjc := MultiJobConstructor{
		Name:         readinessMultiJobName,
		JobGroupName: readinessGroup,
	}

	return fx.Options(
		JobGroupConstructor{Name: readinessGroup}.Annotate(),
		mjc.Annotate(),
	)
}

// SelfChecksIn holds parameters for RegisterSelfChecks.
type SelfChecksIn struct {
	fx.In

	Liveness  *MultiJob `name:"liveness.service"`
	Readiness *MultiJob `name:"readiness.service"`
}

// RegisterSelfChecks registers self check jobs (liveness, readiness) for internal service components.
func registerSelfChecks(sc SelfChecksIn) {
	liveness := &BasicJob{
		JobBase: JobBase{
			JobName: "job-module",
		},
		JobFunc: checkSelfLiveness,
	}
	err := sc.Liveness.RegisterJob(liveness)
	if err != nil {
		log.Error().Err(err).Msg("Unable to register liveness job")
		return
	}

	readiness := &BasicJob{
		JobBase: JobBase{
			JobName: "job-module",
		},
		JobFunc: checkSelfReadiness,
	}
	err = sc.Readiness.RegisterJob(readiness)
	if err != nil {
		log.Error().Err(err).Msg("Unable to register readiness job")
		return
	}
}

func checkSelfLiveness(context.Context) (proto.Message, error) {
	message := wrapperspb.String("Liveness job module is alive")
	return message, nil
}

func checkSelfReadiness(context.Context) (proto.Message, error) {
	message := wrapperspb.String("Readiness job module is alive")
	return message, nil
}
