package handlers

// AlertHandler interface defines the method that all alert handlers must implement
type AlertHandler interface {
	HandleAlert(payload []byte) error
}
