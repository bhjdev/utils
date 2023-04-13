package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bhjdev/ulog"
	"github.com/gin-gonic/gin"
)

type Res struct {
	C *gin.Context
}

func (r *Res) Ok(code int, msg interface{}, data interface{}) {

	res := map[string]interface{}{
		"data":   data,
		"status": "success",
		"code":   0,
		"msg":    msg,
	}

	// return res
	r.C.JSON(http.StatusOK, res)

}

func (r *Res) Fail(code int, msg interface{}, data interface{}) {

	ul := &ulog.Log{
		C: r.C,
	}

	type TypeLog struct {
		code int
		msg  interface{}
		data interface{}
	}

	ul.Error(TypeLog{
		code: code,
		msg:  msg,
		data: data,
	})

	res := map[string]interface{}{
		"code":   code,
		"status": "fail",
		"msg":    msg,
		"data":   data,
	}

	r.C.JSON(http.StatusOK, res)

}

// struct to map
func StructToMap(_struct interface{}) (map[string]interface{}, error) {

	_jsonEncode, err := json.Marshal(&_struct)

	if err != nil {
		return nil, err
	}

	_map := make(map[string]interface{})

	json.Unmarshal(_jsonEncode, &_map)

	return _map, nil
}

type TypeConfig struct {
	Env      string
	Debug    bool
	Server   string
	Cdn      string
	Database struct {
		Addr                 string
		User                 string
		Passwd               string
		DBName               string
		AllowNativePasswords bool
	}
}

// get general config
func GetConfig() (TypeConfig, error) {

	// get project path
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var config TypeConfig

	configBytes, err := ioutil.ReadFile(dir + "/../config.json")

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return TypeConfig{}, err
	}

	return config, nil
}

func CatchPanic(res *Res) {
	if err := recover(); err != nil {

		errString := err.(string)

		var errCode int
		var errMsg string

		if strings.Contains(errString, "__") {
			errArr := strings.Split(errString, "__")

			number, _ := strconv.Atoi(errArr[01])

			errCode = number
			errMsg = errArr[0]

		} else {
			errCode = -1
			errMsg = errString
		}

		res.Fail(errCode, errMsg, nil)
	}
}

// format timeStamp to string2
// timeInput: timestamp or "now"
func FormatDatatime(timeTime time.Time, formatType string) string {

	// format current time
	formats := map[string]string{
		"DDMMYYYYhhmmss": "2006-01-02 15:04:05",
		"DDMMYYYYhhmm":   "2006-01-02 15:04",
		"YYYY":           "2006",
		"MM":             "01",
		"DD":             "02",
		"hhmmss":         "15:04:05",
		"hhmm":           "15:04",
		"hh":             "15",
		"mm":             "04",
		"ss":             "05",
	}

	timeOutput := timeTime.Format(formats[formatType])

	return timeOutput

}

func GeneratePassword(originPassword string) string {
	// generate password
	passwordBytes := sha256.Sum256([]byte(originPassword))

	password := string(base64.StdEncoding.EncodeToString(passwordBytes[:]))

	return password
}

func GetOffset(page int, limit int) int {
	return (page - 1) * limit
}
