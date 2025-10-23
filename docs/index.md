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
      version = "2.0.3"
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

- `auth-url` (String) Authentication URL for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_AUTH_URL environment variable.
- `default-region` (String) Default region configuration for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_DEFAULT_REGION environment variable.

### Optional
- `microversion-check-timeout` (Int) timeout of an api call that inquires about the endpoint status of a product that provides terraform. The value can be adjusted according to the Internet environment. Default 15 seconds
- `max-remain-days`  (Int) Reference value to compare from the user's point of view the date of the product's SDK that is no longer guaranteed. The user can flexibly adjust the check for the sdk's deprecated date through this value. Default 90 days

### auth-url example

| **Environment**     | **env value** | **Example Service URL**              |
|---------------|---------------|--------------------------------------|
| for Samsung   | s             | https://iam.s.samsungsdscloud.com/v1 |
| for Sovereign | g             | https://iam.g.samsungsdscloud.com/v1 |
| for Enterprise | e             | https://iam.e.samsungsdscloud.com/v1 |


### Add your credentials
Insert following parameters in ```.scpconf/credentials.json``` file
```
{
    "access-key": "xxxxxxxxxxxxxx",
    "secret-key": "xxxxxxxxxxxxxx"
}
```
- `access_key` (String) Access key for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_ACCESS_KEY environment variable.
- `secret_key` (String) Secret key for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_SECRET_KEY environment variable.

## Schema
### Optional
- `access_key` (String) Access key for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_ACCESS_KEY environment variable.
- `auth_url` (String) Authentication URL for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_AUTH_URL environment variable.
- `default_region` (String) Default region configuration for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_DEFAULT_REGION environment variable.
- `secret_key` (String) Secret key for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_SECRET_KEY environment variable.
- `max_remain_days` (Int) Max Remain Days configuration for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_MAX_REMAIN_DAYS environment variable.
- `microversion_check_timeout` (Int) Microversion Check Timeout for calling samsungcloudplatformv2 API. May also be provided via SCP_TF_MICROVERSION_CHECK_TIMEOUT environment variable.
