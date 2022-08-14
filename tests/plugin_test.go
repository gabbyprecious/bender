package plugin

import (
	"os"
	"testing"

	"github.com/vincenzopalazzo/cln4go/client"
)

func TestCallFistMethod(t *testing.T) {
	path := os.Getenv("CLN_UNIX_SOCKET")
	client, err := client.NewUnix(path)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	password_response, err := client.Call("bender_set_password", map[string]any{"password": "alibaba"})
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	_, found := password_response["message"]
	if !found {
		t.Error("The message is not found")
	}

	response, err := client.Call("bender_run_server", make(map[string]interface{}))
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	message, found := response["message"]
	if !found {
		t.Error("The message is not found")
	}

	if message != "Server up and running,/ listen and serve on 0.0.0.0:9080" {
		t.Errorf("message received %s different from expected %s", message, "Server up and running,/ listen and serve on 0.0.0.0:8080")
	}
}
