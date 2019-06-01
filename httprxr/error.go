package httprxr

import "fmt"

const InvalidParamErrorCode = "invalid_param"

//export
func InvalidParamError(param string, value ...interface{}) ResponseMessage {
	message := ""
	if len(value) > 0 {
		message = fmt.Sprintf("invalid param [%s] with value [%v]", param, value[0])
	} else {
		message = fmt.Sprintf("invalid param : %s", param)
	}
	return NewErrorMessage(InvalidParamErrorCode, message)
}
