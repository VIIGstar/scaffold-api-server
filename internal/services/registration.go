package services

import (
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
)

// AppManagement Use remote config like consul to manage service, it'd be better to access registered app
var AppManagement map[string]interface{}

func init() {
	if AppManagement == nil {
		AppManagement = make(map[string]interface{})
		logrus.Info("Init new app management")
	}
}

func RegisterApp(name string, app interface{}) {
	mux := sync.RWMutex{}
	mux.Lock()
	defer mux.Unlock()
	name = strings.ToLower(name)
	AppManagement[name] = app
}

func GetApp(name string) interface{} {
	name = strings.ToLower(name)
	return AppManagement[name]
}
