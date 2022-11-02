package connutil

type Listener interface {
	OnConnected(c *Client)
	OnDisconnected(c *Client)
	OnMessage(c *Client, data []byte)
	OnReconnected(c *Client)
	OnRemove(c *Client)
}

type EmptyListener struct{}

func (e *EmptyListener) OnConnected(c *Client)            {}
func (e *EmptyListener) OnDisconnected(c *Client)         {}
func (e *EmptyListener) OnMessage(c *Client, data []byte) {}
func (e *EmptyListener) OnReconnected(c *Client)          {}
func (e *EmptyListener) OnRemove(c *Client)               {}

var DefaultListener = new(EmptyListener)

type listenerFunc struct {
	onConnected    func(c *Client)
	onReconnected  func(c *Client)
	onMessage      func(c *Client, data []byte)
	onDisconnected func(c *Client)
	onRemove       func(c *Client)
}

func (l *listenerFunc) OnConnected(c *Client)            { l.onConnected(c) }
func (l *listenerFunc) OnDisconnected(c *Client)         { l.onDisconnected(c) }
func (l *listenerFunc) OnMessage(c *Client, data []byte) { l.onMessage(c, data) }
func (l *listenerFunc) OnReconnected(c *Client)          { l.onReconnected(c) }
func (l *listenerFunc) OnRemove(c *Client)               { l.onRemove(c) }

func NewListener(
	onConnected func(c *Client),
	onDisconnected func(c *Client),
	onReconnected func(c *Client),
	onMessage func(c *Client, data []byte),
	onRemove func(c *Client),
) Listener {
	return &listenerFunc{
		onConnected:    onConnected,
		onDisconnected: onDisconnected,
		onReconnected:  onReconnected,
		onMessage:      onMessage,
		onRemove:       onRemove,
	}
}
