package extractors

import (
	"encoding/json"

	"aperture.tech/aperture/pkg/log"
)

var inputTemplate = func() interface{} {
	var out interface{}
	if err := json.Unmarshal([]byte(inputTemplateJSON), &out); err != nil {
		log.Panic().Err(err).Msgf("Failed to unmarshal: %v", err)
	}
	return out
}()

const inputTemplateJSON = `
{
    "destination": {
        "address": {
            "socketAddress": {
                "portValue": 8000,
                "address": "10.25.95.68"
            }
        }
    },
    "metadataContext": {
        "filterMetadata": {
            "envoy.filters.http.header_to_metadata": {
                "policy_type": "ingres"
            }
        }
    },
    "request": {
        "http": {
            "headers": {
                ":authority": "example-app",
                ":method": "POST",
                ":path": "/pets/dogs",
                "accept": "*/*",
                "authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWxpY2lhIFNtaXRoc29uaWFuIiwicm9sZXMiOlsicmVhZGVyIiwid3JpdGVyIl0sInVzZXJuYW1lIjoiYWxpY2UifQ.md2KPJFH9OgBq-N0RonGdf5doGYRO_1miN8ugTSeTYc",
                "content-length": "0",
                "user-agent": "curl/7.68.0-DEV",
                "x-ext-auth-allow": "yes",
                "x-forwarded-proto": "http",
                "x-request-id": "1455bbb0-0623-4810-a2c6-df73ffd8863a"
            },
            "host": "example-app",
            "id": "8306787481883314548",
            "method": "POST",
            "path": "/pets/dogs",
            "protocol": "HTTP/1.1",
            "body": "{\"query\":\"SELECT name FROM users WHERE password = @pass\",\"arguments\":{\"pass\":\"foobar1\"}}",
            "bearer": "special attribute that pulls Bearer value from authorization header"
        }
    },
    "source": {
        "address": {
            "socketAddress": {
                "portValue": 33772,
                "address": "10.25.95.69"
            }
        }
    }
}
`
