---
page_title: "samsungcloudplatformv2_filestorage_snapshot_schedule Resource - samsungcloudplatformv2"
subcategory: File Storage Snapshot Schedule
description: |-
  Lists of SnapshotSchedules.
---

# samsungcloudplatformv2_filestorage_snapshot_schedule (Resource)

Lists of SnapshotSchedules.

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_filestorage_snapshot_schedule" "snapshot_schedule" {
  volume_id = var.volume_id
  snapshot_retention_count = var.snapshot_retention_count
  snapshot_schedule = var.snapshot_schedule
}

output "snapshot_schedule_output" {
  value = samsungcloudplatformv2_filestorage_snapshot_schedule.snapshot_schedule
}

variable "volume_id" {
  type    = string
  default = ""
}

variable "snapshot_retention_count" {
  type    = number
  default = 0
}

variable "snapshot_schedule" {
  type = object({
    frequency   = string
    day_of_week = string
    hour        = string
  })
  default = {
    day_of_week = ""
    frequency   = ""
    hour        = ""
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `volume_id` (String) Volume ID 
  - example : 'bfdbabf2-04d9-4e8b-a205-020f8e6da438'

### Optional

- `snapshot_retention_count` (Number) Snapshot Retention Count (If not entered, 10 will be applied.) 
  - example : 1 
  - maximum : 128 
  - minimum : 1
- `snapshot_schedule` (Attributes) Snapshot Schedule (see [below for nested schema](#nestedatt--snapshot_schedule))

### Read-Only

- `snapshot_policy_enabled` (Boolean) Snapshot Policy Enabled 
  - example : 'true'

<a id="nestedatt--snapshot_schedule"></a>
### Nested Schema for `snapshot_schedule`

Optional:

- `day_of_week` (String) Day Of Week 
  - example : 'MON' 
  - pattern: '^(SUN|MON|TUE|WED|THU|FRI|SAT)$'
- `frequency` (String) Frequency 
  - example : 'DAILY' 
  - pattern: '^(WEEKLY|DAILY)$'
- `hour` (String) Hour 
  - example : '0' 
  - maximum : 23 
  - minimum : 0  
  - pattern: '^([0-9]|1[0-9]|2[0-3])$'

Read-Only:

- `id` (String) ID