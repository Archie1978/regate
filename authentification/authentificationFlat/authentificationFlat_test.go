package authentificationFlat

import (
	"fmt"
	"testing"
	"time"
)

// Test get code and application not started
func TestCodeAnswer(t *testing.T) {
	fmt.Println("Get Code without server")
	var authentficationFlat AuthentificationFlat
	code, err := authentficationFlat.GetCode()
	if err != ErrAppNotStarted {
		t.Fatal("Test Get code wthout server error:", err)
	}
	if code != "" {
		t.Fatal("Code must be empty")
	}

	// Generate code
	go func() {
		authentficationFlat.engineGenerateCode()
	}()

	for {
		<-time.After(time.Second)
		code, err = authentficationFlat.GetCode()
		fmt.Println("code, err", code, err)
	}
}
