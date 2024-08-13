// pubsub.go
package alert

import (
	"sync"

	"github.com/Climatik-Project/Climatik-Project/api/v1alpha1"
)

type PubSub struct {
	subscribers map[string][]AlertManager
	mu          sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]AlertManager),
	}
}

func (ps *PubSub) Subscribe(topic string, subscriber AlertManager) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.subscribers[topic] = append(ps.subscribers[topic], subscriber)
}

func (ps *PubSub) Publish(topic string, podName string, powerCapValue int, devices map[string]string, config *v1alpha1.PowerCappingConfig) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	for _, subscriber := range ps.subscribers[topic] {
		go subscriber.CreateAlert(podName, powerCapValue, devices, config)
	}
}
