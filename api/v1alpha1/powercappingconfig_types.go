/*
Copyright 2024.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PowerCappingSpecKind string

const (
	NoPowerCappingSpec                                        PowerCappingSpecKind = "NoPowerCappingSpec"
	AbsolutePowerCapInWatts                                   PowerCappingSpecKind = "AbsolutePowerCapInWatts"
	RelativePowerCapOfPeakPowerConsumptionInPercentage        PowerCappingSpecKind = "RelativePowerCapOfPeakPowerConsumptionInPercentage"
	RelativePowerCappingOfAveragePowerConsumptionInPercentage PowerCappingSpecKind = "RelativePowerCappingOfAveragePowerConsumptionInPercentage"
)

type AbsolutePowerCapInWattsSpec struct {
	PowerCapInWatts int `json:"powerCapInWatts,omitempty"`
}

type RelativePowerCapInPercentageSpec struct {
	PowerCapPercentage int `json:"powerCapPercentage,omitempty"` // Power cap in percentage of peak power consumption
	SampleWindow       int `json:"sampleWindow,omitempty"`       // Sample window in seconds
}

// PowerCappingSpec specifies the kind of PowerCappingConfig
type PowerCappingSpec struct {
	Kind                             PowerCappingSpecKind `json:"kind,omitempty"`
	AbsolutePowerCapInWattsSpec      `json:"absolutePowerCapInWatts,omitempty"`
	RelativePowerCapInPercentageSpec `json:"relativePowerCapInPercentage,omitempty"`
}

type TemperatureThresholdKind string

const (
	NoTemperatureThreshold                                       TemperatureThresholdKind = "NoTemperatureThreshold"
	AbsoluteTemperatureThresholdInCelsius                        TemperatureThresholdKind = "AbsoluteTemperatureThresholdInCelsius"
	RelativeTemperatureThresholdOfPeakTemperatureInPercentage    TemperatureThresholdKind = "RelativeTemperatureThresholdOfPeakTemperatureInPercentage"
	RelativeTemperatureThresholdOfAverageTemperatureInPercentage TemperatureThresholdKind = "RelativeTemperatureThresholdOfAverageTemperatureInPercentage"
)

type AbsoluteTemperatureThresholdInCelsiusSpec struct {
	TemperatureThresholdInCelsius int `json:"temperatureThresholdInCelsius,omitempty"`
}

type RelativeTemperatureThresholdInPercentageSpec struct {
	TemperatureThresholdPercentage int `json:"temperatureThresholdPercentage,omitempty"` // Temperature threshold in percentage of peak temperature
	SampleWindow                   int `json:"sampleWindow,omitempty"`                   // Sample window in seconds
}

// TemperatureThresholdSpec specifies the kind of TemperatureThresholdConfig
type TemperatureThresholdSpec struct {
	Kind                                         TemperatureThresholdKind `json:"kind,omitempty"`
	AbsoluteTemperatureThresholdInCelsiusSpec    `json:"absoluteTemperatureThresholdInCelsius,omitempty"`
	RelativeTemperatureThresholdInPercentageSpec `json:"relativeTemperatureThresholdInPercentage,omitempty"`
}

// PowerCappingConfigSpec defines the desired state of PowerCappingConfig
type PowerCappingConfigSpec struct {
	WorkloadType             string                   `json:"workloadType,omitempty"`             // "training" or "inference"
	EfficiencyLevel          string                   `json:"efficiencyLevel,omitempty"`          // "low", "medium", "high"
	PowerCappingSpec         PowerCappingSpec         `json:"powerCappingSpec,omitempty"`         // Power capping specification
	TemperatureThresholdSpec TemperatureThresholdSpec `json:"temperatureThresholdSpec,omitempty"` // Temperature threshold specification
}

// PowerCappingConfigStatus is the status for a PowerCappingConfig resource
type PowerCappingConfigStatus struct {
	CurrentPowerConsumption  int `json:"currentPowerConsumption,omitempty"`
	ForecastPowerConsumption int `json:"forecastPowerConsumption,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PowerCappingConfig is the Schema for the powercappingconfigs API
type PowerCappingConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PowerCappingConfigSpec   `json:"spec,omitempty"`
	Status PowerCappingConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PowerCappingConfigList contains a list of PowerCappingConfig
type PowerCappingConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []PowerCappingConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PowerCappingConfig{}, &PowerCappingConfigList{})
}
