---
page_title: "samsungcloudplatformv2 Provider"
subcategory: ""
description: |-
    Interact with samsungcloudplatformv2.
---

# samsungcloudplatformv2 Provider

The samsungcloudplatformv2 provider is used to interact with Samsung Cloud Platform services.

The provider needs to be configured with the proper credentials before it can be used.

## Example Usage

```terraform
terraform {
  required_providers {
    samsungcloudplatformv2 = {
      version = "1.0.3"
      source = "SamsungSDSCloud/samsungcloudplatformv2"
    }
  }
  required_version = ">= 1.11"
}

provider "samsungcloudplatformv2" {
}

//Create a new virtual server instance
resource "samsungcloudplatformv2_virtualserver_server" "server" {
  name           = var.name
  state          = var.state
  image_id       = var.image_id
  server_type_id = var.server_type_id
  #...
}
```

## Setup credentials


### Create local setting file
Create .scpconf directory in your OS account home


```
cd %USERPROFILE%
mkdir ".scpconf"
```

Create config.json and credentials.json in .scpconf directory

### Add Samsungcloudplatform configuration
Insert following parameters in ```.scpconf/config.json``` file
```
{
    "auth-url": "https://iam.dev2.samsungsdscloud.com/v1",
    "default-region": "kr-west1"
}
```

### auth-url example

| **Environment**     | **env value** | **Example Service URL**              |
|---------------|---------------|--------------------------------------|
| for Samsung   | s             | https://iam.s.samsungsdscloud.com/v1 |
| for Sovereign | g             | https://iam.g.samsungsdscloud.com/v1 |


### Add your credentials
Insert following parameters in ```.scpconf/credentials.json``` file
```
{
    "access-key": "xxxxxxxxxxxxxx",
    "secret-key": "xxxxxxxxxxxxxx"
}
```

## Schema

### Optional

- `access_key` (String) Access key for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_ACCESS_KEY environment variable.
- `auth_url` (String) Authentication URL for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_AUTH_URL environment variable.
- `default_region` (String) Default region configuration for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_DEFAULT_REGION environment variable.
- `secret_key` (String) Secret key for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_SECRET_KEY environment variable.
