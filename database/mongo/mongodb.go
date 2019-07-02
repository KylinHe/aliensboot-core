package mongo

import (
	"container/heap"
	"github.com/KylinHe/aliensboot-core/common/util"
	"github.com/KylinHe/aliensboot-core/config"
	"github.com/KylinHe/aliensboot-core/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

// session
type Session struct {
	*mgo.Session
	ref   int
	index int
}

// session heap
type SessionHeap []*Session

func (h SessionHeap) Len() int {
	return len(h)
}

func (h SessionHeap) Less(i, j int) bool {
	return h[i].ref < h[j].ref
}

func (h SessionHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *SessionHeap) Push(s interface{}) {
	s.(*Session).index = len(*h)
	*h = append(*h, s.(*Session))
}

func (h *SessionHeap) Pop() interface{} {
	l := len(*h)
	s := (*h)[l-1]
	s.index = -1
	*h = (*h)[:l-1]
	return s
}

type DialContext struct {
	sync.Mutex
	sessions SessionHeap
}

// goroutine safe
//func Dial() (*DialContext, error) {
//
//	, , mode mgo.Mode
//	c, err := DialWithTimeout(config.Address, int(config.MaxSession), 10*time.Second, 5*time.Minute)
//	return c, err
//}

// goroutine safe
func Dial(config config.DBConfig) (*DialContext, error) {
	s, err := mgo.DialWithTimeout(config.Address, util.GetSecondDuration(config.DialTimeout))
	if err != nil {
		return nil, err
	}
	//最终一致性
	//s.SetMode(mgo.Eventual, false)
	if config.Mode != nil {
		s.SetMode(mgo.Mode(*config.Mode), false)
	}
	s.SetSyncTimeout(util.GetSecondDuration(config.SyncTimeout))
	s.SetSocketTimeout(util.GetSecondDuration(config.SocketTimeout))

	c := new(DialContext)

	// sessions
	c.sessions = make(SessionHeap, config.MaxSession)
	c.sessions[0] = &Session{s, 0, 0}
	for i := 1; i < int(config.MaxSession); i++ {
		c.sessions[i] = &Session{s.New(), 0, i}
	}
	heap.Init(&c.sessions)
	return c, nil
}

// goroutine safe
func (c *DialContext) Close() {
	c.Lock()
	for _, s := range c.sessions {
		s.Close()
		if s.ref != 0 {
			log.Errorf("session ref = %v", s.ref)
		}
	}
	c.Unlock()
}

// goroutine safe
func (c *DialContext) Ref() *Session {
	c.Lock()
	s := c.sessions[0]
	if s.ref == 0 {
		s.Refresh()
	}
	s.ref++
	heap.Fix(&c.sessions, 0)
	c.Unlock()

	return s
}

// goroutine safe
func (c *DialContext) UnRef(s *Session) {
	c.Lock()
	s.ref--
	heap.Fix(&c.sessions, s.index)
	c.Unlock()
}

// goroutine safe
func (c *DialContext) EnsureCounter(db string, collection string, id string, startId int) error {
	s := c.Ref()
	defer c.UnRef(s)

	err := s.DB(db).C(collection).Insert(bson.M{
		"_id": id,
		"seq": startId,
	})
	if mgo.IsDup(err) {
		return nil
	} else {
		return err
	}
}

// goroutine safe
func (c *DialContext) NextSeq(db string, collection string, id string) (int64, error) {
	s := c.Ref()
	defer c.UnRef(s)

	var res struct {
		Seq int64
	}
	_, err := s.DB(db).C(collection).FindId(id).Apply(mgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": 1}},
		ReturnNew: true,
	}, &res)

	return res.Seq, err
}

// goroutine safe
func (c *DialContext) EnsureIndex(db string, collection string, key []string) error {
	s := c.Ref()
	defer c.UnRef(s)

	return s.DB(db).C(collection).EnsureIndex(mgo.Index{
		Key:    key,
		Unique: false,
		Sparse: true,
	})
}

// goroutine safe
func (c *DialContext) EnsureUniqueIndex(db string, collection string, key []string) error {
	s := c.Ref()
	defer c.UnRef(s)

	return s.DB(db).C(collection).EnsureIndex(mgo.Index{
		Key:    key,
		Unique: true,
		Sparse: true,
	})
}
