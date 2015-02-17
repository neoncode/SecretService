package DataAccess

type IDataAccess interface {
	Remove(key string) (error)
	Set(key string, value interface{}) (error)
	Get(key string, value interface{}) (error)
}



