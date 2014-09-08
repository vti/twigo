package action

import "github.com/vti/twigo/app"

type BaseAction struct {
	Context *app.Context
}

func (action *BaseAction) SetContext(context *app.Context) {
	action.Context = context
}
