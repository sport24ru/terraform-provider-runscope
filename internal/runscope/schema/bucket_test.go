package schema

import (
	"encoding/json"
	"testing"
)

func TestUnmarshallBucketListResponse(t *testing.T) {
	var resp BucketListResponse
	err := json.Unmarshal([]byte(runscopeBucketListOkResponse), &resp)
	if err != nil {
		t.Error(err)
	}
}

const runscopeBucketListOkResponse = `{
  "meta": {
    "status": "success"
  },
  "data": [
    {
      "name": "terraform-provider-test",
      "key": "ymdbe56klm54",
      "auth_token": null,
      "default": false,
      "verify_ssl": true,
      "team": {
        "name": "Home",
        "id": "c8ffd67b-c281-45d3-9735-3f40ee567a02"
      },
      "collections_url": "https://api.runscope.com/buckets/ymdbe56klm54/collections",
      "messages_url": "https://api.runscope.com/buckets/ymdbe56klm54/stream",
      "tests_url": "https://api.runscope.com/buckets/ymdbe56klm54/tests",
      "trigger_url": "https://api.runscope.com/radar/bucket/93a953d2-02cf-477b-a998-6562eb7873d3/trigger"
    },
    {
      "name": "terraform-provider-test",
      "key": "bx8wuwh0g8wm",
      "auth_token": null,
      "default": false,
      "verify_ssl": true,
      "team": {
        "name": "Home",
        "id": "c8ffd67b-c281-45d3-9735-3f40ee567a02"
      },
      "collections_url": "https://api.runscope.com/buckets/bx8wuwh0g8wm/collections",
      "messages_url": "https://api.runscope.com/buckets/bx8wuwh0g8wm/stream",
      "tests_url": "https://api.runscope.com/buckets/bx8wuwh0g8wm/tests",
      "trigger_url": "https://api.runscope.com/radar/bucket/521a0b5c-1c61-4d77-bd9b-721973151b68/trigger"
    },
    {
      "name": "terraform-provider-test",
      "key": "frckbsi057xt",
      "auth_token": null,
      "default": false,
      "verify_ssl": true,
      "team": {
        "name": "Home",
        "id": "c8ffd67b-c281-45d3-9735-3f40ee567a02"
      },
      "collections_url": "https://api.runscope.com/buckets/frckbsi057xt/collections",
      "messages_url": "https://api.runscope.com/buckets/frckbsi057xt/stream",
      "tests_url": "https://api.runscope.com/buckets/frckbsi057xt/tests",
      "trigger_url": "https://api.runscope.com/radar/bucket/7ca4013e-5edd-40a0-a037-2a25c09b865c/trigger"
    },
    {
      "name": "terraform-provider-test",
      "key": "b897o6h4arz6",
      "auth_token": null,
      "default": false,
      "verify_ssl": true,
      "team": {
        "name": "Home",
        "id": "c8ffd67b-c281-45d3-9735-3f40ee567a02"
      },
      "collections_url": "https://api.runscope.com/buckets/b897o6h4arz6/collections",
      "messages_url": "https://api.runscope.com/buckets/b897o6h4arz6/stream",
      "tests_url": "https://api.runscope.com/buckets/b897o6h4arz6/tests",
      "trigger_url": "https://api.runscope.com/radar/bucket/f785fc5b-8103-436a-97ba-8da41b80b0b7/trigger"
    }
  ],
  "error": null
}
`
