package sessions

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	return &RedisStore{Client: client, SessionDuration: sessionDuration}
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	// marshal the sessionState to JSON
	value, err := json.Marshal(sessionState)
	if err != nil {
		return err
	}

	// save to redis database
	client := rs.Client
	time := rs.SessionDuration
	key := sid.getRedisKey()
	err = client.Set(key, value, time).Err()
	if err != nil {
		return err
	}

	return nil
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	// get the session state from redis
	client := rs.Client
	key := sid.getRedisKey()
	sess, err := client.Get(key).Result()
	if err != nil {
		return ErrStateNotFound
	}

	// reset expiry time
	_, err = client.Expire(key, rs.SessionDuration).Result()
	if err != nil {
		return err
	}

	// unmarshall value back to sessionState parameter
	buffer := []byte(sess)
	return json.Unmarshal(buffer, sessionState)
}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	// delete the data stored in redis for the provided SessionID
	key := sid.getRedisKey()
	client := rs.Client

	_, err := client.Del(key).Result()
	if err != nil {
		return err
	}

	return nil
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	//convert the SessionID to a string and add the prefix "sid:" to keep
	//SessionID keys separate from other keys that might end up in this
	//redis instance
	return "sid:" + sid.String()
}
