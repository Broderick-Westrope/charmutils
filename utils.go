package btutils

import (
	"fmt"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
)

func UpdateTypedModel[T tea.Model](model *T, msg tea.Msg) (tea.Cmd, error) {
	var ok bool
	newModel, cmd := (*model).Update(msg)
	*model, ok = newModel.(T)
	if !ok {
		return nil, fmt.Errorf("failed to update model of type %q: %w", reflect.TypeOf(model), ErrInvalidTypeAssertion)
	}
	return cmd, nil
}
