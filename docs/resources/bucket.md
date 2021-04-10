# Resource `runscope_bucket`

A [bucket](https://www.runscope.com/docs/api/buckets) resource.
[Buckets](https://www.runscope.com/docs/buckets) are a simple way to organize your requests and tests.

## Example Usage

```hcl
# Add a bucket to your runscope account
resource "runscope_bucket" "main" {
    name      = "a-bucket"
    team_uuid = "870ed937-bc6e-4d8b-a9a5-d7f9f2412fa3"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (String, Required) The name of this bucket.
* `team_uuid` - (String, Required) Unique identifier for the team this bucket is being created for.

## Attributes Reference

The following attributes are exported:

* `name` - The name of this bucket.
* `id` - The ID of this bucket.
* `team_uuid` - Unique identifier for the team this bucket belongs to.
* `auth_token` - Bucket auth token if set.
* `default` - `true` if this bucket is the 'default' for a team.
* `verify_ssl` - `true` if this bucket is configured to verify ssl for requests made to it.
* `trigger_url` - URL to trigger a test run for all tests within a bucket.

## Import

Buckets can be imported using the bucket `key`, e.g.

```
$ terraform import runscope_bucket.example t2f4bkvnggcx
```
