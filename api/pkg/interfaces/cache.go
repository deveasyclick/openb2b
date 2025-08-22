package interfaces

type Cache interface {
	Get(key string) interface{}
	Set(key string, value interface{}, ttl int)
	Delete(key string)
	Lock(key string, ttl int)
	Unlock(key string)
}
