// service.go
package alert

import "github.com/Climatik-Project/Climatik-Project/api/v1alpha1"

type AlertService struct {
	Pubsub *PubSub
}

func NewAlertService(config map[string]map[string]string) (*AlertService, error) {
	pubsub, err := CreateAlertService(config)
	if err != nil {
		return nil, err
	}
	return &AlertService{Pubsub: pubsub}, nil
}

func (s *AlertService) SendAlert(podName string, powerCapValue int, devices map[string]string, config *v1alpha1.PowerCappingConfig) error {
	s.Pubsub.Publish("alerts", podName, powerCapValue, devices, config)
	return nil
}
