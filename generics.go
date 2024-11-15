package charmutils

import (
	"fmt"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
)

// UpdateTypedModel performs an update on the model using the given msg.
// This enables easily storing models of a concrete type without having the clutter of frequent type assertions.
func UpdateTypedModel[T tea.Model](model *T, msg tea.Msg) (tea.Cmd, error) {
	var ok bool
	newModel, cmd := (*model).Update(msg)
	*model, ok = newModel.(T)
	if !ok {
		return nil, fmt.Errorf("failed to update model of type %q: %w", reflect.TypeOf(model), ErrInvalidTypeAssertion)
	}
	return cmd, nil
}
