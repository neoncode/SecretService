package DataAccess

type MockDAL struct { 
	internalMap map[string]interface{}
}

func GetMockDAL() (IDataAccess) { 
	return new(MockDAL)
}

func (dal *MockDAL) Get(key string, value interface{}) (error) {
	value = dal.internalMap[key]
	return nil
}

func (dal *MockDAL) Set(key string, value interface{}) (error) {
	dal.internalMap[key] = value
	return nil
}

func (dal *MockDAL) Remove(key string) (error) {
	delete(dal.internalMap, key)
	return nil
}
