package framework

import (
	"fmt"
	"github.com/monkiato/game-server-core/pkg/framework/net"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type Plugin interface {
	GetApiName() string
	OnLoad() error
	Initialize(serverManager *ServerManager) error
}

type User struct {
	Id string
	Token string
}

type Context struct {
	User *User
}

type ServerManager struct {
	RestApiModule *net.RestApiModule
	GrpcModule    *net.GrpcModule
}

type TriggerFunction func(ctx *Context, payload string) (string, error)

func NewServerManager(router *net.RestApiModule, grpcServer *net.GrpcModule) *ServerManager {
	return &ServerManager{
		RestApiModule: router,
		GrpcModule:    grpcServer,
	}
}

func (sm ServerManager) RegisterEndpoint(plugin Plugin, methodName string, triggerFunction TriggerFunction) {
	//TODO: validate method name

	sm.RestApiModule.Router.HandleFunc(fmt.Sprintf("/%s/%s", plugin.GetApiName(), methodName), func(writer http.ResponseWriter, request *http.Request) {
		ctx := sm.createContext()

		data, err := ioutil.ReadAll(request.Body)
		if err != nil {
			logrus.Debugf("unable to read body")
			writer.WriteHeader(http.StatusInternalServerError)
			//TODO: handle error message and response body
		}

		response, err := triggerFunction(ctx, string(data))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			//TODO: handle error message and response body
		}

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(response))
	}).Methods(http.MethodPost)
}

func (sm ServerManager) createContext() *Context {
	//TODO: populate context data
	ctx := &Context{}
	return ctx
}
