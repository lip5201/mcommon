package mcommon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
	"github.com/twinj/uuid"
	"gopkg.in/go-playground/validator.v8"
)

type GinResp struct {
	ErrCode int64  `json:"error"`
	ErrMsg  string `json:"error_msg"`
	Data    gin.H  `json:"data,omitempty"`
}

// GinRespSuccess 成功返回
var GinRespSuccess = GinResp{
	ErrCode: ErrorSuccess,
	ErrMsg:  ErrorSuccessMsg,
}

// GinRespInternalErr 成功返回
var GinRespInternalErr = GinResp{
	ErrCode: ErrorInternal,
	ErrMsg:  ErrorInternalMsg,
}

// GetUUIDStr 获取唯一字符串
func GetUUIDStr() string {
	u1 := uuid.NewV4()
	return strings.Replace(u1.String(), "-", "", -1)
}

// IsStringInSlice 字符串是否在数组中
func IsStringInSlice(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// IsIntInSlice 数字是否在数组中
func IsIntInSlice(arr []int64, str int64) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// GinRepeatReadBody 创建可重复度body
func GinRepeatReadBody(c *gin.Context) error {
	var err error
	var body []byte
	if cb, ok := c.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			body = cbb
		}
	}
	if body == nil {
		body, err = ioutil.ReadAll(c.Request.Body)
		if err != nil {
			Log.Errorf("err: [%T] %s", err, err.Error())
			c.Abort()
			return err
		}
		c.Set(gin.BodyBytesKey, body)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return nil
}

// GinShouldBindRepeat 可重复绑定参数
func GinShouldBindRepeat(c *gin.Context, obj interface{}) error {
	err := GinRepeatReadBody(c)
	if err != nil {
		return err
	}
	return c.ShouldBind(obj)
}

// GinFillBindError 检测gin输入绑定错误
func GinFillBindError(c *gin.Context, err error) {
	validatorError, ok := err.(validator.ValidationErrors)
	if ok {
		errMsgList := make([]string, 0, 16)
		for _, v := range validatorError {
			errMsgList = append(errMsgList, fmt.Sprintf("[%s] is %s", strcase.ToSnake(v.Field), v.ActualTag))
		}
		c.JSON(
			http.StatusOK,
			GinResp{
				ErrCode: ErrorBind,
				ErrMsg:  strings.Join(errMsgList, ", "),
				Data:    nil,
			},
		)
		return
	}
	unmarshalError, ok := err.(*json.UnmarshalTypeError)
	if ok {
		c.JSON(
			http.StatusOK,
			GinResp{
				ErrCode: ErrorBind,
				ErrMsg:  fmt.Sprintf("[%s] type error", unmarshalError.Field),
			},
		)
		return
	}
	if err == io.EOF {
		c.JSON(
			http.StatusOK,
			GinResp{
				ErrCode: ErrorBind,
				ErrMsg:  fmt.Sprintf("empty body"),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		GinResp{
			ErrCode: ErrorInternal,
			ErrMsg:  ErrorInternalMsg,
		},
	)
}

// GinFillSuccessData 填充返回数据
func GinFillSuccessData(data gin.H) GinResp {
	return GinResp{
		ErrCode: ErrorSuccess,
		ErrMsg:  ErrorSuccessMsg,
		Data:    data,
	}
}

// GinDoRespSuccess 返回成功信息
func GinDoRespSuccess(c *gin.Context, data gin.H) {
	c.JSON(http.StatusOK, GinResp{
		ErrCode: ErrorSuccess,
		ErrMsg:  ErrorSuccessMsg,
		Data:    data,
	})
}

// 返回错误信息
func GinDoRespInternalErr(c *gin.Context) {
	c.JSON(http.StatusOK, GinResp{
		ErrCode: ErrorInternal,
		ErrMsg:  ErrorInternalMsg,
	})
}

// GinDoRespErr 返回特殊错误
func GinDoRespErr(c *gin.Context, code int64, msg string, data gin.H) {
	c.JSON(http.StatusOK, GinResp{
		ErrCode: code,
		ErrMsg:  msg,
		Data:    data,
	})
}
