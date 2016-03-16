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

package notifier

import (
	"github.com/astaxie/beego"
	"github.com/cloudawan/cloudone_gui/controllers/utility/guimessagedisplay"
	"github.com/cloudawan/cloudone_utility/restclient"
	"sort"
	"time"
)

type ListController struct {
	beego.Controller
}

type ReplicationControllerNotifier struct {
	Check             bool
	CoolDownDuration  time.Duration
	RemainingCoolDown time.Duration
	KubeapiHost       string
	KubeapiPort       int
	Namespace         string
	Kind              string
	Name              string
	NotifierSlice     []Notifier
	IndicatorSlice    []Indicator
}

type Notifier struct {
	Kind string
	Data string
}

type NotifierSMSNexmo struct {
	Destination         string
	Sender              string
	ReceiverNumberSlice []string
}

type NotifierEmail struct {
	Destination          string
	ReceiverAccountSlice []string
}

type Indicator struct {
	Type                  string
	AboveAllOrOne         bool
	AbovePercentageOfData float64
	AboveThreshold        int64
	BelowAllOrOne         bool
	BelowPercentageOfData float64
	BelowThreshold        int64
}

type ByReplicationControllerNotifier []ReplicationControllerNotifier

func (b ByReplicationControllerNotifier) Len() int           { return len(b) }
func (b ByReplicationControllerNotifier) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByReplicationControllerNotifier) Less(i, j int) bool { return b[i].Name < b[j].Name }

func (c *ListController) Get() {
	c.TplName = "notification/notifier/list.html"
	guimessage := guimessagedisplay.GetGUIMessage(c)

	cloudoneProtocol := beego.AppConfig.String("cloudoneProtocol")
	cloudoneHost := beego.AppConfig.String("cloudoneHost")
	cloudonePort := beego.AppConfig.String("cloudonePort")

	url := cloudoneProtocol + "://" + cloudoneHost + ":" + cloudonePort +
		"/api/v1/notifiers/"

	replicationControllerNotifierSlice := make([]ReplicationControllerNotifier, 0)

	_, err := restclient.RequestGetWithStructure(url, &replicationControllerNotifierSlice)

	if err != nil {
		// Error
		guimessage.AddDanger(err.Error())
	} else {
		sort.Sort(ByReplicationControllerNotifier(replicationControllerNotifierSlice))
		c.Data["replicationControllerNotifierSlice"] = replicationControllerNotifierSlice
	}

	guimessage.OutputMessage(c.Data)
}
