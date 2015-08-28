// Copyright 2015 CloudAwan LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package identity

import (
	"encoding/json"
	"github.com/astaxie/beego/context"
)

const (
	loginURL = "/api/v1/identity/login"
)

func FilterToken(ctx *context.Context) {
	if (ctx.Input.IsGet() || ctx.Input.IsPost()) && ctx.Input.Url() == loginURL {
		// Don't redirect itself to prevent the circle
	} else {
		token := ctx.Input.Header("token")
		err := verifyToken(token)
		if err != nil {
			jsonMap := make(map[string]interface{})
			jsonMap["error"] = err.Error()
			byteSlice, _ := json.Marshal(jsonMap)
			ctx.Output.SetStatus(401)
			ctx.Output.Body(byteSlice)
			//ctx.Redirect(302, loginURL)
		}
	}
}
