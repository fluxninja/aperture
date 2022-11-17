package log_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"

	"github.com/fluxninja/aperture/pkg/log"
)

var _ = Describe("Autosample", func() {
	It("ratelimits repeated logs", func() {
		var count Count
		for i := 0; i < 5; i++ {
			log.Autosample().Hook(&count).Info().Msg("a")
		}
		Expect(int(count)).To(Equal(1))
	})

	It("doesn't ratelimit unrelated logs", func() {
		var count Count
		log.Autosample().Hook(&count).Info().Msg("a")
		log.Autosample().Hook(&count).Info().Msg("a")
		log.Autosample().Hook(&count).Info().Msg("a")
		log.Autosample().Hook(&count).Info().Msg("bb")
		Expect(int(count)).To(Equal(4))
	})
})

var _ = Describe("logger.Autosample", func() {
	var logger *log.Logger
	BeforeEach(func() {
		logger = log.GetGlobalLogger()
	})

	It("ratelimits repeated logs", func() {
		var count Count
		for i := 0; i < 5; i++ {
			logger.Autosample().Hook(&count).Info().Msg("a")
		}
		Expect(int(count)).To(Equal(1))
	})

	It("doesn't ratelimit unrelated logs", func() {
		var count Count
		logger.Autosample().Hook(&count).Info().Msg("a")
		logger.Autosample().Hook(&count).Info().Msg("a")
		logger.Autosample().Hook(&count).Info().Msg("a")
		logger.Autosample().Hook(&count).Info().Msg("bb")
		Expect(int(count)).To(Equal(4))
	})
})

var _ = Describe("logger.Bug", func() {
	var logger *log.Logger
	BeforeEach(func() {
		logger = log.GetGlobalLogger()
	})

	It("ratelimits repeated logs", func() {
		var count Count
		logger = logger.Hook(&count)
		for i := 0; i < 5; i++ {
			logger.Bug().Msg("a")
		}
		Expect(int(count)).To(Equal(1))
	})

	It("doesn't ratelimit unrelated logs", func() {
		var count Count
		logger = logger.Hook(&count)
		logger.Bug().Msg("a")
		logger.Bug().Msg("a")
		logger.Bug().Msg("a")
		logger.Bug().Msg("bb")
		Expect(int(count)).To(Equal(4))
	})

	It("uses warn level", func() {
		var lvl SaveLevel
		logger = logger.Hook(&lvl)
		logger.Bug().Msg("a")
		Expect(lvl.last).To(Equal(zerolog.WarnLevel))
	})
})

type SaveLevel struct{ last zerolog.Level }

func (s *SaveLevel) Run(e *zerolog.Event, lvl zerolog.Level, msg string) { s.last = lvl }

type Count int

func (c *Count) Run(e *zerolog.Event, lvl zerolog.Level, msg string) { *c++ }
