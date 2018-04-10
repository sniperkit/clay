package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/model"
)

type commandExecutionController struct {
	BaseController
}

func newCommandExecutionController() *commandExecutionController {
	return CreateController(&commandExecutionController{}, model.NewCommandExecution()).(*commandExecutionController)
}

func (receiver *commandExecutionController) GetResourceSingleURL() string {
	commandResourceName := extension.GetAssociateResourceNameWithModel(model.NewCommand())
	return fmt.Sprintf("%s/:id/execution", commandResourceName)
}

func (receiver *commandExecutionController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodPost: {
			extension.URLSingle: receiver.Create,
		},
		extension.MethodDelete: {
			extension.URLSingle: receiver.Delete,
		},
	}
	return routeMap
}

func init() {
	extension.RegisterController(newCommandExecutionController())
}