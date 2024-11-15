package charmutils

import (
	"github.com/charmbracelet/bubbles/key"
)

// ConstructKeyBinding creates a new binding with the given keys and a help message with the
// description and each of the keys.
func ConstructKeyBinding(keys []string, desc string) key.Binding {
	var helpKeys string
	for _, k := range keys {
		if k == " " {
			k = "space"
		}
		helpKeys += k + ", "
	}
	helpKeys = helpKeys[:len(helpKeys)-2]

	return key.NewBinding(key.WithKeys(keys...), key.WithHelp(helpKeys, desc))
}
