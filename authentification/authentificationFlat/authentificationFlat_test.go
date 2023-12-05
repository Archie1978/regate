package authentificationFlat

import (
	"fmt"
	"testing"
	"time"
)

// Test get code and application not started
func TestCodeAnswer(t *testing.T) {
	fmt.Println("Get Code without server")
	var authentificationFlat AuthentificationFlat
	code, err := authentificationFlat.GetCode()
	if err != ErrAppNotStarted {
		t.Fatal("Test Get code wthout server error:", err)
	}
	if code != "" {
		t.Fatal("Code must be empty")
	}

	// Generate code
	go func() {
		authentificationFlat.engineGenerateCode()
	}()

	for {
		<-time.After(time.Second)
		code, err = authentificationFlat.GetCode()
		fmt.Println("code, err", code, err)
	}
}
