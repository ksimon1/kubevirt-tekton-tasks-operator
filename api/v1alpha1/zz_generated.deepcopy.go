//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Pipelines) DeepCopyInto(out *Pipelines) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Pipelines.
func (in *Pipelines) DeepCopy() *Pipelines {
	if in == nil {
		return nil
	}
	out := new(Pipelines)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TektonTasks) DeepCopyInto(out *TektonTasks) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TektonTasks.
func (in *TektonTasks) DeepCopy() *TektonTasks {
	if in == nil {
		return nil
	}
	out := new(TektonTasks)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TektonTasks) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TektonTasksList) DeepCopyInto(out *TektonTasksList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TektonTasks, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TektonTasksList.
func (in *TektonTasksList) DeepCopy() *TektonTasksList {
	if in == nil {
		return nil
	}
	out := new(TektonTasksList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TektonTasksList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TektonTasksSpec) DeepCopyInto(out *TektonTasksSpec) {
	*out = *in
	out.Pipelines = in.Pipelines
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TektonTasksSpec.
func (in *TektonTasksSpec) DeepCopy() *TektonTasksSpec {
	if in == nil {
		return nil
	}
	out := new(TektonTasksSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TektonTasksStatus) DeepCopyInto(out *TektonTasksStatus) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TektonTasksStatus.
func (in *TektonTasksStatus) DeepCopy() *TektonTasksStatus {
	if in == nil {
		return nil
	}
	out := new(TektonTasksStatus)
	in.DeepCopyInto(out)
	return out
}
