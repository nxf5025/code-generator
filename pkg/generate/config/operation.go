// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package config

import (
	"encoding/json"

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"

	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

// StringArray is a type that can be represented in JSON as *either* a string
// *or* an array of strings
type StringArray []string

// OperationConfig represents instructions to the ACK code generator to
// specify the overriding values for API operation parameters and its custom implementation.
type OperationConfig struct {
	CustomImplementation                   string            `json:"custom_implementation,omitempty"`
	CustomCheckRequiredFieldsMissingMethod string            `json:"custom_check_required_fields_missing_method,omitempty"`
	OverrideValues                         map[string]string `json:"override_values"`
	// SetOutputCustomMethodName provides the name of the custom method on the
	// `resourceManager` struct that will set fields on a `resource` struct
	// depending on the output of the operation.
	SetOutputCustomMethodName string `json:"set_output_custom_method_name,omitempty"`
	// OutputWrapperFieldPath provides the JSON-Path like to the struct field containing
	// information that will be merged into a `resource` object.
	OutputWrapperFieldPath string `json:"output_wrapper_field_path,omitempty"`
	// Override for resource name in case of heuristic failure
	// An example of this is correcting stutter when the resource logic doesn't properly determine the resource name
	ResourceName string `json:"resource_name"`
	// Override for operation type in case of heuristic failure
	// An example of this is `Put...` or `Register...` API operations not being correctly classified as `Create` op type
	// OperationType []string `json:"operation_type"`
	OperationType StringArray `json:"operation_type"`
}

// IsIgnoredOperation returns true if Operation Name is configured to be ignored
// in generator config for the AWS service
func (c *Config) IsIgnoredOperation(operation *awssdkmodel.Operation) bool {
	if c == nil {
		return false
	}
	if operation == nil {
		return true
	}
	return util.InStrings(operation.Name, c.Ignore.Operations)
}

// ListOpMatchFieldNames returns a slice of strings representing the field
// names in the List operation's Output shape's element Shape that we should
// check a corresponding value in the target Spec exists.
func (c *Config) ListOpMatchFieldNames(
	resName string,
) []string {
	res := []string{}
	if c == nil {
		return res
	}
	rConfig, found := c.Resources[resName]
	if !found {
		return res
	}
	if rConfig.ListOperation == nil {
		return res
	}
	return rConfig.ListOperation.MatchFields
}

// UnmarshalJSON parses input for a either a string or
// or a list and returns a StringArray.
func (a *StringArray) UnmarshalJSON(b []byte) error {
	var multi []string
	err := json.Unmarshal(b, &multi)
	if err != nil {
		var single string
		err := json.Unmarshal(b, &single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}
