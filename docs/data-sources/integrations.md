# Data Source `runscope_integration`

Use this data source to list all of your [integrations](https://www.runscope.com/docs/api/integrations)
that you can use with other runscope resources.

## Example Usage

```hcl
data "runscope_integrations" "slack" {
  team_uuid = "d26553c0-3537-40a8-9d3c-64b0453262a9"
  filter = {
    name   = "type"
    values = ["slack"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Filter to reduce the list of integrations returned.

Variables (`filter`) supports the following:

* `name` - The name of the field to filter on, currently either: `id`, `type` or `description`.
* `values` - The list of values to match against

## Attributes Reference
The following attributes are exported:

* `ids` - Set of identifiers of the found integrations.
