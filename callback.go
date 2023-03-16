package connutil

type Callback interface {
	OnMessage(string, []byte)
	OnConnected(string)
	OnDisconnected(string)
	OnShutdown(string)
}

type EmptyCallback struct{}

func (c *EmptyCallback) OnMessage(string, []byte) {}
func (c *EmptyCallback) OnConnected(string)       {}
func (c *EmptyCallback) OnDisconnected(string)    {}
func (c *EmptyCallback) OnShutdown(string)        {}
