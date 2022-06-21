package message_store_db

import (
	"database/sql"
	"fmt"
	"github.com/glide-im/glide/pkg/messages"
	"github.com/glide-im/im-service/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

type ChatMessageStore struct {
	db *sql.DB
}

func New(conf *config.MySqlConf) (*ChatMessageStore, error) {
	mysqlUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.Db)
	db, err := sql.Open("mysql", mysqlUrl)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	m := &ChatMessageStore{
		db: db,
	}
	return m, nil
}

func (D *ChatMessageStore) StoreMessage(m *messages.ChatMessage) error {

	from, err := strconv.ParseInt(m.From, 10, 64)
	if err != nil {
		return err
	}
	to, err := strconv.ParseInt(m.To, 10, 64)
	if err != nil {
		return err
	}

	lg := from
	sm := to
	if lg < sm {
		lg, sm = sm, lg
	}
	sid := fmt.Sprintf("%d_%d", lg, sm)

	// todo update the type of user id to string
	//mysql only
	_, err = D.db.Exec(
		"INSERT INTO im_chat_message (`m_id`, `session_id`, `from`, `to`, `type`, `content`, `send_at`, `create_at`, `cli_seq`, `status`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)ON DUPLICATE KEY UPDATE send_at=?",
		m.Mid, sid, from, to, m.Type, m.Content, m.SendAt, time.Now().Unix(), 0, 0,
		m.SendAt)
	return err
}

type IdleChatMessageStore struct {
}

func (i *IdleChatMessageStore) StoreMessage(message *messages.ChatMessage) error {
	return nil
}
