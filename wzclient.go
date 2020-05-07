package whizz

import (
	"time"

	wzlib_crypto "github.com/infra-whizz/wzlib/crypto"

	"github.com/infra-whizz/wzlib"
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_transport "github.com/infra-whizz/wzlib/transport"
	"github.com/nats-io/nats.go"
)

type WzChannels struct {
	console *nats.Subscription
}

type WzClient struct {
	channels     *WzChannels
	events       *WzEvents
	transport    *wzlib_transport.WzdPubSub
	cryptobundle *wzlib_crypto.WzCryptoBundle
	wzlib_logger.WzLogger
}

func NewWhizzClient() *WzClient {
	wzc := new(WzClient)
	wzc.transport = wzlib_transport.NewWizPubSub()
	wzc.channels = new(WzChannels)
	wzc.events = NewWzEvents()
	wzc.cryptobundle = wzlib_crypto.NewWzCryptoBundle()

	return wzc
}

// GetCryptoBundle returns a cryptobundle with AES, RSA and utils API
func (wzc *WzClient) GetCryptoBundle() *wzlib_crypto.WzCryptoBundle {
	return wzc.cryptobundle
}

func (wzc *WzClient) initialise() {
	var err error
	wzc.transport.Start()
	wzc.channels.console, err = wzc.transport.
		GetSubscriber().
		Subscribe(wzlib.CHANNEL_CONTROLLER, wzc.events.onControllerReplyEvent)
	if err != nil {
		wzc.GetLogger().Panicf("Unable to subscribe to console channel: %s\n", err.Error())
	}
}

func (wzc *WzClient) start() {
}

// Boot the whizz client
func (wzc *WzClient) Boot() {
	wzc.initialise()
	wzc.start()
}

// Stop and teardown everything: disconnect, cleanup etc.
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
			if len(wzc.events.replies) == int(wzc.events.expectedReplies) && wzc.events.expectedReplies > 0 {
				return
			}
		}
		wzc.GetLogger().Debugln("Waiting", sec-cs, "seconds")
		cs++
	}
}
