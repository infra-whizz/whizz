package whizz

import (
	"github.com/infra-whizz/wzlib"
	wzlib_transport "github.com/infra-whizz/wzlib/transport"
)

// Send message on behalf of the console
func (wzc *WzClient) send(evp *wzlib_transport.WzGenericMessage) {
	out, err := evp.Serialise()
	if err != nil {
		wzc.GetLogger().Errorln("Error serialising envelope:", err.Error())
	} else {
		if err := wzc.transport.GetPublisher().Publish(wzlib.CHANNEL_CONSOLE, out); err != nil {
			wzc.GetLogger().Errorln("Error sending message:", err.Error())
		}
	}
}

func (wzc *WzClient) Call() {
}

func (wzc *WzClient) Accept(fingerprints ...string) {
}

func (wzc *WzClient) Reject(fingerprints ...string) {
}

// ListNew clients returning info about each client.
func (wzc *WzClient) ListNew() []map[string]interface{} {
	envelope := wzlib_transport.NewWzMessage(wzlib_transport.MSGTYPE_CLIENT)
	envelope.Payload[wzlib_transport.PAYLOAD_COMMAND] = "list.clients.new"
	wzc.send(envelope)
	wzc.Wait(5)

	clients := make([]map[string]interface{}, 0)
	for _, msg := range wzc.replies {
		payload := msg.Payload[wzlib_transport.PAYLOAD_FUNC_RET].(map[string]interface{})["registered"]
		if payload != nil {
			for _, client := range payload.([]interface{}) {
				clients = append(clients, client.(map[string]interface{}))
			}
		}
	}
	return clients
}

func (wzc *WzClient) ListRejected() {
	envelope := wzlib_transport.NewWzMessage(wzlib_transport.MSGTYPE_CLIENT)
	envelope.Payload[wzlib_transport.PAYLOAD_COMMAND] = "list.clients.rejected"
	wzc.send(envelope)
	wzc.Wait(5)
}
