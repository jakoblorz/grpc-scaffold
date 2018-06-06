package endpoint

type ControllerWrapper struct {
	addr        string
	controllers []Controller
}

func NewControllerWrapper(addr string) ControllerWrapper {
	return ControllerWrapper{
		addr:        addr,
		controllers: make([]Controller, 0),
	}
}

func (c ControllerWrapper) serverAddr() string {
	return c.addr
}

func (c ControllerWrapper) registerController(s Controller) {

}
