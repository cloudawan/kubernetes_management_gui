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

package appservice

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/identity"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/rbac"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
	"strconv"
)

type DeployInformation struct {
	Namespace                 string
	ImageInformationName      string
	CurrentVersion            string
	CurrentVersionDescription string
	Description               string
	ReplicaAmount             int
}
type ByDeployInformation []DeployInformation

func (b ByDeployInformation) Len() int           { return len(b) }
func (b ByDeployInformation) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByDeployInformation) Less(i, j int) bool { return b.getIdentifier(i) < b.getIdentifier(j) }
func (b ByDeployInformation) getIdentifier(i int) string {
	return b[i].ImageInformationName + "_" + b[i].Namespace + "_" + b[i].CurrentVersion
}

type DeployClusterApplication struct {
	Name                           string
	Namespace                      string
	Size                           int
	ServiceName                    string
	ReplicationControllerNameSlice []string
}

type ByDeployClusterApplication []DeployClusterApplication

func (b ByDeployClusterApplication) Len() int      { return len(b) }
func (b ByDeployClusterApplication) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByDeployClusterApplication) Less(i, j int) bool {
	return b.getIdentifier(i) < b.getIdentifier(j)
}
func (b ByDeployClusterApplication) getIdentifier(i int) string {
	return b[i].Name + "_" + b[i].Namespace
}

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	guimessage := guimessagedisplay.GetGUIMessage(c)
	c.TplName = "dashboard/deploy/index.html"

	// Authorization for web page display
	c.Data["layoutMenu"] = c.GetSession("layoutMenu")
	// Dashboard tab menu
	user, _ := c.GetSession("user").(*rbac.User)
	c.Data["dashboardTabMenu"] = identity.GetDashboardTabMenu(user, "deploy")

	cloudoneGUIProtocol := beego.AppConfig.String("cloudoneGUIProtocol")
	cloudoneGUIHost := c.Ctx.Input.Host()
	cloudoneGUIPort := c.Ctx.Input.Port()

	c.Data["cloudoneGUIProtocol"] = cloudoneGUIProtocol
	c.Data["cloudoneGUIHost"] = cloudoneGUIHost
	c.Data["cloudoneGUIPort"] = cloudoneGUIPort

	guimessage.OutputMessage(c.Data)
}

type DataController struct {
	beego.Controller
}

func (c *DataController) Get() {
	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	// Json
	c.Data["json"] = make(map[string]interface{})
	c.Data["json"].(map[string]interface{})["applicationView"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["thirdpartyView"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["errorMap"] = make(map[string]interface{})

	// Application view
	applicationJsonMap := make(map[string]interface{})
	applicationJsonMap["name"] = "App View"
	applicationJsonMap["children"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["applicationView"] = append(c.Data["json"].(map[string]interface{})["applicationView"].([]interface{}), applicationJsonMap)

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deploys/"

	deployInformationSlice := make([]DeployInformation, 0)

	tokenHeaderMap, _ := c.GetSession("tokenHeaderMap").(map[string]string)

	_, err := restclient.RequestGetWithStructure(url, &deployInformationSlice, tokenHeaderMap)

	if err != nil {
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.ServeJSON()
		return
	}

	sort.Sort(ByDeployInformation(deployInformationSlice))

	deployInformationMap := make(map[string][]DeployInformation)
	for _, deployInformation := range deployInformationSlice {
		if deployInformationMap[deployInformation.ImageInformationName] == nil {
			deployInformationMap[deployInformation.ImageInformationName] = make([]DeployInformation, 0)
		}
		deployInformationMap[deployInformation.ImageInformationName] = append(deployInformationMap[deployInformation.ImageInformationName], deployInformation)
	}

	applicationViewLeafAmount := 0
	for key, deployInformationSlice := range deployInformationMap {
		deployInformationJsonMap := make(map[string]interface{})
		deployInformationJsonMap["name"] = key
		deployInformationJsonMap["children"] = make([]interface{}, 0)

		for _, deployInformation := range deployInformationSlice {
			namespaceJsonMap := make(map[string]interface{})
			namespaceJsonMap["name"] = deployInformation.Namespace + " (" + strconv.Itoa(deployInformation.ReplicaAmount) + ")"
			namespaceJsonMap["children"] = make([]interface{}, 0)

			versionJsonMap := make(map[string]interface{})
			versionJsonMap["name"] = deployInformation.CurrentVersion + " " + deployInformation.CurrentVersionDescription
			versionJsonMap["children"] = make([]interface{}, 0)

			namespaceJsonMap["children"] = append(namespaceJsonMap["children"].([]interface{}), versionJsonMap)
			deployInformationJsonMap["children"] = append(deployInformationJsonMap["children"].([]interface{}), namespaceJsonMap)
		}

		applicationJsonMap["children"] = append(applicationJsonMap["children"].([]interface{}), deployInformationJsonMap)
		applicationViewLeafAmount += 1
	}

	// Third-party view
	thirdpartyJsonMap := make(map[string]interface{})
	thirdpartyJsonMap["name"] = "3rd party View"
	thirdpartyJsonMap["children"] = make([]interface{}, 0)
	c.Data["json"].(map[string]interface{})["thirdpartyView"] = append(c.Data["json"].(map[string]interface{})["thirdpartyView"].([]interface{}), thirdpartyJsonMap)

	url = cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/deployclusterapplications/"

	deployClusterApplicationSlice := make([]DeployClusterApplication, 0)

	_, err = restclient.RequestGetWithStructure(url, &deployClusterApplicationSlice, tokenHeaderMap)

	if identity.IsTokenInvalidAndRedirect(c, c.Ctx, err) {
		return
	}

	if err != nil {
		c.Data["json"].(map[string]interface{})["error"] = err.Error()
		c.ServeJSON()
		return
	}

	sort.Sort(ByDeployClusterApplication(deployClusterApplicationSlice))

	deployClusterApplicationMap := make(map[string][]DeployClusterApplication)
	for _, deployClusterApplication := range deployClusterApplicationSlice {
		if deployClusterApplicationMap[deployClusterApplication.Name] == nil {
			deployClusterApplicationMap[deployClusterApplication.Name] = make([]DeployClusterApplication, 0)
		}
		deployClusterApplicationMap[deployClusterApplication.Name] = append(deployClusterApplicationMap[deployClusterApplication.Name], deployClusterApplication)
	}

	thirdpartyViewLeafAmount := 0
	for deployClusterApplicationName, deployClusterApplicationSlice := range deployClusterApplicationMap {
		deployClusterApplicationJsonMap := make(map[string]interface{})
		deployClusterApplicationJsonMap["name"] = deployClusterApplicationName
		deployClusterApplicationJsonMap["children"] = make([]interface{}, 0)

		for _, deployClusterApplication := range deployClusterApplicationSlice {
			namespaceJsonMap := make(map[string]interface{})
			namespaceJsonMap["name"] = deployClusterApplication.Namespace + " (" + strconv.Itoa(deployClusterApplication.Size) + ")"
			namespaceJsonMap["children"] = make([]interface{}, 0)

			for _, replicationControllerName := range deployClusterApplication.ReplicationControllerNameSlice {
				replicationControllerNameJsonMap := make(map[string]interface{})
				replicationControllerNameJsonMap["name"] = replicationControllerName
				replicationControllerNameJsonMap["children"] = make([]interface{}, 0)

				namespaceJsonMap["children"] = append(namespaceJsonMap["children"].([]interface{}), replicationControllerNameJsonMap)

				thirdpartyViewLeafAmount += 1
			}

			deployClusterApplicationJsonMap["children"] = append(deployClusterApplicationJsonMap["children"].([]interface{}), namespaceJsonMap)
		}

		thirdpartyJsonMap["children"] = append(thirdpartyJsonMap["children"].([]interface{}), deployClusterApplicationJsonMap)
	}

	c.Data["json"].(map[string]interface{})["applicationViewLeafAmount"] = applicationViewLeafAmount
	c.Data["json"].(map[string]interface{})["thirdpartyViewLeafAmount"] = thirdpartyViewLeafAmount

	c.ServeJSON()
}
