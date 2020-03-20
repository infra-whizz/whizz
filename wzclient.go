package whizz

import (
	"log"

	"github.com/infra-whizz/wzlib"
	wzlib_transport "github.com/infra-whizz/wzlib/transport"
	"github.com/nats-io/nats.go"
)

type WzChannels struct {
	console *nats.Subscription
}

type WzClient struct {
	channels  *WzChannels
	transport *wzlib_transport.WzdPubSub
}

func NewWhizzClient() *WzClient {
	wzc := new(WzClient)
	wzc.transport = wzlib_transport.NewWizPubSub()
	wzc.channels = new(WzChannels)
	return wzc
}

func (wzc *WzClient) onControllerReplyEvent(msg *nats.Msg) {
	log.Println("received from console channel:", len(msg.Data), "bytes")
	envelope := wzlib.NewWzConsoleReplyMessage()
	envelope.LoadBytes(msg.Data)
	log.Println("Jid:", envelope.Jid)
}

func (wzc *WzClient) initialise() {
	var err error
	wzc.transport.Start()
	wzc.channels.console, err = wzc.transport.
		GetSubscriber().
		Subscribe(wzlib.CHANNEL_CONSOLE, wzc.onControllerReplyEvent)
	if err != nil {
		log.Panicf("Unable to subscribe to console channel: %s\n", err.Error())
	}
}

func (wzc *WzClient) start() {
	envelope := wzlib.NewWzConsoleMessage()
	envelope.Jid = "some-shit-jid"
	out, err := envelope.Serialise()
	if err != nil {
		log.Println("Error serialising:", err.Error())
	} else {
		wzc.transport.GetPublisher().Publish(wzlib.CHANNEL_CONSOLE, out)
	}
}

func (wzc *WzClient) stop() {
	wzc.transport.Disconnect()
}

func (wzc *WzClient) Call() {
	wzc.initialise()
	wzc.start()
	wzc.stop()
}
