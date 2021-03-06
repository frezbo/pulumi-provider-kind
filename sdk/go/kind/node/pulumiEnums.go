// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package node

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type RoleType string

const (
	// node that hosts Kubernetes control-plane components
	RoleTypeControlPlane = RoleType("control-plane")
	// node that hosts Kubernetes worker
	RoleTypeWorker = RoleType("worker")
	// node that hosts an external load balancer
	RoleTypeLoadBalancer = RoleType("external-load-balancer")
)

func (RoleType) ElementType() reflect.Type {
	return reflect.TypeOf((*RoleType)(nil)).Elem()
}

func (e RoleType) ToRoleTypeOutput() RoleTypeOutput {
	return pulumi.ToOutput(e).(RoleTypeOutput)
}

func (e RoleType) ToRoleTypeOutputWithContext(ctx context.Context) RoleTypeOutput {
	return pulumi.ToOutputWithContext(ctx, e).(RoleTypeOutput)
}

func (e RoleType) ToRoleTypePtrOutput() RoleTypePtrOutput {
	return e.ToRoleTypePtrOutputWithContext(context.Background())
}

func (e RoleType) ToRoleTypePtrOutputWithContext(ctx context.Context) RoleTypePtrOutput {
	return RoleType(e).ToRoleTypeOutputWithContext(ctx).ToRoleTypePtrOutputWithContext(ctx)
}

func (e RoleType) ToStringOutput() pulumi.StringOutput {
	return pulumi.ToOutput(pulumi.String(e)).(pulumi.StringOutput)
}

func (e RoleType) ToStringOutputWithContext(ctx context.Context) pulumi.StringOutput {
	return pulumi.ToOutputWithContext(ctx, pulumi.String(e)).(pulumi.StringOutput)
}

func (e RoleType) ToStringPtrOutput() pulumi.StringPtrOutput {
	return pulumi.String(e).ToStringPtrOutputWithContext(context.Background())
}

func (e RoleType) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return pulumi.String(e).ToStringOutputWithContext(ctx).ToStringPtrOutputWithContext(ctx)
}

type RoleTypeOutput struct{ *pulumi.OutputState }

func (RoleTypeOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*RoleType)(nil)).Elem()
}

func (o RoleTypeOutput) ToRoleTypeOutput() RoleTypeOutput {
	return o
}

func (o RoleTypeOutput) ToRoleTypeOutputWithContext(ctx context.Context) RoleTypeOutput {
	return o
}

func (o RoleTypeOutput) ToRoleTypePtrOutput() RoleTypePtrOutput {
	return o.ToRoleTypePtrOutputWithContext(context.Background())
}

func (o RoleTypeOutput) ToRoleTypePtrOutputWithContext(ctx context.Context) RoleTypePtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v RoleType) *RoleType {
		return &v
	}).(RoleTypePtrOutput)
}

func (o RoleTypeOutput) ToStringOutput() pulumi.StringOutput {
	return o.ToStringOutputWithContext(context.Background())
}

func (o RoleTypeOutput) ToStringOutputWithContext(ctx context.Context) pulumi.StringOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e RoleType) string {
		return string(e)
	}).(pulumi.StringOutput)
}

func (o RoleTypeOutput) ToStringPtrOutput() pulumi.StringPtrOutput {
	return o.ToStringPtrOutputWithContext(context.Background())
}

func (o RoleTypeOutput) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e RoleType) *string {
		v := string(e)
		return &v
	}).(pulumi.StringPtrOutput)
}

type RoleTypePtrOutput struct{ *pulumi.OutputState }

func (RoleTypePtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**RoleType)(nil)).Elem()
}

func (o RoleTypePtrOutput) ToRoleTypePtrOutput() RoleTypePtrOutput {
	return o
}

func (o RoleTypePtrOutput) ToRoleTypePtrOutputWithContext(ctx context.Context) RoleTypePtrOutput {
	return o
}

func (o RoleTypePtrOutput) Elem() RoleTypeOutput {
	return o.ApplyT(func(v *RoleType) RoleType {
		if v != nil {
			return *v
		}
		var ret RoleType
		return ret
	}).(RoleTypeOutput)
}

func (o RoleTypePtrOutput) ToStringPtrOutput() pulumi.StringPtrOutput {
	return o.ToStringPtrOutputWithContext(context.Background())
}

func (o RoleTypePtrOutput) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e *RoleType) *string {
		if e == nil {
			return nil
		}
		v := string(*e)
		return &v
	}).(pulumi.StringPtrOutput)
}

// RoleTypeInput is an input type that accepts RoleTypeArgs and RoleTypeOutput values.
// You can construct a concrete instance of `RoleTypeInput` via:
//
//          RoleTypeArgs{...}
type RoleTypeInput interface {
	pulumi.Input

	ToRoleTypeOutput() RoleTypeOutput
	ToRoleTypeOutputWithContext(context.Context) RoleTypeOutput
}

var roleTypePtrType = reflect.TypeOf((**RoleType)(nil)).Elem()

type RoleTypePtrInput interface {
	pulumi.Input

	ToRoleTypePtrOutput() RoleTypePtrOutput
	ToRoleTypePtrOutputWithContext(context.Context) RoleTypePtrOutput
}

type roleTypePtr string

func RoleTypePtr(v string) RoleTypePtrInput {
	return (*roleTypePtr)(&v)
}

func (*roleTypePtr) ElementType() reflect.Type {
	return roleTypePtrType
}

func (in *roleTypePtr) ToRoleTypePtrOutput() RoleTypePtrOutput {
	return pulumi.ToOutput(in).(RoleTypePtrOutput)
}

func (in *roleTypePtr) ToRoleTypePtrOutputWithContext(ctx context.Context) RoleTypePtrOutput {
	return pulumi.ToOutputWithContext(ctx, in).(RoleTypePtrOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*RoleTypeInput)(nil)).Elem(), RoleType("control-plane"))
	pulumi.RegisterInputType(reflect.TypeOf((*RoleTypePtrInput)(nil)).Elem(), RoleType("control-plane"))
	pulumi.RegisterOutputType(RoleTypeOutput{})
	pulumi.RegisterOutputType(RoleTypePtrOutput{})
}
