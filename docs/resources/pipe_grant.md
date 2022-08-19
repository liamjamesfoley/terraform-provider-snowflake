---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "snowflake_pipe_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# snowflake_pipe_grant (Resource)



## Example Usage

```terraform
resource snowflake_pipe_grant grant {
  database_name = "db"
  schema_name   = "schema"
  pipe_name     = "pipe"

  privilege = "operate"
  roles = [
    "role1",
    "role2",
  ]

  on_future         = false
  with_grant_option = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `database_name` (String) The name of the database containing the current or future pipes on which to grant privileges.
- `schema_name` (String) The name of the schema containing the current or future pipes on which to grant privileges.

### Optional

- `enable_multiple_grants` (Boolean) When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.
- `on_future` (Boolean) When this is set to true and a schema_name is provided, apply this grant on all future pipes in the given schema. When this is true and no schema_name is provided apply this grant on all future pipes in the given database. The pipe_name field must be unset in order to use on_future.
- `pipe_name` (String) The name of the pipe on which to grant privileges immediately (only valid if on_future is false).
- `privilege` (String) The privilege to grant on the current or future pipe.
- `roles` (Set of String) Grants privilege to these roles.
- `with_grant_option` (Boolean) When this is set to true, allows the recipient role to grant the privileges to other roles.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# format is database name | schema name | pipe name | privilege | true/false for with_grant_option
terraform import snowflake_pipe_grant.example 'dbName|schemaName|pipeName|OPERATE|false'
```