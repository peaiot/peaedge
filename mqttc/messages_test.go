package mqttc

import (
	"testing"
)

func TestDecodeMessage(t *testing.T) {
	messageJson, _ := newBootstrapMessage("test").Encode()
	t.Log(messageJson)
	message, err := DecodeMessage[BootstrapMessage](messageJson)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(message)
}
