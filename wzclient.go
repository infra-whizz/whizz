package whizz

import (
	"log"
	"time"

	"github.com/infra-whizz/wzlib"
	wzlib_transport "github.com/infra-whizz/wzlib/transport"
	"github.com/nats-io/nats.go"
)

type WzChannels struct {
	console *nats.Subscription
}

type WzClient struct {
	channels        *WzChannels
	transport       *wzlib_transport.WzdPubSub
	replies         []*wzlib_transport.WzGenericMessage
	expectedReplies int64
}

func NewWhizzClient() *WzClient {
	wzc := new(WzClient)
	wzc.transport = wzlib_transport.NewWizPubSub()
	wzc.channels = new(WzChannels)
	wzc.replies = make([]*wzlib_transport.WzGenericMessage, 0)
	return wzc
}

// Controller replied
func (wzc *WzClient) onControllerReplyEvent(msg *nats.Msg) {
	envelope := wzlib_transport.NewWzGenericMessage()
	if err := envelope.LoadBytes(msg.Data); err != nil {
		log.Println(err.Error())
	} else {
		batchMax, ok := envelope.Payload[wzlib_transport.PAYLOAD_BATCH_SIZE]
		if !ok || batchMax == nil {
			log.Println("Discarding controller reply: no batch.max defined")
		} else {
			wzc.expectedReplies = batchMax.(int64)
			wzc.replies = append(wzc.replies, envelope)
		}
	}
}

func (wzc *WzClient) initialise() {
	var err error
	wzc.transport.Start()
	wzc.channels.console, err = wzc.transport.
		GetSubscriber().
		Subscribe(wzlib.CHANNEL_CONTROLLER, wzc.onControllerReplyEvent)
	if err != nil {
		log.Panicf("Unable to subscribe to console channel: %s\n", err.Error())
	}
}

func (wzc *WzClient) start() {
}

func (wzc *WzClient) Boot() {
	wzc.initialise()
	wzc.start()
}

func (wzc *WzClient) Stop() {
	wzc.transport.Disconnect()
}

// Wait for the client get reply from the cluster with timeout.
func (wzc *WzClient) Wait(sec int) {
	if sec < 1 {
		sec = 30
	}

	cs := 0
	for {
		if cs >= sec {
			return
		}
		ms := 0
		for {
			time.Sleep(time.Millisecond)
			ms++
			if ms > 0x400 {
				break
			}
			if len(wzc.replies) == int(wzc.expectedReplies) && wzc.expectedReplies > 0 {
				return
			}
		}
		log.Println("Waiting", sec-cs, "seconds")
		cs++
	}
}
