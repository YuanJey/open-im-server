package cache

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestParallelSetMessageToCache(t *testing.T) {
	var (
		cid      = fmt.Sprintf("cid-%v", rand.Int63())
		seqFirst = rand.Int63()
		msgs     = []*sdkws.MsgData{}
	)

	for i := 0; i < 100; i++ {
		msgs = append(msgs, &sdkws.MsgData{
			Seq: seqFirst + int64(i),
		})
	}

	testParallelSetMessageToCache(t, cid, msgs)
}

func testParallelSetMessageToCache(t *testing.T, cid string, msgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	ret, err := cacher.ParallelSetMessageToCache(context.Background(), cid, msgs)
	assert.Nil(t, err)
	assert.Equal(t, len(msgs), ret)

	// validate
	for _, msg := range msgs {
		key := cacher.getMessageCacheKey(cid, msg.Seq)
		val, err := rdb.Exists(context.Background(), key).Result()
		assert.Nil(t, err)
		assert.EqualValues(t, 1, val)
	}
}

func TestPipeSetMessageToCache(t *testing.T) {
	var (
		cid      = fmt.Sprintf("cid-%v", rand.Int63())
		seqFirst = rand.Int63()
		msgs     = []*sdkws.MsgData{}
	)

	for i := 0; i < 100; i++ {
		msgs = append(msgs, &sdkws.MsgData{
			Seq: seqFirst + int64(i),
		})
	}

	testPipeSetMessageToCache(t, cid, msgs)
}

func testPipeSetMessageToCache(t *testing.T, cid string, msgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	ret, err := cacher.PipeSetMessageToCache(context.Background(), cid, msgs)
	assert.Nil(t, err)
	assert.Equal(t, len(msgs), ret)

	// validate
	for _, msg := range msgs {
		key := cacher.getMessageCacheKey(cid, msg.Seq)
		val, err := rdb.Exists(context.Background(), key).Result()
		assert.Nil(t, err)
		assert.EqualValues(t, 1, val)
	}
}

func TestGetMessagesBySeq(t *testing.T) {
	var (
		cid      = fmt.Sprintf("cid-%v", rand.Int63())
		seqFirst = rand.Int63()
		msgs     = []*sdkws.MsgData{}
	)

	seqs := []int64{}
	for i := 0; i < 100; i++ {
		msgs = append(msgs, &sdkws.MsgData{
			Seq:    seqFirst + int64(i),
			SendID: fmt.Sprintf("fake-sendid-%v", i),
		})
		seqs = append(seqs, seqFirst+int64(i))
	}

	// set data to cache
	testPipeSetMessageToCache(t, cid, msgs)

	// get data from cache with parallet mode
	testParallelGetMessagesBySeq(t, cid, seqs, msgs)

	// get data from cache with pipeline mode
	testPipeGetMessagesBySeq(t, cid, seqs, msgs)
}

func testParallelGetMessagesBySeq(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	respMsgs, failedSeqs, err := cacher.ParallelGetMessagesBySeq(context.Background(), cid, seqs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(failedSeqs))
	assert.Equal(t, len(respMsgs), len(seqs))

	// validate
	for idx, msg := range respMsgs {
		assert.Equal(t, msg.Seq, inputMsgs[idx].Seq)
		assert.Equal(t, msg.SendID, inputMsgs[idx].SendID)
	}
}

func testPipeGetMessagesBySeq(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	respMsgs, failedSeqs, err := cacher.PipeGetMessagesBySeq(context.Background(), cid, seqs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(failedSeqs))
	assert.Equal(t, len(respMsgs), len(seqs))

	// validate
	for idx, msg := range respMsgs {
		assert.Equal(t, msg.Seq, inputMsgs[idx].Seq)
		assert.Equal(t, msg.SendID, inputMsgs[idx].SendID)
	}
}

func TestGetMessagesBySeqWithEmptySeqs(t *testing.T) {
	var (
		cid            = fmt.Sprintf("cid-%v", rand.Int63())
		seqFirst int64 = 0
		msgs           = []*sdkws.MsgData{}
	)

	seqs := []int64{}
	for i := 0; i < 100; i++ {
		msgs = append(msgs, &sdkws.MsgData{
			Seq:    seqFirst + int64(i),
			SendID: fmt.Sprintf("fake-sendid-%v", i),
		})
		seqs = append(seqs, seqFirst+int64(i))
	}

	// don't set cache, only get data from cache.

	// get data from cache with parallet mode
	testParallelGetMessagesBySeqWithEmptry(t, cid, seqs, msgs)

	// get data from cache with pipeline mode
	testPipeGetMessagesBySeqWithEmptry(t, cid, seqs, msgs)
}

func testParallelGetMessagesBySeqWithEmptry(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	respMsgs, failedSeqs, err := cacher.ParallelGetMessagesBySeq(context.Background(), cid, seqs)
	assert.Nil(t, err)
	assert.Equal(t, len(seqs), len(failedSeqs))
	assert.Equal(t, 0, len(respMsgs))
}

func testPipeGetMessagesBySeqWithEmptry(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	respMsgs, failedSeqs, err := cacher.PipeGetMessagesBySeq(context.Background(), cid, seqs)
	assert.Equal(t, err, redis.Nil)
	assert.Equal(t, len(seqs), len(failedSeqs))
	assert.Equal(t, 0, len(respMsgs))
}

func TestGetMessagesBySeqWithLostHalfSeqs(t *testing.T) {
	var (
		cid            = fmt.Sprintf("cid-%v", rand.Int63())
		seqFirst int64 = 0
		msgs           = []*sdkws.MsgData{}
	)

	seqs := []int64{}
	for i := 0; i < 100; i++ {
		msgs = append(msgs, &sdkws.MsgData{
			Seq:    seqFirst + int64(i),
			SendID: fmt.Sprintf("fake-sendid-%v", i),
		})
		seqs = append(seqs, seqFirst+int64(i))
	}

	// Only set half the number of messages.
	testParallelSetMessageToCache(t, cid, msgs[:50])

	// get data from cache with parallet mode
	testParallelGetMessagesBySeqWithLostHalfSeqs(t, cid, seqs, msgs)

	// get data from cache with pipeline mode
	testPipeGetMessagesBySeqWithLostHalfSeqs(t, cid, seqs, msgs)
}

func testParallelGetMessagesBySeqWithLostHalfSeqs(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	respMsgs, failedSeqs, err := cacher.ParallelGetMessagesBySeq(context.Background(), cid, seqs)
	assert.Nil(t, err)
	assert.Equal(t, len(seqs)/2, len(failedSeqs))
	assert.Equal(t, len(seqs)/2, len(respMsgs))

	for idx, msg := range respMsgs {
		assert.Equal(t, msg.Seq, seqs[idx])
	}
}

func testPipeGetMessagesBySeqWithLostHalfSeqs(t *testing.T, cid string, seqs []int64, inputMsgs []*sdkws.MsgData) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	cacher := msgCache{rdb: rdb}

	respMsgs, failedSeqs, err := cacher.PipeGetMessagesBySeq(context.Background(), cid, seqs)
	assert.Nil(t, err)
	assert.Equal(t, len(seqs)/2, len(failedSeqs))
	assert.Equal(t, len(seqs)/2, len(respMsgs))

	for idx, msg := range respMsgs {
		assert.Equal(t, msg.Seq, seqs[idx])
	}
}
