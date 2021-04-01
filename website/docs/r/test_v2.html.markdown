---
layout: "runscope"
page_title: "Runscope: runscope_test_v2"
sidebar_current: "docs-runscope-resource-test-v2"
description: |-
  Provides a Runscope test resource with configurable default environment.
---

# runscope\_test\_v2

A [test](https://www.runscope.com/docs/api/tests) resource.
[Tests](https://www.runscope.com/docs/buckets) are made up of
a collection of [test steps](step.html) and an
[environment](environment.html).

The difference with `runscope_test` is the ability to configure
default environment.

## Example Usage

```hcl
resource "runscope_bucket" "main" {
    name      = "terraform-ftw"
    team_uuid = "870ed937-bc6e-4d8b-a9a5-d7f9f2412fa3"
}

resource "runscope_environment" "prod" {
    bucket_id = runscope_bucket.prod.id
    name      = "main"
}

resource "runscope_test_v2" "api" {
  name                   = "api-test"
  description            = "checks the api is up and running"
  bucket_id              = runscope_bucket.main.id
  default_environment_id = runscope_environment.prod.id
}
```

## Argument Reference

The following arguments are supported:

* `bucket_id` - (Required) ID of bucket containing test.
* `name` - (Required) The name of this test.
* `description` - (Optional) Human-readable description of the new test.
  is being created for.
* `default_environment_id` - (Required) ID of default environment of test. If you don't need to configure it, just use `runscope_test`.

## Import

Test can be imported using the bucket ID and the test UUID, e.g.

```
$ terraform import runscope_test_v2.example t2f4bkvnggcx/ea37dff1-36e1-44ae-aa7e-48693f235660
```
