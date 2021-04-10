package schema

import (
	"encoding/json"
	"testing"
)

func TestUnmarshallTestCreateResponse(t *testing.T) {
	var resp TestCreateResponse
	err := json.Unmarshal([]byte(runscopeTestCreateOkResponse), &resp)
	if err != nil {
		t.Error(err)
	}

	expectedId := "626a024c-f75e-4f57-82d4-104fe443c0f3"
	if resp.Test.Id != expectedId {
		t.Errorf("expected Id '%s', got '%s'", expectedId, resp.Test.Id)
	}

	expectedName := "Sample Name"
	if resp.Test.Name != expectedName {
		t.Errorf("expected Name '%s', got '%s'", expectedName, resp.Test.Name)
	}

	expectedDescription := ""
	if resp.Test.Description != expectedDescription {
		t.Errorf("expected Description '%s', got '%s'", expectedDescription, resp.Test.Description)
	}

	expectedDefaultEnvironmentId := "a50b63cc-c377-4823-9a95-8b91f12326f2"
	if resp.Test.DefaultEnvironmentId != expectedDefaultEnvironmentId {
		t.Errorf("expected DefaultEnvironmentId '%s', got '%s'", expectedDefaultEnvironmentId, resp.Test.DefaultEnvironmentId)
	}

	var expectedCreatedAt int64 = 1438832081
	if resp.CreatedAt != expectedCreatedAt {
		t.Errorf("expected CreatedAt '%d', got '%d'", expectedCreatedAt, resp.CreatedAt)
	}

	expectedCreatedBy := CreatedBy{
		Id:    "4ee15ecc-7fe1-43cb-aa12-ef50420f2cf9",
		Name:  "Grace Hopper",
		Email: "grace@example.com",
	}
	if resp.CreatedBy != expectedCreatedBy {
		t.Errorf("expected CreatedBy '%+v', got '%+v'", expectedCreatedBy, resp.CreatedBy)
	}
}

const runscopeTestCreateOkResponse = `{
    "data": {
        "created_at": 1438832081,
        "created_by": {
            "email": "grace@example.com",
            "name": "Grace Hopper",
            "id": "4ee15ecc-7fe1-43cb-aa12-ef50420f2cf9"
        },
        "default_environment_id": "a50b63cc-c377-4823-9a95-8b91f12326f2",
        "description": null,
        "environments": [
            {
                "emails": {
                    "notify_all": false,
                    "notify_on": "all",
                    "notify_threshold": 1,
                    "recipients": []
                },
                "initial_variables": {
                    "base_url": "https://api.example.com"
                },
                "integrations": [
                    {
                        "description": "Pagerduty Account",
                        "integration_type": "pagerduty",
                        "id": "53776d9a-4f34-4f1f-9gff-c155dfb6692e"
                    }
                ],
                "name": "Test Settings",
                "parent_environment_id": null,
                "preserve_cookies": false,
                "regions": [
                    "us1"
                ],
                "remote_agents": [],
                "script": "",
                "test_id": "626a024c-f75e-4f57-82d4-104fe443c0f3",
                "id": "a50b63cc-c377-4823-9a95-8b91f12326f2",
                "verify_ssl": true,
                "webhooks": null
            }
        ],
        "last_run": null,
        "name": "Sample Name",
        "schedules": [],
        "steps": [
            {
                "assertions": [
                    {
                        "comparison": "is_equal",
                        "source": "response_status",
                        "value": 200
                    }
                ],
                "auth": {},
                "body": "",
                "form": {},
                "headers": {},
                "method": "GET",
                "note": "",
                "step_type": "request",
                "url": "https://yourapihere.com/",
                "id": "53f8e1fd-0989-491a-9f15-cc055f27d097",
                "variables": []
            }
        ],
        "trigger_url": "http://api.runscope.com/radar/b96ecee2-cce6-4d80-8f07-33ac22a22ebd/trigger",
        "id": "626a024c-f75e-4f57-82d4-104fe443c0f3"
    },
    "error": null,
    "meta": {
        "status": "success"
    }
}`
