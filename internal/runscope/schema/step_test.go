package schema

import (
	"encoding/json"
	"testing"
)

func TestStep_CreateResponse(t *testing.T) {
	var resp StepCreateResponse

	err := json.Unmarshal([]byte(runscopeStepCreateResponse), &resp)
	if err != nil {
		t.Error(err)
	}
}

const runscopeStepCreateResponse = `
{
  "meta": {
    "status": "success"
  },
  "data": [
    {
      "id": "d7363d46-2c07-42db-bd2e-54b37e0094cc",
      "step_type": "request",
      "skipped": false,
      "note": "Testing step, single step test",
      "method": "GET",
      "multipart_form": null,
      "headers": {
        "Accept-Encoding": [
          "application/json",
          "application/xml"
        ],
        "Authorization": [
          "Bearer 8a95deb9-3240-4a54-a045-ea186f6dc5c5"
        ]
      },
      "url": "http://example.com",
      "auth": {
        "username": "user",
        "password": "password1",
        "auth_type": "basic"
      },
      "assertions": [
        {
          "comparison": "equal_number",
          "source": "response_status",
          "value": "200"
        },
        {
          "comparison": "equal",
          "source": "response_json",
          "value": "0d813889-46cf-47ce-9bac-66e58d3372ed",
          "property": ""
        }
      ],
      "variables": [
        {
          "source": "response_status",
          "name": "httpStatus"
        },
        {
          "source": "response_header",
          "name": "httpContentEncoding",
          "property": "Content-Encoding"
        }
      ],
      "scripts": [
        "log(\"script 1\");",
        "log(\"script 2\");"
      ],
      "before_scripts": [
        "log(\"before script\");"
      ]
    }
  ],
  "error": null
}
`
