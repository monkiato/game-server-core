package framework

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/monkiato/game-server-core/pkg/framework/net"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Plugin interface must be used to create extending plugins
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
	DB *gorm.DB
}

type DBConfig struct {
	Host string
	Port int
	User string
	DBName string
	Password string
}

type TriggerFunction func(ctx *Context, payload string) (string, error)

func NewServerManager(router *net.RestApiModule, grpcServer *net.GrpcModule, dbConfig *DBConfig) *ServerManager {
	db, err := connectDB(dbConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database.\n%v", err))
	}

	return &ServerManager{
		RestApiModule: router,
		GrpcModule:    grpcServer,
		DB: db,
	}
}

func connectDB(config *DBConfig) (*gorm.DB, error) {
	return gorm.Open( "postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		config.Host,
		config.Port,
		config.User,
		config.DBName,
		config.Password))
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
