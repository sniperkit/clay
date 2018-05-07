package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
)

type ephemeralBinaryObjectRawController struct {
	BaseController
}

func newEphemeralBinaryObjectRawController() *ephemeralBinaryObjectRawController {
	return CreateController(&ephemeralBinaryObjectRawController{}, model.NewEphemeralBinaryObjectRaw()).(*ephemeralBinaryObjectRawController)
}

func (receiver *ephemeralBinaryObjectRawController) GetResourceSingleURL() (string, error) {
	modelKey, err := receiver.model.GetModelKey(receiver.model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	ephemeralBinaryObjectResourceName, err := extension.GetAssociatedResourceNameWithModel(model.NewEphemeralBinaryObject())
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/:%s/raw", ephemeralBinaryObjectResourceName, modelKey.KeyParameter), nil
}

func (receiver *ephemeralBinaryObjectRawController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
		},
	}
	return routeMap
}

func (receiver *ephemeralBinaryObjectRawController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	data := result.([]byte)

	contentType := c.Request.URL.Query().Get("content_type")

	switch contentType {
	case "":
		c.Data(code, "application/octet-stream", []byte(data))
	default:
		c.Data(code, contentType, []byte(data))
	}
}

func init() {
	extension.RegisterController(newEphemeralBinaryObjectRawController())
}
