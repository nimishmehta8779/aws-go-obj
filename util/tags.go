package util

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewNameTags(ctx *pulumi.Context, name string) map[string]string {
	tags := NewCommonTags(ctx)
	tags["Name"] = fmt.Sprintf("%s-%s-%s", ctx.Project(), ctx.Stack(), name)
	return tags
}

func NewCommonTags(ctx *pulumi.Context) map[string]string {
	return map[string]string{
		"Project":     ctx.Project(),
		"Environment": ctx.Stack(),
		"Name":        fmt.Sprintf("%s-%s", ctx.Project(), ctx.Stack()),
	}
}
