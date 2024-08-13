package alert

import (
	v1alpha1 "github.com/Climatik-Project/Climatik-Project/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewMockPowerCappingConfig returns a new instance of a mock PowerCappingConfig
func NewMockPowerCappingConfig() *v1alpha1.PowerCappingConfig {
	return &v1alpha1.PowerCappingConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-powercapping-config",
			Namespace: "default",
		},
		Spec: v1alpha1.PowerCappingConfigSpec{
			WorkloadType:    "training",
			EfficiencyLevel: "High",
			PowerCappingSpec: v1alpha1.PowerCappingSpec{
				Kind: v1alpha1.RelativePowerCapOfPeakPowerConsumptionInPercentage,
				RelativePowerCapInPercentageSpec: v1alpha1.RelativePowerCapInPercentageSpec{
					PowerCapPercentage: 80,
					SampleWindow:       30,
				},
			},
		},
		Status: v1alpha1.PowerCappingConfigStatus{
			CurrentPowerConsumption:  50,
			ForecastPowerConsumption: 75,
		},
	}
}
