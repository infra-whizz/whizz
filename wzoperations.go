package whizz

import (
	"fmt"

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

// Request a batch of replies to a specific command
func (wzc *WzClient) getRepliesOnCommand(command string) {
	envelope := wzlib_transport.NewWzMessage(wzlib_transport.MSGTYPE_CLIENT)
	envelope.Payload[wzlib_transport.PAYLOAD_COMMAND] = command
	wzc.send(envelope)
	wzc.Wait(5)
}

// Join all the batch chunks on a return point
func (wzc *WzClient) gatherChunksOn(returnpoint string, clients *[]map[string]interface{}) {
	for _, msg := range wzc.replies {
		payload := msg.Payload[wzlib_transport.PAYLOAD_FUNC_RET].(map[string]interface{})[returnpoint]
		if payload != nil {
			for _, client := range payload.([]interface{}) {
				*clients = append(*clients, client.(map[string]interface{}))
			}
		}
	}
}

// Call query to target registered client (main mode)
func (wzc *WzClient) Call() {
}

// Accept unaccepted (new) clients. If no fingerprints, all accepted
func (wzc *WzClient) Accept(fingerprints ...string) (missing []string) {
	var msg string
	if len(fingerprints) > 0 {
		msg = fmt.Sprintf("%d client machines", len(fingerprints))
	} else {
		msg = "all new client machines"
		fingerprints = make([]string, 0)
	}
	wzc.GetLogger().Debugf("Accepting %s", msg)

	envelope := wzlib_transport.NewWzMessage(wzlib_transport.MSGTYPE_CLIENT)
	envelope.Payload[wzlib_transport.PAYLOAD_COMMAND] = "clients.accept"
	envelope.Payload[wzlib_transport.PAYLOAD_COMMAND_PARAMS] = map[string]interface{}{"fingerprints": fingerprints}

	wzc.send(envelope)
	wzc.Wait(5)

	missing = make([]string, 0)
	for _, msg := range wzc.replies {
		payload := msg.Payload[wzlib_transport.PAYLOAD_FUNC_RET].(map[string]interface{})["accepted.missing"]
		if payload != nil {
			for _, missingFp := range payload.([]interface{}) {
				missing = append(missing, missingFp.(string))
			}
		}
	}
	return
}

// Reject unaccepted (new) clients
func (wzc *WzClient) Reject(fingerprints ...string) (missing []string) {
	var msg string
	if len(fingerprints) > 0 {
		msg = fmt.Sprintf("%d client machines", len(fingerprints))
	} else {
		msg = "all new client machines"
	}
	wzc.GetLogger().Debugf("Rejecting %s", msg)

	envelope := wzlib_transport.NewWzMessage(wzlib_transport.MSGTYPE_CLIENT)
	envelope.Payload[wzlib_transport.PAYLOAD_COMMAND] = "clients.reject"
	envelope.Payload[wzlib_transport.PAYLOAD_COMMAND_PARAMS] = map[string]interface{}{"fingerprints": fingerprints}

	wzc.send(envelope)
	wzc.Wait(5)

	missing = make([]string, 0)
	for _, msg := range wzc.replies {
		payload := msg.Payload[wzlib_transport.PAYLOAD_FUNC_RET].(map[string]interface{})["rejected.missing"]
		if payload != nil {
			for _, missingFp := range payload.([]interface{}) {
				missing = append(missing, missingFp.(string))
			}
		}
	}
	return

}

// Delete any clients regardless what their status is, by fingerprints
func (wzc *WzClient) Delete(fingerprints ...string) {
	wzc.GetLogger().Debugf("Deleting %d client machines", len(fingerprints))
}

// ListNew clients returning info about each client.
func (wzc *WzClient) ListNew() []map[string]interface{} {
	clients := make([]map[string]interface{}, 0)
	wzc.getRepliesOnCommand("list.clients.new")
	wzc.gatherChunksOn("registered", &clients)
	return clients
}

// ListRejected clients
func (wzc *WzClient) ListRejected() []map[string]interface{} {
	clients := make([]map[string]interface{}, 0)
	wzc.getRepliesOnCommand("list.clients.rejected")
	wzc.gatherChunksOn("rejected", &clients)
	return clients
}
