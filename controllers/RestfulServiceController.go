package Controllers

import(
	"github.com/cummingsi1993@gmail.com/SecretService/Interfaces"
	"github.com/go-martini/martini"
)

func RegisterRestfulService(restModel Interfaces.IRestful, classicMartini *martini.ClassicMartini){
	path := restModel.GetPath()

	classicMartini.Get(path, restModel.RestfulGet)
	classicMartini.Get(path+"/:id", restModel.RestfulGet)

	classicMartini.Put(path, restModel.RestfulPut)
	classicMartini.Put(path+"/:id", restModel.RestfulPut)

	classicMartini.Post(path, restModel.RestfulPost)
	classicMartini.Post(path+"/:id", restModel.RestfulPost)

	classicMartini.Delete(path, restModel.RestfulDelete)
	classicMartini.Delete(path+"/:id", restModel.RestfulDelete)
}