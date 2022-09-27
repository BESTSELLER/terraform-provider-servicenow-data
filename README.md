# terraform-provider-servicenow-data

Terraform provider that can manage data in any SN table it can Read/Write
## Limitations

This provider only lets you manage individual rows in Tables, you cannot manage the tables themselves.

## Set-up

The provider requires you to create and configure the tables and Service-Now Beforehand.
You will also need an account with API Access.,

The provider can also be configured using Environmental variables:

```
SN_API_URL
SN_API_USER
SN_API_PASS
```

See examples for how to use it.

