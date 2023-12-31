// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ssoadmin

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssoadmin"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tfslices "github.com/hashicorp/terraform-provider-aws/internal/slices"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKResource("aws_ssoadmin_customer_managed_policy_attachment")
func ResourceCustomerManagedPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCustomerManagedPolicyAttachmentCreate,
		ReadWithoutTimeout:   resourceCustomerManagedPolicyAttachmentRead,
		DeleteWithoutTimeout: resourceCustomerManagedPolicyAttachmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"customer_managed_policy_reference": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(0, 128),
						},
						"path": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "/",
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(0, 512),
						},
					},
				},
			},
			"instance_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"permission_set_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
		},
	}
}

func resourceCustomerManagedPolicyAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSOAdminConn(ctx)

	tfMap := d.Get("customer_managed_policy_reference").([]interface{})[0].(map[string]interface{})
	policyName := tfMap["name"].(string)
	policyPath := tfMap["path"].(string)
	instanceARN := d.Get("instance_arn").(string)
	permissionSetARN := d.Get("permission_set_arn").(string)
	id := CustomerManagedPolicyAttachmentCreateResourceID(policyName, policyPath, permissionSetARN, instanceARN)
	input := &ssoadmin.AttachCustomerManagedPolicyReferenceToPermissionSetInput{
		CustomerManagedPolicyReference: expandCustomerManagedPolicyReference(tfMap),
		InstanceArn:                    aws.String(instanceARN),
		PermissionSetArn:               aws.String(permissionSetARN),
	}

	_, err := conn.AttachCustomerManagedPolicyReferenceToPermissionSetWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating SSO Customer Managed Policy Attachment (%s): %s", id, err)
	}

	d.SetId(id)

	// After the policy has been attached to the permission set, provision in all accounts that use this permission set.
	if err := provisionPermissionSet(ctx, conn, permissionSetARN, instanceARN, d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	return append(diags, resourceCustomerManagedPolicyAttachmentRead(ctx, d, meta)...)
}

func resourceCustomerManagedPolicyAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSOAdminConn(ctx)

	policyName, policyPath, permissionSetARN, instanceARN, err := CustomerManagedPolicyAttachmentParseResourceID(d.Id())
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	policy, err := FindCustomerManagedPolicy(ctx, conn, policyName, policyPath, permissionSetARN, instanceARN)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] SSO Customer Managed Policy Attachment (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading SSO Customer Managed Policy Attachment (%s): %s", d.Id(), err)
	}

	if err := d.Set("customer_managed_policy_reference", []interface{}{flattenCustomerManagedPolicyReference(policy)}); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting customer_managed_policy_reference: %s", err)
	}
	d.Set("instance_arn", instanceARN)
	d.Set("permission_set_arn", permissionSetARN)

	return diags
}

func resourceCustomerManagedPolicyAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSOAdminConn(ctx)

	policyName, policyPath, permissionSetARN, instanceARN, err := CustomerManagedPolicyAttachmentParseResourceID(d.Id())
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	input := &ssoadmin.DetachCustomerManagedPolicyReferenceFromPermissionSetInput{
		CustomerManagedPolicyReference: &ssoadmin.CustomerManagedPolicyReference{
			Name: aws.String(policyName),
			Path: aws.String(policyPath),
		},
		InstanceArn:      aws.String(instanceARN),
		PermissionSetArn: aws.String(permissionSetARN),
	}

	log.Printf("[INFO] Deleting SSO Customer Managed Policy Attachment: %s", d.Id())
	_, err = conn.DetachCustomerManagedPolicyReferenceFromPermissionSetWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, ssoadmin.ErrCodeResourceNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting SSO Customer Managed Policy Attachment (%s): %s", d.Id(), err)
	}

	// After the policy has been detached from the permission set, provision in all accounts that use this permission set.
	if err := provisionPermissionSet(ctx, conn, permissionSetARN, instanceARN, d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	return diags
}

const customerManagedPolicyAttachmentIDSeparator = ","

func CustomerManagedPolicyAttachmentCreateResourceID(policyName, policyPath, permissionSetARN, instanceARN string) string {
	parts := []string{policyName, policyPath, permissionSetARN, instanceARN}
	id := strings.Join(parts, customerManagedPolicyAttachmentIDSeparator)

	return id
}

func CustomerManagedPolicyAttachmentParseResourceID(id string) (string, string, string, string, error) {
	parts := strings.Split(id, customerManagedPolicyAttachmentIDSeparator)

	if len(parts) == 4 && parts[0] != "" && parts[1] != "" && parts[2] != "" && parts[3] != "" {
		return parts[0], parts[1], parts[2], parts[3], nil
	}

	return "", "", "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected CUSTOMER_MANAGED_POLICY_NAME%[2]sCUSTOMER_MANAGED_POLICY_PATH%[2]sPERMISSION_SET_ARN%[2]sINSTANCE_ARN", id, customerManagedPolicyAttachmentIDSeparator)
}

func FindCustomerManagedPolicy(ctx context.Context, conn *ssoadmin.SSOAdmin, policyName, policyPath, permissionSetARN, instanceARN string) (*ssoadmin.CustomerManagedPolicyReference, error) {
	input := &ssoadmin.ListCustomerManagedPolicyReferencesInPermissionSetInput{
		InstanceArn:      aws.String(instanceARN),
		PermissionSetArn: aws.String(permissionSetARN),
	}
	filter := func(c *ssoadmin.CustomerManagedPolicyReference) bool {
		return aws.StringValue(c.Name) == policyName && aws.StringValue(c.Path) == policyPath
	}

	return findCustomerManagedPolicyReference(ctx, conn, input, filter)
}

func findCustomerManagedPolicyReference(ctx context.Context, conn *ssoadmin.SSOAdmin, input *ssoadmin.ListCustomerManagedPolicyReferencesInPermissionSetInput, filter tfslices.Predicate[*ssoadmin.CustomerManagedPolicyReference]) (*ssoadmin.CustomerManagedPolicyReference, error) {
	output, err := findCustomerManagedPolicyReferences(ctx, conn, input, filter)

	if err != nil {
		return nil, err
	}

	return tfresource.AssertSinglePtrResult(output)
}

func findCustomerManagedPolicyReferences(ctx context.Context, conn *ssoadmin.SSOAdmin, input *ssoadmin.ListCustomerManagedPolicyReferencesInPermissionSetInput, filter tfslices.Predicate[*ssoadmin.CustomerManagedPolicyReference]) ([]*ssoadmin.CustomerManagedPolicyReference, error) {
	var output []*ssoadmin.CustomerManagedPolicyReference

	err := conn.ListCustomerManagedPolicyReferencesInPermissionSetPagesWithContext(ctx, input, func(page *ssoadmin.ListCustomerManagedPolicyReferencesInPermissionSetOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.CustomerManagedPolicyReferences {
			if v != nil && filter(v) {
				output = append(output, v)
			}
		}

		return !lastPage
	})

	if tfawserr.ErrCodeEquals(err, ssoadmin.ErrCodeResourceNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	return output, nil
}

func expandCustomerManagedPolicyReference(tfMap map[string]interface{}) *ssoadmin.CustomerManagedPolicyReference {
	if tfMap == nil {
		return nil
	}

	apiObject := &ssoadmin.CustomerManagedPolicyReference{}

	if v, ok := tfMap["name"].(string); ok && v != "" {
		apiObject.Name = aws.String(v)
	}

	if v, ok := tfMap["path"].(string); ok && v != "" {
		apiObject.Path = aws.String(v)
	}

	return apiObject
}

func flattenCustomerManagedPolicyReference(apiObject *ssoadmin.CustomerManagedPolicyReference) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Name; v != nil {
		tfMap["name"] = aws.StringValue(v)
	}

	if v := apiObject.Path; v != nil {
		tfMap["path"] = aws.StringValue(v)
	}

	return tfMap
}
