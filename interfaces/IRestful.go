package Interfaces

import(
	"net/http"
	"github.com/go-martini/martini"
)

type IRestful interface {
	GetPath() string

	RestfulDelete(params martini.Params)(int, string)

	RestfulGet(params martini.Params)(int, string)

	RestfulPut(params martini.Params, req *http.Request)(int, string)

	RestfulPost(params martini.Params, req *http.Request)(int, string)
}