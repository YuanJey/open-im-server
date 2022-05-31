package db

import (
	pbChat "Open_IM/pkg/proto/chat"
	"context"
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SetTokenMapByUidPid(t *testing.T) {
	m := make(map[string]int, 0)
	m["test1"] = 1
	m["test2"] = 2
	m["2332"] = 4
	_ = DB.SetTokenMapByUidPid("1234", 2, m)

}
func Test_GetTokenMapByUidPid(t *testing.T) {
	m, err := DB.GetTokenMapByUidPid("1234", "Android")
	assert.Nil(t, err)
	fmt.Println(m)
}

func TestDataBases_GetMultiConversationMsgOpt(t *testing.T) {
	m, err := DB.GetMultiConversationMsgOpt("fg", []string{"user", "age", "color"})
	assert.Nil(t, err)
	fmt.Println(m)
}
func Test_GetKeyTTL(t *testing.T) {
	ctx := context.Background()
	key := flag.String("key", "key", "key value")
	flag.Parse()
	ttl, err := DB.rdb.TTL(ctx, *key).Result()
	assert.Nil(t, err)
	fmt.Println(ttl)
}
func Test_HGetAll(t *testing.T) {
	ctx := context.Background()
	key := flag.String("key", "key", "key value")
	flag.Parse()
	ttl, err := DB.rdb.TTL(ctx, *key).Result()
	assert.Nil(t, err)
	fmt.Println(ttl)
}
func Test_NewSetMessageToCache(t *testing.T) {
	var msg pbChat.MsgDataToMQ
	uid := "test_uid"
	msg.MsgData.Seq = 11
	msg.MsgData.ClientMsgID = "23jwhjsdf"
	messageList := []*pbChat.MsgDataToMQ{&msg}
	err := DB.NewSetMessageToCache(messageList, uid, "test")
	assert.Nil(t, err)

}
