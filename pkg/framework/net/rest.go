package net

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

type RestApiModule struct {
	Router *mux.Router
	port   int
}

func (m *RestApiModule) Run() {
	//print all available routes
	walkRoute(m.Router)

	addr := fmt.Sprintf(":%d", m.port)
	srv := &http.Server{
		Handler:      m.Router,
		Addr:         addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Printf("running http server at %s", addr)
	log.Fatal(srv.ListenAndServe())
}

func NewRestApiModule(port int) *RestApiModule {
	router := mux.NewRouter()
	router.Use(DecompressMiddleware)
	router.Use(JsonResponseTypeMiddleware)
	addListRoutesEndpoint(router)
	return &RestApiModule {
		Router: router,
		port: port,
	}
}

func addListRoutesEndpoint(route *mux.Router) {
	logrus.Debug("adding all routes list...")
	route.HandleFunc("/routes", func(writer http.ResponseWriter, request *http.Request) {
		routesMap := walkRoute(route)

		data, _ := json.Marshal(routesMap)
		writer.Write(data)
	})
}

func walkRoute(route *mux.Router) map[string][]string {
	routesMap := map[string][]string{}

	route.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		for _, method := range methods {
			routesMap[tpl] = append(routesMap[tpl], method)
			logrus.Debugf("%s	%s", method, tpl)
		}
		return nil
	})

	return routesMap
}
