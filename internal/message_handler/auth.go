package message_handler

import (
	"errors"
	"github.com/glide-im/glide/pkg/auth"
	"github.com/glide-im/glide/pkg/auth/jwt_auth"
	"github.com/glide-im/glide/pkg/gate"
	"github.com/glide-im/glide/pkg/gate/gateway"
	"github.com/glide-im/glide/pkg/logger"
	"github.com/glide-im/glide/pkg/messages"
)

type AuthRequest struct {
}

func (d *MessageHandler) handleAuth(c *gate.Info, msg *messages.GlideMessage) error {

	t := auth.Token{}
	e := msg.Data.Deserialize(&t)
	if e != nil {
		resp := messages.NewMessage(0, ActionApiFailed, "invalid token")
		d.enqueueMessage(c.ID, resp)
		return nil
	}

	info := jwt_auth.JwtAuthInfo{
		UID:    msg.From,
		Device: c.ID.Device(),
	}
	r, err := d.auth.Auth(&info, &t)

	if err != nil {
		return err
	}

	if r.Success {
		respMsg := messages.NewMessage(msg.Seq, ActionApiSuccess, r.Response)
		jwtResp, ok := r.Response.(*jwt_auth.Response)
		if !ok {
			resp := messages.NewMessage(msg.Seq, ActionApiFailed, "internal error")
			d.enqueueMessage(c.ID, resp)
			return errors.New("invalid response type: expected *jwt_auth.Response")
		}

		newID := gate.NewID("", jwtResp.Uid, jwtResp.Device)
		err = d.def.GetClientInterface().SetClientID(c.ID, newID)
		if gateway.IsIDAlreadyExist(err) {
			tempId, err := gate.GenTempID(newID.Gateway())
			if err != nil {
				return err
			}
			err = d.def.GetClientInterface().SetClientID(newID, tempId)
			if err != nil {
				return errors.New("failed to set temp id:" + err.Error())
			}
			d.enqueueMessage(tempId, createKickOutMessage(c))
			err = d.def.GetClientInterface().SetClientID(c.ID, newID)
			if err != nil {
				return errors.New("failed to set new id:" + err.Error())
			}
		} else if gateway.IsClientNotExist(err) {
			return errors.New("auth client not exist")
		} else if err != nil {
			return err
		}

		logger.D("auth success: %s", newID)
		d.enqueueMessage(newID, respMsg)
	} else {
		resp := messages.NewMessage(msg.Seq, ActionApiFailed, r.Msg)
		d.enqueueMessage(c.ID, resp)
	}
	return nil
}
