package fatplugin

type Reactor interface {
	// React checks if the message is relevant to this Reactor, and if so, returns true and executes an action
	React(chatID int, message string) bool
}
