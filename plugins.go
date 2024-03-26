package godantic

type ValidationPlugin interface {
	Validate() *CustomErr
}
