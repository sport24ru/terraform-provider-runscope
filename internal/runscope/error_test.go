package runscope

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestError_Error(t *testing.T) {
	var err Error

	json.Unmarshal([]byte(runscopeInvalidTokenResponse), &err)
	e := Error{
		Response: &http.Response{
			Status:     "FORBIDDEN",
			StatusCode: 403,
		},
		E: err.E,
	}
	expectedError := "403 You must provide a valid Authorization header to use the Runscope API."
	if e.Error() != expectedError {
		t.Errorf("Expected %s error message, got %s", expectedError, e.Error())
	}
}

const runscopeInvalidTokenResponse = `
{
  "data": {},
  "meta": {
    "status": "error"
  },
  "error": {
    "status": 403,
    "message": "You must provide a valid Authorization header to use the Runscope API.",
    "more_info": "https://www.runscope.com/docs/api/authentication"
  }
}
`
