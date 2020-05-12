package whizz

import (
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_transport "github.com/infra-whizz/wzlib/transport"
	"github.com/nats-io/nats.go"
)

type WzEvents struct {
	expectedReplies int64
	replies         []*wzlib_transport.WzGenericMessage
	wzlib_logger.WzLogger
}

func NewWzEvents() *WzEvents {
	wze := new(WzEvents)
	wze.replies = make([]*wzlib_transport.WzGenericMessage, 0)

	return wze
}

// XXX: this should be a part of transport!
func (wze *WzEvents) getMessageError(envelope *wzlib_transport.WzGenericMessage) string {
	fRetBlock, ex := envelope.Payload[wzlib_transport.PAYLOAD_FUNC_RET]
	if ex {
		err, ex := fRetBlock.(map[string]interface{})["error"]
		if ex {
			return err.(string)
		}
	}
	return ""
}

// onControllerReplyEvent function is fired when a controller has been replied
func (wze *WzEvents) onControllerReplyEvent(msg *nats.Msg) {
	envelope := wzlib_transport.NewWzGenericMessage()
	if err := envelope.LoadBytes(msg.Data); err != nil {
		wze.GetLogger().Errorln(err.Error())
	} else {
		if errmsg := wze.getMessageError(envelope); errmsg != "" {
			wze.GetLogger().Errorf("Error: %s", errmsg)
			return
		}
		batchMax, ok := envelope.Payload[wzlib_transport.PAYLOAD_BATCH_SIZE]
		if !ok || batchMax == nil {
			wze.GetLogger().Warningln("Discarding controller reply: no batch.max defined")
		} else {
			wze.expectedReplies = batchMax.(int64)
			wze.replies = append(wze.replies, envelope)
		}
	}
}
