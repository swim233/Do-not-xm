package module

import (
	"encoding/json"

	"github.com/swim233/do_not_xm/logger"
)

func ParserData2Json() {

}

func ParseJson2Data[T any](k string) (T, error) {
	var data T
	var err error
	if err = json.Unmarshal([]byte(k), &data); err != nil {
		logger.Suger.Errorf("Fail to parser json to data : %s", err)
		logger.Suger.Errorf("Json information : %s", k)
	}
	return data, err

}
