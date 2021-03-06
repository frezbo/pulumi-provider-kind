// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package mount

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// KIND Mount type
type Mount struct {
	ContainerPath  *string `pulumi:"containerPath"`
	HostPath       *string `pulumi:"hostPath"`
	Propagation    *string `pulumi:"propagation"`
	ReadOnly       *bool   `pulumi:"readOnly"`
	SelinuxRelabel *bool   `pulumi:"selinuxRelabel"`
}

// MountInput is an input type that accepts MountArgs and MountOutput values.
// You can construct a concrete instance of `MountInput` via:
//
//          MountArgs{...}
type MountInput interface {
	pulumi.Input

	ToMountOutput() MountOutput
	ToMountOutputWithContext(context.Context) MountOutput
}

// KIND Mount type
type MountArgs struct {
	ContainerPath  pulumi.StringPtrInput `pulumi:"containerPath"`
	HostPath       pulumi.StringPtrInput `pulumi:"hostPath"`
	Propagation    pulumi.StringPtrInput `pulumi:"propagation"`
	ReadOnly       pulumi.BoolPtrInput   `pulumi:"readOnly"`
	SelinuxRelabel pulumi.BoolPtrInput   `pulumi:"selinuxRelabel"`
}

func (MountArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*Mount)(nil)).Elem()
}

func (i MountArgs) ToMountOutput() MountOutput {
	return i.ToMountOutputWithContext(context.Background())
}

func (i MountArgs) ToMountOutputWithContext(ctx context.Context) MountOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MountOutput)
}

// MountArrayInput is an input type that accepts MountArray and MountArrayOutput values.
// You can construct a concrete instance of `MountArrayInput` via:
//
//          MountArray{ MountArgs{...} }
type MountArrayInput interface {
	pulumi.Input

	ToMountArrayOutput() MountArrayOutput
	ToMountArrayOutputWithContext(context.Context) MountArrayOutput
}

type MountArray []MountInput

func (MountArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]Mount)(nil)).Elem()
}

func (i MountArray) ToMountArrayOutput() MountArrayOutput {
	return i.ToMountArrayOutputWithContext(context.Background())
}

func (i MountArray) ToMountArrayOutputWithContext(ctx context.Context) MountArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MountArrayOutput)
}

// KIND Mount type
type MountOutput struct{ *pulumi.OutputState }

func (MountOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Mount)(nil)).Elem()
}

func (o MountOutput) ToMountOutput() MountOutput {
	return o
}

func (o MountOutput) ToMountOutputWithContext(ctx context.Context) MountOutput {
	return o
}

func (o MountOutput) ContainerPath() pulumi.StringPtrOutput {
	return o.ApplyT(func(v Mount) *string { return v.ContainerPath }).(pulumi.StringPtrOutput)
}

func (o MountOutput) HostPath() pulumi.StringPtrOutput {
	return o.ApplyT(func(v Mount) *string { return v.HostPath }).(pulumi.StringPtrOutput)
}

func (o MountOutput) Propagation() pulumi.StringPtrOutput {
	return o.ApplyT(func(v Mount) *string { return v.Propagation }).(pulumi.StringPtrOutput)
}

func (o MountOutput) ReadOnly() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v Mount) *bool { return v.ReadOnly }).(pulumi.BoolPtrOutput)
}

func (o MountOutput) SelinuxRelabel() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v Mount) *bool { return v.SelinuxRelabel }).(pulumi.BoolPtrOutput)
}

type MountArrayOutput struct{ *pulumi.OutputState }

func (MountArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]Mount)(nil)).Elem()
}

func (o MountArrayOutput) ToMountArrayOutput() MountArrayOutput {
	return o
}

func (o MountArrayOutput) ToMountArrayOutputWithContext(ctx context.Context) MountArrayOutput {
	return o
}

func (o MountArrayOutput) Index(i pulumi.IntInput) MountOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) Mount {
		return vs[0].([]Mount)[vs[1].(int)]
	}).(MountOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*MountInput)(nil)).Elem(), MountArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*MountArrayInput)(nil)).Elem(), MountArray{})
	pulumi.RegisterOutputType(MountOutput{})
	pulumi.RegisterOutputType(MountArrayOutput{})
}
