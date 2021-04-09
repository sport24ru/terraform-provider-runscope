# Resource `runscope_test`

A [test](https://www.runscope.com/docs/api/tests) resource.
[Tests](https://www.runscope.com/docs/buckets) are made up of
a collection of [test steps](step.html) and an
[environment](environment.html).

## Example Usage

```hcl
resource "runscope_bucket" "main" {
  name      = "terraform-ftw"
  team_uuid = "870ed937-bc6e-4d8b-a9a5-d7f9f2412fa3"
}

resource "runscope_test" "api" {
    name        = "api-test"
    description = "checks the api is up and running"
    bucket_id   = runscope_bucket.main.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of this test.
* `description` - (Optional) Human-readable description of the new test.
  is being created for.

## Attribute Reference

The following attributes are exported:

* `id` - The unique identifier for the test.
* `name` - The name of this test.
* `description` - Human-readable description of the new test.

## Import

Test can be imported using the bucket ID and the test UUID, e.g.

```
$ terraform import runscope_test.example t2f4bkvnggcx/ea37dff1-36e1-44ae-aa7e-48693f235660
```
