package assembly

import (
	"strconv"
	"wacc_32/assembly/builtins"
	"wacc_32/assembly/instructions"
)

func (cg *CodeGenerator) addStringToBSS(str string) string {
	id := "msg_" + strconv.Itoa(len(cg.bssVars))
	cg.bssVars[id] = instructions.NewStringLiteral(id, str)
	return id
}

func (cg *CodeGenerator) addErrMsgBSS(errType builtins.RuntimeErrType) string {
	id := "err_" + errType.String()
	cg.bssVars[id] = instructions.NewStringLiteral(id, errType.GetMsg())
	return id
}
