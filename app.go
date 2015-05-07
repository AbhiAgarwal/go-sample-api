package main

import (
	"net/http"
	"fmt"
	"github.com/julienschmidt/httprouter"
)

type Handler func(res http.ResponseWriter, req *http.Request, ps httprouter.Params)
type Resource map[string]Handler
type Actions struct {
	HandlerName string
	OneHandler Handler
}

var resources = make(map[string]Resource)

func ListHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	resource, exists := resources[ps.ByName("resource")]
	if exists == false {
		//todo: return not found
		return
	}
	action, exists := resource["list"]
	if exists == false {
		//todo: return not found
		return
	}
	action(res, req, ps)
}

func RestDispatch(res http.ResponseWriter, req *http.Request, ps httprouter.Params, actionName string) {
	resource, exists := resources[ps.ByName("resource")]
	if exists == false {
		//todo: return not found
		return
	}
	action, exists := resource[actionName]
	if exists == false {
		//todo: return not found
		return
	}
	action(res, req, ps)
}

func RegisterAction(ResourceName string, actions ...Actions) {
	for _, action := range actions {
		resourceName, actionName := ResourceName, action.HandlerName
		resource, exists := resources[resourceName]
		if exists == false {
			resource = make(Resource)
			resources[resourceName] = resource
		}
		resource[actionName] = action.OneHandler
	}
}

func UserList(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res.Write([]byte("Hello, World"))
	return
}

func main() {
	router := httprouter.New()
	router.GET("/:resource", ListHandler)
	userList := Actions{
		HandlerName: "list",
		OneHandler: UserList,
	}
	RegisterAction("users", userList)
	fmt.Println(http.ListenAndServe(":8080", router))
}
