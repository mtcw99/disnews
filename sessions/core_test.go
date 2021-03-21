package sessions

import (
	"testing"

	"github.com/google/uuid"
)

func TestSessionNewValid(t *testing.T) {
	sessions := New()
	uuid, err := sessions.NewSession("user")
	if err != nil {
		t.Fatal("Error has occured: ", err)
	}

	valid := sessions.ValidateSession(uuid)
	if !valid {
		t.Fatal("Invalid session when it should be valid")
	}
}

func TestSessionInvalid(t *testing.T) {
	sessions := New()
	valid := sessions.ValidateSession(uuid.NewString())
	if valid {
		t.Fatal("Valid session when it should be invalid")
	}
}
