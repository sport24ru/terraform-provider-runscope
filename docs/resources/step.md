# Resource `runscope_step`

A [step](https://www.runscope.com/docs/api/steps) resource.
API tests are comprised of a series of steps, most often HTTP requests.
In addition to requests, you can also add additional types of steps to
your tests like pauses and conditions.

## Example Usage

```hcl
resource "runscope_bucket" "bucket" {
    name      = "terraform-provider-test"
    team_uuid = "dfb75aac-eeb3-4451-8675-3a37ab421e4f"
}

resource "runscope_test" "test" {
  bucket_id   = runscope_bucket.bucket.id
  name        = "runscope test"
  description = "This is a test test..."
}

resource "runscope_step" "main_page" {
    bucket_id = runscope_bucket.bucket.id
    test_id   = runscope_test.test.id
    step_type = "request"
    url       = "http://example.com"
    note      = "A comment for the test step"
    method    = "GET"
    variable {
        name   = "httpStatus"
        source = "response_status"
    }
    variable {
        name     = "httpContentEncoding"
        source   = "response_header"
        property = "Content-Encoding"
    }

    assertion {
        source     = "response_status"
        comparison = "equal_number"
        value      = "200"
    }
    assertion {
        source     = "response_json"
        comparison = "equal"
        value      = "c5baeb4a-2379-478a-9cda-1b671de77cf9"
        property   = "data.id"
    }

    auth {
        username  = "myUsername"
        auth_type = "basic"
        password  = "myPassword"
    }
    before_scripts = [<<EOF
       var endVar = new Date();
       var startVar = new Date();
       alert('this is a multi-line before script')
    EOF
    ]
    scripts = [<<EOF
       var endVar = new Date();
       var startVar = new Date();
       alert('this is a multi-line after script')
    EOF
    ]
    header {
        header = "Accept-Encoding"
        value  = "application/json"
    }
    header {
        header = "Accept-Encoding"
        value  = "application/xml"
    }
    header {
        header = "Authorization"
        value  = "Bearer bb74fe7b-b9f2-48bd-9445-bdc60e1edc6a"
    }

}
```

## Argument Reference

The following arguments are supported:

* `bucket_id` - (Required) The id of the bucket to associate this step with.
* `test_id` - (Required) The id of the test to associate this step with.
* `note` = (Optional) A comment attached to the test step.
* `step_type` - (Required) The type of step.
  * [request](#request-steps)
  * pause
  * condition
  * ghost
  * subtest

### Request steps

When creating a `request` type of step the additional arguments also apply:

* `method` - (Required) The HTTP method for this request step.
* `variable` - (Optional) Block describing variable to extract out of the HTTP response from this request. May be declared multiple times. Variable documented below.
* `assertion` - (Optional) Block describing assertion to apply to the HTTP response from this request. May be declared multiple times. Assertion documented below.
* `header` - (Optional) Block describing header to apply to the request. May be declared multiple times. Header documented below.
* `body` - (Optional) A string to use as the body of the request.
* `auth` - (Optional) The credentials used to authenticate the request
* `before_script` - (Optional) Runs a script before the request is made
* `script` - (Optional) Runs a script after the request is made

Variable (`variable`) supports the following:

* `name` - (Required) Name of the variable to define.
* `property` - (Required) The name of the source property. i.e. header name or json path
* `source` - (Required) The variable source, for list of allowed values see: https://api.blazemeter.com/api-monitoring/#variable-sources-list

Assertion (`assertion`) supports the following:

* `source` - (Required) The assertion source, for list of allowed values see: https://api.blazemeter.com/api-monitoring/#assertion-sources-list
* `property` - (Optional) The name of the source property. i.e. header name or json path
* `comparison` - (Required) The assertion comparison to make i.e. `equals`, for list of allowed values see: https://api.blazemeter.com/api-monitoring/#assertion-comparisons-list
* `value` - (Optional) The value the `comparison` will use

**Example Assertion**

Status Code == 200

```
"assertion": [
    {
        "source": "response_status",
        "comparison": "equal_number",
        "value": 200
    }
]
```

JSON element 'address' contains the text "avenue"

```
"assertion": [
    {
        "source": "response_json",
        "property": "address",
        "comparison": "contains",
        "value": "avenue"
    }
]
```

Response Time is faster than 1 second.

```
"assertion": [
    {
        "source": "response_time",
        "comparison": "is_less_than",
        "value": 1000
    }
]
```

The `headers` list supports the following:

* `header` - (Required) The name of the header
* `value` - (Required) The name header value

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the step.

## Import

Test can be imported using the bucket ID, test ID and step ID e.g.

```
$ terraform import runscope_test.example t2f4bkvnggcx/ea37dff1-36e1-44ae-aa7e-48693f235660/ad41345e-8874-4bcd-904a-d85901239789
```

or you may use position of step in test (starting from 1) e.g.

```
$ terraform import runscope_test.example t2f4bkvnggcx/ea37dff1-36e1-44ae-aa7e-48693f235660#1
```
