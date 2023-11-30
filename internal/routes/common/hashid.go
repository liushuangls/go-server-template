package common

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/liushuangls/go-server-template/internal/clients"
)

var hashID *clients.HashID

func SetHashID(h *clients.HashID) {
	hashID = h
}

func ShouldBindWithHashID(c *gin.Context, obj any, b ...binding.Binding) error {
	if len(b) > 0 {
		if err := c.ShouldBindWith(obj, b[0]); err != nil {
			return err
		}
	} else {
		if err := c.ShouldBind(obj); err != nil {
			return err
		}
	}
	return ParseHashID(obj)
}

// example: UserHashID string `hashID:"target=UserID,type=User"`
func ParseHashID(obj any) error {
	val := reflect.ValueOf(obj).Elem()
	for i := 0; i < val.NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get("hashID")
		if tag == "" {
			continue
		}

		hashStr, ok := val.Field(i).Interface().(string)
		if !ok || hashStr == "" {
			continue
		}

		targetField, targetValue, err := parseHashTag(tag, hashStr)
		if err != nil {
			return err
		}
		if targetField != "" && targetValue != 0 {
			if target := val.FieldByName(targetField); target.CanSet() && target.CanInt() {
				target.SetInt(int64(targetValue))
			}
		}
	}
	return nil
}

func parseHashTag(tag, hashStr string) (targetField string, targetVal int, err error) {
	tags := strings.Split(tag, ",")
	for _, t := range tags {
		items := strings.Split(t, "=")
		if len(items) != 2 {
			continue
		}
		k, v := items[0], items[1]
		switch k {
		case "target":
			targetField = v
		case "type":
			hashCli := hashID.GetHashID(v)
			if hashCli != nil {
				targetVal, err = hashCli.Decode(hashStr)
				if err != nil {
					return
				}
			}
		}
	}
	return
}
