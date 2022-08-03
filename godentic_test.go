package godantic

import "testing"
import "github.com/stretchr/testify/assert"

func TestExtraFields(t *testing.T) {
	structure := map[string]interface{}{
		"name":    "ismael",
		"surname": "GraHms",
	}

	body := map[string]interface{}{
		"name":    "ismael",
		"surname": "GraHms",
		"foo":     "bar",
	}
	g := New()
	err := g.ForbidExtraFields(body, structure, "")
	assert.Equal(t, err.Error(), "foo")
}
