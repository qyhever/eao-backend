package templates

import (
	_ "embed"
	"strings"
)

const registerCodePlaceholder = "{{CODE}}"

//go:embed send.html
var registerCodeTemplateA1 string

func RenderRegisterCodeTemplateA1(code string) string {
	return strings.ReplaceAll(registerCodeTemplateA1, registerCodePlaceholder, code)
}
