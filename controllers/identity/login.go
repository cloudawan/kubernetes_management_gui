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
	"github.com/astaxie/beego"
	"github.com/cloudawan/kubernetes_management_gui/controllers/utility/guimessagedisplay"
	"strconv"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.TplNames = "identity/login.html"
}

func (c *LoginController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	username := c.GetString("username")
	password := c.GetString("password")
	timeZoneOffset, err := c.GetInt("timeZoneOffset")
	if err != nil {
		guimessage.AddDanger("Fail to get browser time zone offset. Use UTC instead")
	} else {
		hourOffset := float64(timeZoneOffset) / 60.0
		sign := "-"
		if hourOffset < 0 {
			sign = "+"
		}
		guimessage.AddSuccess("Browser time zone is " + sign + strconv.FormatFloat(hourOffset, 'f', -1, 64) + " from UTC")
		c.SetSession("timeZoneOffset", timeZoneOffset)
	}

	// TODO RBAC
	if username == "admin" && password == "password" {
		// Username
		c.SetSession("username", username)
		// Namespace
		namespace := beego.AppConfig.String("namespace")
		c.SetSession("namespace", namespace)

		guimessage.AddSuccess("User " + username + " login")
	}

	c.Ctx.Redirect(302, "/gui/dashboard/topology/")

	guimessage.RedirectMessage(c)
}
