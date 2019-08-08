package redis

//type Lock struct {
//	resource string
//	token    string
//	conn     redis.Conn
//	timeout  int
//}
//
//func (lock *Lock) tryLock() (ok bool, err error) {
//	_, err = redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int(lock.timeout), "NX"))
//	if err == redis.ErrNil {
//		// The lock was not successful, it already exists.
//		return false, nil
//	}
//	if err != nil {
//		return false, err
//	}
//	return true, nil
//}
//
//func (lock *Lock) Unlock() (err error) {
//	_, err = lock.conn.Do("del", lock.key())
//	return
//}
//
//func (lock *Lock) key() string {
//	return fmt.Sprintf("redislock:%s", lock.resource)
//}
//
//func (lock *Lock) AddTimeout(expire int64) (ok bool, err error) {
//	ttl, err := redis.Int64(lock.conn.Do("TTL", lock.key()))
//	if err != nil {
//		log.Fatal("redis get failed:", err)
//	}
//	if ttl > 0 {
//		_, err := redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int(ttl+expire)))
//		if err == redis.ErrNil {
//			return false, nil
//		}
//		if err != nil {
//			return false, err
//		}
//	}
//	return false, nil
//}
//
//func TryLock(conn redis.Conn, resource string, token string, DefaulTimeout int) (lock *Lock, ok bool, err error) {
//	return TryLockWithTimeout(conn, resource, token, DefaulTimeout)
//}
//
//func TryLockWithTimeout(conn redis.Conn, resource string, token string, timeout int) (lock *Lock, ok bool, err error) {
//	lock = &Lock{resource, token, conn, timeout}
//
//	ok, err = lock.tryLock()
//
//	if !ok || err != nil {
//		lock = nil
//	}
//
//	return
//}
//
//func main() {
//	fmt.Println("start")
//	DefaultTimeout := 10
//	conn, err := redis.Dial("tcp", "localhost:6379")
//
//	lock, ok, err := TryLock(conn, "xiaoru.cc", "token", int(DefaultTimeout))
//	if err != nil {
//		log.Fatal("Error while attempting lock")
//	}
//	if !ok {
//		log.Fatal("bug")
//	}
//	lock.AddTimeout(100)
//
//	time.Sleep(time.Duration(DefaultTimeout) * time.Second)
//	fmt.Println("end")
//	defer lock.Unlock()
//}
