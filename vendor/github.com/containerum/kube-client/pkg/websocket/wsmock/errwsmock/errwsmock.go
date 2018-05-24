// Code generated by noice. DO NOT EDIT.
package errwsmock

import (
	bytes "bytes"
	template "text/template"

	cherry "github.com/containerum/cherry"
)

const ()

func ErrUpgradeFailed(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "upgrade failed", StatusHTTP: 500, ID: cherry.ErrID{SID: "wsmock", Kind: 0x1}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}
func renderTemplate(templText string) string {
	buf := &bytes.Buffer{}
	templ, err := template.New("").Parse(templText)
	if err != nil {
		return err.Error()
	}
	err = templ.Execute(buf, map[string]string{})
	if err != nil {
		return err.Error()
	}
	return buf.String()
}
