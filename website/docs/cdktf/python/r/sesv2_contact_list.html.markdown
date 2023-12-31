---
subcategory: "SESv2 (Simple Email V2)"
layout: "aws"
page_title: "AWS: aws_sesv2_contact_list"
description: |-
  Terraform resource for managing an AWS SESv2 (Simple Email V2) Contact List.
---


<!-- Please do not edit this file, it is generated. -->
# Resource: aws_sesv2_contact_list

Terraform resource for managing an AWS SESv2 (Simple Email V2) Contact List.

## Example Usage

### Basic Usage

```python
# DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.aws.sesv2_contact_list import Sesv2ContactList
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        Sesv2ContactList(self, "example",
            contact_list_name="example"
        )
```

### Extended Usage

```python
# DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.aws.sesv2_contact_list import Sesv2ContactList
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        Sesv2ContactList(self, "example",
            contact_list_name="example",
            description="description",
            topic=[Sesv2ContactListTopic(
                default_subscription_status="OPT_IN",
                description="topic description",
                display_name="Example Topic",
                topic_name="example-topic"
            )
            ]
        )
```

## Argument Reference

The following arguments are required:

* `contact_list_name` - (Required) The name of the contact list.

The following arguments are optional:

* `description` - (Optional) A description of what the contact list is about.
* `tags` - (Optional) Key-value map of resource tags for the contact list. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.
* `topic` - (Optional) Configuration block(s) with topic for the contact list. Detailed below.

### topic

The following arguments are required:

* `default_subscription_status` - (Required) The default subscription status to be applied to a contact if the contact has not noted their preference for subscribing to a topic.
* `display_name` - (Required) The name of the topic the contact will see.
* `topic_name` - (Required) The name of the topic.

The following arguments are optional:

* `description` - (Optional) A description of what the topic is about, which the contact will see.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `created_timestamp` - A timestamp noting when the contact list was created in ISO 8601 format.
* `last_updated_timestamp` - A timestamp noting the last time the contact list was updated in ISO 8601 format.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import SESv2 (Simple Email V2) Contact List using the `example_id_arg`. For example:

```python
# DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import TerraformStack
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
```

Using `terraform import`, import SESv2 (Simple Email V2) Contact List using the `example_id_arg`. For example:

```console
% terraform import aws_sesv2_contact_list.example example
```

<!-- cache-key: cdktf-0.18.0 input-bf6ee9d89685af78292cb0d5392a6357f5844d7e092a1f578086fbf6539e59a5 -->