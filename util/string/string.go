package string

import (
	"fmt"
	"github.com/qb0C80aE/clay/logging"
	mapstructutilpkg "github.com/qb0C80aE/clay/util/mapstruct"
	"strings"
)

var utility = &Utility{}

// Utility handles string operation
type Utility struct {
}

// GetUtility returns the instance of utility
func GetUtility() *Utility {
	return utility
}

// Split splits string by given separator into an slice
func (receiver *Utility) Split(value interface{}, separator string) interface{} {
	data := fmt.Sprintf("%v", value)
	return strings.Split(data, separator)
}

// Join joins given slice with give separator into a string
func (receiver *Utility) Join(slice interface{}, separator string) (interface{}, error) {
	interfaceSlice, err := mapstructutilpkg.GetUtility().SliceToInterfaceSlice(slice)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	stringSlice := make([]string, len(interfaceSlice))

	for index, item := range interfaceSlice {
		stringSlice[index] = fmt.Sprintf("%v", item)
	}

	return strings.Join(stringSlice, separator), nil
}

// Trim eliminate cutset string from both side of given string
func (receiver *Utility) Trim(value interface{}, cutset string) interface{} {
	data := fmt.Sprintf("%v", value)
	return strings.Trim(data, cutset)
}
