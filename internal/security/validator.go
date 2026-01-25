package security

import (
	"fmt"
	"strings"
)

var (
	protectedBranches = []string{"master", "main", "test", "develop", "release"}
)

type OperationType int

const (
	OperationRead OperationType = iota
	OperationCreate
	OperationUpdate
	OperationDelete
	OperationMerge
)

type ValidationResult struct {
	Allowed bool
	Reason  string
}

type Operation struct {
	Type            OperationType
	TargetBranch    string
	TargetProjectID string
	OperationName   string
}

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateOperation(op Operation) *ValidationResult {
	switch op.Type {
	case OperationRead:
		return v.validateReadOperation(op)
	case OperationCreate:
		return v.validateCreateOperation(op)
	case OperationUpdate:
		return v.validateUpdateOperation(op)
	case OperationDelete:
		return v.validateDeleteOperation(op)
	case OperationMerge:
		return v.validateMergeOperation(op)
	default:
		return &ValidationResult{Allowed: false, Reason: "unknown operation type"}
	}
}

func (v *Validator) validateReadOperation(op Operation) *ValidationResult {
	return &ValidationResult{Allowed: true, Reason: ""}
}

func (v *Validator) validateCreateOperation(op Operation) *ValidationResult {
	if v.isProtectedBranch(op.TargetBranch) {
		return &ValidationResult{
			Allowed: false,
			Reason: fmt.Sprintf("creating or modifying protected branch '%s' is not allowed", op.TargetBranch),
		}
	}

	if op.OperationName == "create_branch" {
		if v.isProtectedBranch(op.TargetBranch) {
			return &ValidationResult{
				Allowed: false,
				Reason: fmt.Sprintf("creating branch '%s' is not allowed (protected branch)", op.TargetBranch),
			}
		}
	}

	return &ValidationResult{Allowed: true, Reason: ""}
}

func (v *Validator) validateUpdateOperation(op Operation) *ValidationResult {
	if v.isProtectedBranch(op.TargetBranch) {
		return &ValidationResult{
			Allowed: false,
			Reason: fmt.Sprintf("modifying protected branch '%s' is not allowed", op.TargetBranch),
		}
	}

	if op.OperationName == "update_file" && v.isProtectedBranch(op.TargetBranch) {
		return &ValidationResult{
			Allowed: false,
			Reason: fmt.Sprintf("updating files on protected branch '%s' is not allowed", op.TargetBranch),
		}
	}

	if op.OperationName == "update_mr" {
		return &ValidationResult{
			Allowed: false,
			Reason: "modifying merge request settings is not allowed",
		}
	}

	return &ValidationResult{Allowed: true, Reason: ""}
}

func (v *Validator) validateDeleteOperation(op Operation) *ValidationResult {
	return &ValidationResult{
		Allowed: false,
		Reason: "delete operations are not allowed",
	}
}

func (v *Validator) validateMergeOperation(op Operation) *ValidationResult {
	return &ValidationResult{
		Allowed: false,
		Reason: "merge operations are not allowed",
	}
}

func (v *Validator) isProtectedBranch(branch string) bool {
	if branch == "" {
		return false
	}
	branchLower := strings.ToLower(branch)
	for _, protected := range protectedBranches {
		if branchLower == protected {
			return true
		}
	}
	return false
}

func (v *Validator) CanCreateMR(sourceBranch, targetBranch string) *ValidationResult {
	if v.isProtectedBranch(sourceBranch) {
		return &ValidationResult{
			Allowed: false,
			Reason: fmt.Sprintf("creating MR from protected branch '%s' is not allowed", sourceBranch),
		}
	}
	if v.isProtectedBranch(targetBranch) {
		return &ValidationResult{
			Allowed: false,
			Reason: fmt.Sprintf("creating MR to protected branch '%s' is not allowed", targetBranch),
		}
	}
	return &ValidationResult{Allowed: true, Reason: ""}
}

func (v *Validator) CanPushCode(branch string) *ValidationResult {
	if v.isProtectedBranch(branch) {
		return &ValidationResult{
			Allowed: false,
			Reason: fmt.Sprintf("pushing code to protected branch '%s' is not allowed", branch),
		}
	}
	return &ValidationResult{Allowed: true, Reason: ""}
}

func (v *Validator) CanCreateBranch(branch string) *ValidationResult {
	if v.isProtectedBranch(branch) {
		return &ValidationResult{
			Allowed: false,
			Reason: fmt.Sprintf("creating protected branch '%s' is not allowed", branch),
		}
	}
	return &ValidationResult{Allowed: true, Reason: ""}
}

func (v *Validator) CanDeleteResource(resourceType string) *ValidationResult {
	return &ValidationResult{
		Allowed: false,
		Reason: fmt.Sprintf("deleting %s is not allowed", resourceType),
	}
}

func (v *Validator) CanMergeMR() *ValidationResult {
	return &ValidationResult{
		Allowed: false,
		Reason: "merging merge requests is not allowed",
	}
}

func (v *Validator) CanModifyProjectSettings() *ValidationResult {
	return &ValidationResult{
		Allowed: false,
		Reason: "modifying project settings is not allowed",
	}
}
