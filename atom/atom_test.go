package atom_test

import (
	"testing"

	"github.com/stretchr/pat/atom"
	"github.com/stretchr/testify/assert"
)

func TestAtom(t *testing.T) {

	obj1 := map[string]interface{}{}
	obj2 := map[string]interface{}{}
	a := atom.New(obj1)

	assert.Equal(t, obj1, a.Value())
	a.Set(obj2)
	assert.Equal(t, obj2, a.Value())

}
