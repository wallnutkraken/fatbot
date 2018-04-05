package fatplugin

type Cleaner interface {
	// Clean cleans the given text and returns it
	Clean(text string) string
}
