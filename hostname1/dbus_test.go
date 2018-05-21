package hostname1

import (
	"testing"
)

// TestNew ensures that New() works without errors.
func TestNew(t *testing.T) {
	_, err := New()

	if err != nil {
		t.Fatal(err)
	}
}

func TestProperties(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatal(err)
	}

	properties, err := c.GetProperties()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("GetProperties: %#v", properties)
}
