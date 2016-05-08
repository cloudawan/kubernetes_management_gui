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

package deploybluegreen

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/configuration"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
)

type SelectController struct {
	beego.Controller
}

func (c *SelectController) Get() {
	c.TplName = "deploy/deploybluegreen/select.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")

	imageInformation := c.GetString("imageInformation")
	source := c.GetString("source")

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger("Fail to get deployable namespace")

		// Redirect to list
		if source == "repository" {
			c.Ctx.Redirect(302, "/gui/repository/imageinformation/list")
		} else {
			c.Ctx.Redirect(302, "/gui/deploy/deploybluegreen/list")
		}

		guimessage.RedirectMessage(c)

		return
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploybluegreens/" + imageInformation

	deployBlueGreen := DeployBlueGreen{}

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	statusCode, _, _, err := restclient.RequestWithStructure("GET", url, nil, &deployBlueGreen, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if source == "repository" {
		c.Data["buttonUrlCancel"] = "/gui/repository/imageinformation/list"
	} else {
		c.Data["buttonUrlCancel"] = "/gui/deploy/deploybluegreen/list"
	}

	if err != nil {
		// Other errors
		guimessage.AddDanger("Fail to get deployable namespace")

		// Redirect to list
		if source == "repository" {
			c.Ctx.Redirect(302, "/gui/repository/imageinformation/list")
		} else {
			c.Ctx.Redirect(302, "/gui/deploy/deploybluegreen/list")
		}

		guimessage.RedirectMessage(c)

		return

	} else if statusCode == 404 {
		// Not existing, create
		c.Data["actionButtonValue"] = "Create"
		c.Data["pageHeader"] = "Create Blue Green Deployment"
		c.Data["imageInformation"] = imageInformation
	} else {
		// Update
		c.Data["actionButtonValue"] = "Update"
		c.Data["pageHeader"] = "Update Blue Green Deployment"
		c.Data["imageInformation"] = imageInformation

		autoGeneratedNodePort := false
		if deployBlueGreen.NodePort == 0 {
			autoGeneratedNodePort = true
		}

		c.Data["currentNamespace"] = deployBlueGreen.Namespace
		c.Data["description"] = deployBlueGreen.Description

		if autoGeneratedNodePort {
			c.Data["nodePort"] = ""
			c.Data["checkTagAutoGeneratedNodePort"] = "checked"
			c.Data["hiddenTagNodePort"] = "hidden"
		} else {
			c.Data["nodePort"] = deployBlueGreen.NodePort
			c.Data["checkTagAutoGeneratedNodePort"] = ""
			c.Data["hiddenTagNodePort"] = ""
		}
	}

	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploybluegreens/deployable/" + imageInformation + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	namespaceSlice := make([]string, 0)

	tokenHeaderMap, _ = c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestGetWithStructure(url, &namespaceSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		guimessage.AddDanger("Fail to get deployable namespace")

		// Redirect to list
		if source == "repository" {
			c.Ctx.Redirect(302, "/gui/repository/imageinformation/list")
		} else {
			c.Ctx.Redirect(302, "/gui/deploy/deploybluegreen/list")
		}

		guimessage.RedirectMessage(c)

		return
	}

	if len(namespaceSlice) == 0 {
		guimessage.AddDanger("No deployed application is detected so there is no namespace to select for blue green deployment")

		// Redirect to list
		if source == "repository" {
			c.Ctx.Redirect(302, "/gui/repository/imageinformation/list")
		} else {
			c.Ctx.Redirect(302, "/gui/deploy/deploybluegreen/list")
		}

		guimessage.RedirectMessage(c)
	} else {
		c.Data["namespaceSlice"] = namespaceSlice

		guimessage.OutputMessage(c.Data)
	}
}

func (c *SelectController) Post() {
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")
	kubeapiHost, kubeapiPort, err := configuration.GetAvailableKubeapiHostAndPort()
	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
		guimessage.RedirectMessage(c)
		c.Ctx.Redirect(302, "/gui/deploy/deploybluegreen/list")
		return
	}

	imageInformation := c.GetString("imageInformation")
	namespace := c.GetString("namespace")
	nodePort, _ := c.GetInt("nodePort")
	description := c.GetString("description")
	sessionAffinity := c.GetString("sessionAffinity")

	deployBlueGreen := DeployBlueGreen{
		imageInformation,
		namespace,
		nodePort,
		description,
		sessionAffinity,
		"",
		"",
		"",
	}

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploybluegreens/" + "?kubeapihost=" + kubeapiHost + "&kubeapiport=" + strconv.Itoa(kubeapiPort)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err = restclient.RequestPutWithStructure(url, deployBlueGreen, nil, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		guimessage.AddSuccess("Create blue green deployment " + imageInformation + " success")
	}

	c.Ctx.Redirect(302, "/gui/deploy/deploybluegreen/list")

	guimessage.RedirectMessage(c)
}
