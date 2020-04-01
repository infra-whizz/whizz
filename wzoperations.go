package whizz

import (
	"log"

	"github.com/davecgh/go-spew/spew"

	"github.com/infra-whizz/wzlib"
	wzlib_transport "github.com/infra-whizz/wzlib/transport"
)

// Send message on behalf of the console
func (wzc *WzClient) send(evp *wzlib_transport.WzGenericMessage) {
	out, err := evp.Serialise()
	if err != nil {
		log.Println("Error serialising envelope:", err.Error())
	} else {
		if err := wzc.transport.GetPublisher().Publish(wzlib.CHANNEL_CONSOLE, out); err != nil {
			log.Println("Error sending message:", err.Error())
		}
	}
}

func (wzc *WzClient) Call() {
}

func (wzc *WzClient) Accept(fingerprints ...string) {
}

func (wzc *WzClient) Reject(fingerprints ...string) {
}

func (wzc *WzClient) ListNew() {
	envelope := wzlib_transport.NewWzMessage(wzlib_transport.MSGTYPE_CLIENT)
	envelope.Payload[wzlib_transport.PAYLOAD_COMMAND] = "list.clients.new"
	wzc.send(envelope)
	wzc.Wait(5)
	log.Println("-- Begin reading")
	for _, m := range wzc.replies {
		log.Println(".........")
		spew.Dump(m)
	}
	log.Println("-- End reading")
}

func (wzc *WzClient) ListRejected() {
	envelope := wzlib_transport.NewWzMessage(wzlib_transport.MSGTYPE_CLIENT)
	envelope.Payload[wzlib_transport.PAYLOAD_COMMAND] = "list.clients.rejected"
	wzc.send(envelope)
	wzc.Wait(5)
}
