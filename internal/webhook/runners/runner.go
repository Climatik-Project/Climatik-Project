package runners

// Runner interface defines the method that all runners must implement
type Runner interface {
	Run() error
}
