package message_store_db

import (
	"github.com/glide-im/glide/internal/config"
	"github.com/glide-im/glide/pkg/messages"
	"testing"
)

func TestNew2(t *testing.T) {

	store, err := New(&config.MySqlConf{
		Host:     "dengzii.com",
		Port:     3306,
		Username: "go_im_test",
		Password: "N5zJxcFKMXLzWtAc",
		Db:       "go_im_test",
	})
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	err = store.StoreMessage("", &messages.ChatMessage{
		Mid:     1,
		Seq:     0,
		From:    0,
		To:      0,
		Type:    0,
		Content: "-",
		SendAt:  0,
	})
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}
