package Interfaces

type IAccessable interface { 
	Key() string 
}

type IDataAccess interface {
	Create(thingToCreate *IAccessable) (error)
	Delete(key string) (error)
	Update(thingToCreate *IAccessable) (error)
	Get(key string, thingToGet *IAccessable) (error)
}



