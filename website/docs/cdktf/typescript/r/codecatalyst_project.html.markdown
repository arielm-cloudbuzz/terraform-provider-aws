---
subcategory: "CodeCatalyst"
layout: "aws"
page_title: "AWS: aws_codecatalyst_project"
description: |-
  Terraform resource for managing an AWS CodeCatalyst Project.
---


<!-- Please do not edit this file, it is generated. -->
# Resource: aws_codecatalyst_project

Terraform resource for managing an AWS CodeCatalyst Project.

## Example Usage

### Basic Usage

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { CodecatalystProject } from "./.gen/providers/aws/codecatalyst-project";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new CodecatalystProject(this, "test", {
      description: "My CodeCatalyst Project created using Terraform",
      displayName: "MyProject",
      spaceName: "myproject",
    });
  }
}

```

## Argument Reference

The following arguments are required:

* `spaceName` - (Required) The name of the space.
* `displayName` - (Required) The friendly name of the project that will be displayed to users.

The following arguments are optional:

* `description` - (Optional) The description of the project. This description will be displayed to all users of the project. We recommend providing a brief description of the project and its intended purpose.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `name` - The name of the project in the space.

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

* `create` - (Default `60M`)
* `update` - (Default `180M`)
* `delete` - (Default `90M`)

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import CodeCatalyst Project using the `exampleIdArg`. For example:

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
  }
}

```

Using `terraform import`, import CodeCatalyst Project using the `exampleIdArg`. For example:

```console
% terraform import aws_codecatalyst_project.example project-id-12345678
```

<!-- cache-key: cdktf-0.18.0 input-458059ba4aa1ab978fd9d595f17ad727e6f517ced23cac7527a27f8c131037e1 -->