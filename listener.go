package connutil

type Listener interface {
	OnConnected(id string)
	OnDisconnected(id string)
	OnMessage(id string, data []byte)
	OnReconnected(id string)
	OnRemove(id string)
}

type EmptyListener struct{}

func (e *EmptyListener) OnConnected(id string)            {}
func (e *EmptyListener) OnDisconnected(id string)         {}
func (e *EmptyListener) OnMessage(id string, data []byte) {}
func (e *EmptyListener) OnReconnected(id string)          {}
func (e *EmptyListener) OnRemove(id string)               {}

var DefaultListener = new(EmptyListener)

type listenerFunc struct {
	onConnected    func(id string)
	onReconnected  func(id string)
	onMessage      func(id string, data []byte)
	onDisconnected func(id string)
	onRemove       func(id string)
}

func (l *listenerFunc) OnConnected(id string)            { l.onConnected(id) }
func (l *listenerFunc) OnDisconnected(id string)         { l.onDisconnected(id) }
func (l *listenerFunc) OnMessage(id string, data []byte) { l.onMessage(id, data) }
func (l *listenerFunc) OnReconnected(id string)          { l.onReconnected(id) }
func (l *listenerFunc) OnRemove(id string)               { l.onRemove(id) }

func NewListener(
	onConnected func(id string),
	onDisconnected func(id string),
	onReconnected func(id string),
	onMessage func(id string, data []byte),
	onRemove func(id string),
) Listener {
	return &listenerFunc{
		onConnected:    onConnected,
		onDisconnected: onDisconnected,
		onReconnected:  onReconnected,
		onMessage:      onMessage,
		onRemove:       onRemove,
	}
}
