package common

type ActionHandler interface {
	SetValue(value int)
	GetValue()
	StartChRo()
	Leave()
	Disconnect()
	Reconnect()
}
