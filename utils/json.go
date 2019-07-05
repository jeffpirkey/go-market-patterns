package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func ToJsonString(a interface{}) *string {

	var str string
	data, err := json.Marshal(a)
	if err != nil {
		log.Infof("error marshalling data to json string due to %v", err)
	} else {
		str = string(data)
	}

	return &str
}

func ToJsonBytes(a interface{}) []byte {
	data, err := json.Marshal(a)
	if err != nil {
		log.Infof("error marshalling data to json string due to %v", err)
	}
	return data
}
