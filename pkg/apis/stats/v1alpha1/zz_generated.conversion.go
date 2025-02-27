//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Copyright 2021 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	stats "antrea.io/antrea/pkg/apis/stats"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*AntreaClusterNetworkPolicyStats)(nil), (*stats.AntreaClusterNetworkPolicyStats)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_AntreaClusterNetworkPolicyStats_To_stats_AntreaClusterNetworkPolicyStats(a.(*AntreaClusterNetworkPolicyStats), b.(*stats.AntreaClusterNetworkPolicyStats), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*stats.AntreaClusterNetworkPolicyStats)(nil), (*AntreaClusterNetworkPolicyStats)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_stats_AntreaClusterNetworkPolicyStats_To_v1alpha1_AntreaClusterNetworkPolicyStats(a.(*stats.AntreaClusterNetworkPolicyStats), b.(*AntreaClusterNetworkPolicyStats), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*AntreaClusterNetworkPolicyStatsList)(nil), (*stats.AntreaClusterNetworkPolicyStatsList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_AntreaClusterNetworkPolicyStatsList_To_stats_AntreaClusterNetworkPolicyStatsList(a.(*AntreaClusterNetworkPolicyStatsList), b.(*stats.AntreaClusterNetworkPolicyStatsList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*stats.AntreaClusterNetworkPolicyStatsList)(nil), (*AntreaClusterNetworkPolicyStatsList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_stats_AntreaClusterNetworkPolicyStatsList_To_v1alpha1_AntreaClusterNetworkPolicyStatsList(a.(*stats.AntreaClusterNetworkPolicyStatsList), b.(*AntreaClusterNetworkPolicyStatsList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*AntreaNetworkPolicyStats)(nil), (*stats.AntreaNetworkPolicyStats)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_AntreaNetworkPolicyStats_To_stats_AntreaNetworkPolicyStats(a.(*AntreaNetworkPolicyStats), b.(*stats.AntreaNetworkPolicyStats), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*stats.AntreaNetworkPolicyStats)(nil), (*AntreaNetworkPolicyStats)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_stats_AntreaNetworkPolicyStats_To_v1alpha1_AntreaNetworkPolicyStats(a.(*stats.AntreaNetworkPolicyStats), b.(*AntreaNetworkPolicyStats), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*AntreaNetworkPolicyStatsList)(nil), (*stats.AntreaNetworkPolicyStatsList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_AntreaNetworkPolicyStatsList_To_stats_AntreaNetworkPolicyStatsList(a.(*AntreaNetworkPolicyStatsList), b.(*stats.AntreaNetworkPolicyStatsList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*stats.AntreaNetworkPolicyStatsList)(nil), (*AntreaNetworkPolicyStatsList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_stats_AntreaNetworkPolicyStatsList_To_v1alpha1_AntreaNetworkPolicyStatsList(a.(*stats.AntreaNetworkPolicyStatsList), b.(*AntreaNetworkPolicyStatsList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*NetworkPolicyStats)(nil), (*stats.NetworkPolicyStats)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_NetworkPolicyStats_To_stats_NetworkPolicyStats(a.(*NetworkPolicyStats), b.(*stats.NetworkPolicyStats), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*stats.NetworkPolicyStats)(nil), (*NetworkPolicyStats)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_stats_NetworkPolicyStats_To_v1alpha1_NetworkPolicyStats(a.(*stats.NetworkPolicyStats), b.(*NetworkPolicyStats), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*NetworkPolicyStatsList)(nil), (*stats.NetworkPolicyStatsList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_NetworkPolicyStatsList_To_stats_NetworkPolicyStatsList(a.(*NetworkPolicyStatsList), b.(*stats.NetworkPolicyStatsList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*stats.NetworkPolicyStatsList)(nil), (*NetworkPolicyStatsList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_stats_NetworkPolicyStatsList_To_v1alpha1_NetworkPolicyStatsList(a.(*stats.NetworkPolicyStatsList), b.(*NetworkPolicyStatsList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*RuleTrafficStats)(nil), (*stats.RuleTrafficStats)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_RuleTrafficStats_To_stats_RuleTrafficStats(a.(*RuleTrafficStats), b.(*stats.RuleTrafficStats), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*stats.RuleTrafficStats)(nil), (*RuleTrafficStats)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_stats_RuleTrafficStats_To_v1alpha1_RuleTrafficStats(a.(*stats.RuleTrafficStats), b.(*RuleTrafficStats), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*TrafficStats)(nil), (*stats.TrafficStats)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_TrafficStats_To_stats_TrafficStats(a.(*TrafficStats), b.(*stats.TrafficStats), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*stats.TrafficStats)(nil), (*TrafficStats)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_stats_TrafficStats_To_v1alpha1_TrafficStats(a.(*stats.TrafficStats), b.(*TrafficStats), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_AntreaClusterNetworkPolicyStats_To_stats_AntreaClusterNetworkPolicyStats(in *AntreaClusterNetworkPolicyStats, out *stats.AntreaClusterNetworkPolicyStats, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_TrafficStats_To_stats_TrafficStats(&in.TrafficStats, &out.TrafficStats, s); err != nil {
		return err
	}
	out.RuleTrafficStats = *(*[]stats.RuleTrafficStats)(unsafe.Pointer(&in.RuleTrafficStats))
	return nil
}

// Convert_v1alpha1_AntreaClusterNetworkPolicyStats_To_stats_AntreaClusterNetworkPolicyStats is an autogenerated conversion function.
func Convert_v1alpha1_AntreaClusterNetworkPolicyStats_To_stats_AntreaClusterNetworkPolicyStats(in *AntreaClusterNetworkPolicyStats, out *stats.AntreaClusterNetworkPolicyStats, s conversion.Scope) error {
	return autoConvert_v1alpha1_AntreaClusterNetworkPolicyStats_To_stats_AntreaClusterNetworkPolicyStats(in, out, s)
}

func autoConvert_stats_AntreaClusterNetworkPolicyStats_To_v1alpha1_AntreaClusterNetworkPolicyStats(in *stats.AntreaClusterNetworkPolicyStats, out *AntreaClusterNetworkPolicyStats, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_stats_TrafficStats_To_v1alpha1_TrafficStats(&in.TrafficStats, &out.TrafficStats, s); err != nil {
		return err
	}
	out.RuleTrafficStats = *(*[]RuleTrafficStats)(unsafe.Pointer(&in.RuleTrafficStats))
	return nil
}

// Convert_stats_AntreaClusterNetworkPolicyStats_To_v1alpha1_AntreaClusterNetworkPolicyStats is an autogenerated conversion function.
func Convert_stats_AntreaClusterNetworkPolicyStats_To_v1alpha1_AntreaClusterNetworkPolicyStats(in *stats.AntreaClusterNetworkPolicyStats, out *AntreaClusterNetworkPolicyStats, s conversion.Scope) error {
	return autoConvert_stats_AntreaClusterNetworkPolicyStats_To_v1alpha1_AntreaClusterNetworkPolicyStats(in, out, s)
}

func autoConvert_v1alpha1_AntreaClusterNetworkPolicyStatsList_To_stats_AntreaClusterNetworkPolicyStatsList(in *AntreaClusterNetworkPolicyStatsList, out *stats.AntreaClusterNetworkPolicyStatsList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]stats.AntreaClusterNetworkPolicyStats)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_AntreaClusterNetworkPolicyStatsList_To_stats_AntreaClusterNetworkPolicyStatsList is an autogenerated conversion function.
func Convert_v1alpha1_AntreaClusterNetworkPolicyStatsList_To_stats_AntreaClusterNetworkPolicyStatsList(in *AntreaClusterNetworkPolicyStatsList, out *stats.AntreaClusterNetworkPolicyStatsList, s conversion.Scope) error {
	return autoConvert_v1alpha1_AntreaClusterNetworkPolicyStatsList_To_stats_AntreaClusterNetworkPolicyStatsList(in, out, s)
}

func autoConvert_stats_AntreaClusterNetworkPolicyStatsList_To_v1alpha1_AntreaClusterNetworkPolicyStatsList(in *stats.AntreaClusterNetworkPolicyStatsList, out *AntreaClusterNetworkPolicyStatsList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]AntreaClusterNetworkPolicyStats)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_stats_AntreaClusterNetworkPolicyStatsList_To_v1alpha1_AntreaClusterNetworkPolicyStatsList is an autogenerated conversion function.
func Convert_stats_AntreaClusterNetworkPolicyStatsList_To_v1alpha1_AntreaClusterNetworkPolicyStatsList(in *stats.AntreaClusterNetworkPolicyStatsList, out *AntreaClusterNetworkPolicyStatsList, s conversion.Scope) error {
	return autoConvert_stats_AntreaClusterNetworkPolicyStatsList_To_v1alpha1_AntreaClusterNetworkPolicyStatsList(in, out, s)
}

func autoConvert_v1alpha1_AntreaNetworkPolicyStats_To_stats_AntreaNetworkPolicyStats(in *AntreaNetworkPolicyStats, out *stats.AntreaNetworkPolicyStats, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_TrafficStats_To_stats_TrafficStats(&in.TrafficStats, &out.TrafficStats, s); err != nil {
		return err
	}
	out.RuleTrafficStats = *(*[]stats.RuleTrafficStats)(unsafe.Pointer(&in.RuleTrafficStats))
	return nil
}

// Convert_v1alpha1_AntreaNetworkPolicyStats_To_stats_AntreaNetworkPolicyStats is an autogenerated conversion function.
func Convert_v1alpha1_AntreaNetworkPolicyStats_To_stats_AntreaNetworkPolicyStats(in *AntreaNetworkPolicyStats, out *stats.AntreaNetworkPolicyStats, s conversion.Scope) error {
	return autoConvert_v1alpha1_AntreaNetworkPolicyStats_To_stats_AntreaNetworkPolicyStats(in, out, s)
}

func autoConvert_stats_AntreaNetworkPolicyStats_To_v1alpha1_AntreaNetworkPolicyStats(in *stats.AntreaNetworkPolicyStats, out *AntreaNetworkPolicyStats, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_stats_TrafficStats_To_v1alpha1_TrafficStats(&in.TrafficStats, &out.TrafficStats, s); err != nil {
		return err
	}
	out.RuleTrafficStats = *(*[]RuleTrafficStats)(unsafe.Pointer(&in.RuleTrafficStats))
	return nil
}

// Convert_stats_AntreaNetworkPolicyStats_To_v1alpha1_AntreaNetworkPolicyStats is an autogenerated conversion function.
func Convert_stats_AntreaNetworkPolicyStats_To_v1alpha1_AntreaNetworkPolicyStats(in *stats.AntreaNetworkPolicyStats, out *AntreaNetworkPolicyStats, s conversion.Scope) error {
	return autoConvert_stats_AntreaNetworkPolicyStats_To_v1alpha1_AntreaNetworkPolicyStats(in, out, s)
}

func autoConvert_v1alpha1_AntreaNetworkPolicyStatsList_To_stats_AntreaNetworkPolicyStatsList(in *AntreaNetworkPolicyStatsList, out *stats.AntreaNetworkPolicyStatsList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]stats.AntreaNetworkPolicyStats)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_AntreaNetworkPolicyStatsList_To_stats_AntreaNetworkPolicyStatsList is an autogenerated conversion function.
func Convert_v1alpha1_AntreaNetworkPolicyStatsList_To_stats_AntreaNetworkPolicyStatsList(in *AntreaNetworkPolicyStatsList, out *stats.AntreaNetworkPolicyStatsList, s conversion.Scope) error {
	return autoConvert_v1alpha1_AntreaNetworkPolicyStatsList_To_stats_AntreaNetworkPolicyStatsList(in, out, s)
}

func autoConvert_stats_AntreaNetworkPolicyStatsList_To_v1alpha1_AntreaNetworkPolicyStatsList(in *stats.AntreaNetworkPolicyStatsList, out *AntreaNetworkPolicyStatsList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]AntreaNetworkPolicyStats)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_stats_AntreaNetworkPolicyStatsList_To_v1alpha1_AntreaNetworkPolicyStatsList is an autogenerated conversion function.
func Convert_stats_AntreaNetworkPolicyStatsList_To_v1alpha1_AntreaNetworkPolicyStatsList(in *stats.AntreaNetworkPolicyStatsList, out *AntreaNetworkPolicyStatsList, s conversion.Scope) error {
	return autoConvert_stats_AntreaNetworkPolicyStatsList_To_v1alpha1_AntreaNetworkPolicyStatsList(in, out, s)
}

func autoConvert_v1alpha1_NetworkPolicyStats_To_stats_NetworkPolicyStats(in *NetworkPolicyStats, out *stats.NetworkPolicyStats, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_TrafficStats_To_stats_TrafficStats(&in.TrafficStats, &out.TrafficStats, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_NetworkPolicyStats_To_stats_NetworkPolicyStats is an autogenerated conversion function.
func Convert_v1alpha1_NetworkPolicyStats_To_stats_NetworkPolicyStats(in *NetworkPolicyStats, out *stats.NetworkPolicyStats, s conversion.Scope) error {
	return autoConvert_v1alpha1_NetworkPolicyStats_To_stats_NetworkPolicyStats(in, out, s)
}

func autoConvert_stats_NetworkPolicyStats_To_v1alpha1_NetworkPolicyStats(in *stats.NetworkPolicyStats, out *NetworkPolicyStats, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_stats_TrafficStats_To_v1alpha1_TrafficStats(&in.TrafficStats, &out.TrafficStats, s); err != nil {
		return err
	}
	return nil
}

// Convert_stats_NetworkPolicyStats_To_v1alpha1_NetworkPolicyStats is an autogenerated conversion function.
func Convert_stats_NetworkPolicyStats_To_v1alpha1_NetworkPolicyStats(in *stats.NetworkPolicyStats, out *NetworkPolicyStats, s conversion.Scope) error {
	return autoConvert_stats_NetworkPolicyStats_To_v1alpha1_NetworkPolicyStats(in, out, s)
}

func autoConvert_v1alpha1_NetworkPolicyStatsList_To_stats_NetworkPolicyStatsList(in *NetworkPolicyStatsList, out *stats.NetworkPolicyStatsList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]stats.NetworkPolicyStats)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_NetworkPolicyStatsList_To_stats_NetworkPolicyStatsList is an autogenerated conversion function.
func Convert_v1alpha1_NetworkPolicyStatsList_To_stats_NetworkPolicyStatsList(in *NetworkPolicyStatsList, out *stats.NetworkPolicyStatsList, s conversion.Scope) error {
	return autoConvert_v1alpha1_NetworkPolicyStatsList_To_stats_NetworkPolicyStatsList(in, out, s)
}

func autoConvert_stats_NetworkPolicyStatsList_To_v1alpha1_NetworkPolicyStatsList(in *stats.NetworkPolicyStatsList, out *NetworkPolicyStatsList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]NetworkPolicyStats)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_stats_NetworkPolicyStatsList_To_v1alpha1_NetworkPolicyStatsList is an autogenerated conversion function.
func Convert_stats_NetworkPolicyStatsList_To_v1alpha1_NetworkPolicyStatsList(in *stats.NetworkPolicyStatsList, out *NetworkPolicyStatsList, s conversion.Scope) error {
	return autoConvert_stats_NetworkPolicyStatsList_To_v1alpha1_NetworkPolicyStatsList(in, out, s)
}

func autoConvert_v1alpha1_RuleTrafficStats_To_stats_RuleTrafficStats(in *RuleTrafficStats, out *stats.RuleTrafficStats, s conversion.Scope) error {
	out.Name = in.Name
	if err := Convert_v1alpha1_TrafficStats_To_stats_TrafficStats(&in.TrafficStats, &out.TrafficStats, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_RuleTrafficStats_To_stats_RuleTrafficStats is an autogenerated conversion function.
func Convert_v1alpha1_RuleTrafficStats_To_stats_RuleTrafficStats(in *RuleTrafficStats, out *stats.RuleTrafficStats, s conversion.Scope) error {
	return autoConvert_v1alpha1_RuleTrafficStats_To_stats_RuleTrafficStats(in, out, s)
}

func autoConvert_stats_RuleTrafficStats_To_v1alpha1_RuleTrafficStats(in *stats.RuleTrafficStats, out *RuleTrafficStats, s conversion.Scope) error {
	out.Name = in.Name
	if err := Convert_stats_TrafficStats_To_v1alpha1_TrafficStats(&in.TrafficStats, &out.TrafficStats, s); err != nil {
		return err
	}
	return nil
}

// Convert_stats_RuleTrafficStats_To_v1alpha1_RuleTrafficStats is an autogenerated conversion function.
func Convert_stats_RuleTrafficStats_To_v1alpha1_RuleTrafficStats(in *stats.RuleTrafficStats, out *RuleTrafficStats, s conversion.Scope) error {
	return autoConvert_stats_RuleTrafficStats_To_v1alpha1_RuleTrafficStats(in, out, s)
}

func autoConvert_v1alpha1_TrafficStats_To_stats_TrafficStats(in *TrafficStats, out *stats.TrafficStats, s conversion.Scope) error {
	out.Packets = in.Packets
	out.Bytes = in.Bytes
	out.Sessions = in.Sessions
	return nil
}

// Convert_v1alpha1_TrafficStats_To_stats_TrafficStats is an autogenerated conversion function.
func Convert_v1alpha1_TrafficStats_To_stats_TrafficStats(in *TrafficStats, out *stats.TrafficStats, s conversion.Scope) error {
	return autoConvert_v1alpha1_TrafficStats_To_stats_TrafficStats(in, out, s)
}

func autoConvert_stats_TrafficStats_To_v1alpha1_TrafficStats(in *stats.TrafficStats, out *TrafficStats, s conversion.Scope) error {
	out.Packets = in.Packets
	out.Bytes = in.Bytes
	out.Sessions = in.Sessions
	return nil
}

// Convert_stats_TrafficStats_To_v1alpha1_TrafficStats is an autogenerated conversion function.
func Convert_stats_TrafficStats_To_v1alpha1_TrafficStats(in *stats.TrafficStats, out *TrafficStats, s conversion.Scope) error {
	return autoConvert_stats_TrafficStats_To_v1alpha1_TrafficStats(in, out, s)
}
