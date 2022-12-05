package distcache

import (
	"bytes"

	"github.com/fluxninja/aperture/pkg/log"
)

// OlricLogWriter is wrapper around aperture Logger that parse the message before writing to olricConfig.LogOutput.
type OlricLogWriter struct {
	Logger *log.Logger
}

// Write writes the message and logs in aperture Logger format.
func (ol *OlricLogWriter) Write(message []byte) (int, error) {
	message = bytes.TrimSpace(message)
	if len(message) < 2 {
		ol.Logger.Debug().Msg(string(message))
	} else {
		switch message[1] {
		case 'I':
			ol.Logger.Debug().Msg(string(message[7:]))
		case 'W':
			ol.Logger.Warn().Msg(string(message[7:]))
		case 'E':
			ol.Logger.Error().Msg(string(message[8:]))
		case 'D':
			ol.Logger.Debug().Msg(string(message[8:]))
		default:
			ol.Logger.Debug().Msg(string(message))
		}
	}
	return len(message), nil
}
