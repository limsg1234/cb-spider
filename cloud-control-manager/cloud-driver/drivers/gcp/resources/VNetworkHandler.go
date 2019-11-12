// Proof of Concepts of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// This is a Cloud Driver Example for PoC Test.
//
// program by ysjeon@mz.co.kr, 2019.07.
// modify by devunet@mz.co.kr, 2019.11.

package resources

import (
	"context"
	"errors"
	"strconv"

	compute "google.golang.org/api/compute/v1"

	"time"

	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
	"github.com/davecgh/go-spew/spew"
)

type GCPVNetworkHandler struct {
	Region     idrv.RegionInfo
	Ctx        context.Context
	Client     *compute.Service
	Credential idrv.CredentialInfo
}

func (vNetworkHandler *GCPVNetworkHandler) CreateVNetwork(vNetworkReqInfo irs.VNetworkReqInfo) (irs.VNetworkInfo, error) {

	vNetInfo, errVnet := vNetworkHandler.GetVNetwork(GetCBDefaultVNetName())
	spew.Dump(vNetInfo)
	if errVnet == nil {
		if len(vNetInfo.Name) > 0 {
			return irs.VNetworkInfo{}, errors.New("Already exists.")
		}
	}

	// priject id
	projectID := vNetworkHandler.Credential.ProjectID
	//name := vNetworkReqInfo.Name
	name := GetCBDefaultVNetName()

	network := &compute.Network{
		Name: name,
		//Name:                  GetCBDefaultVNetName(),
		AutoCreateSubnetworks: true, // subnet 자동으로 생성됨
	}

	res, err := vNetworkHandler.Client.Networks.Insert(projectID, network).Do()
	if err != nil {
		cblogger.Error(err)
		return irs.VNetworkInfo{}, err
	}
	cblogger.Info(res)

	//생성되는데 시간이 필요 함. 약 20초정도?
	time.Sleep(time.Second * 20)
	info, err2 := vNetworkHandler.Client.Networks.Get(projectID, name).Do()
	if err2 != nil {
		cblogger.Error(err2)
		return irs.VNetworkInfo{}, err2
	}
	networkInfo := irs.VNetworkInfo{
		Name: info.Name,
		Id:   strconv.FormatUint(info.Id, 10),
		KeyValueList: []irs.KeyValue{
			{"SubnetId", info.Name},
		},
	}

	return networkInfo, nil
}

func (vNetworkHandler *GCPVNetworkHandler) ListVNetwork() ([]*irs.VNetworkInfo, error) {
	projectID := vNetworkHandler.Credential.ProjectID

	vNetworkList, err := vNetworkHandler.Client.Networks.List(projectID).Do()
	if err != nil {

		return nil, err
	}
	var vNetworkInfo []*irs.VNetworkInfo
	for _, item := range vNetworkList.Items {
		networkInfo := irs.VNetworkInfo{
			Name: item.Name,
			Id:   strconv.FormatUint(item.Id, 10),
			KeyValueList: []irs.KeyValue{
				{"SubnetId", item.Name},
			},
		}

		vNetworkInfo = append(vNetworkInfo, &networkInfo)

	}

	return vNetworkInfo, nil
}

func (vNetworkHandler *GCPVNetworkHandler) GetVNetwork(vNetworkID string) (irs.VNetworkInfo, error) {

	projectID := vNetworkHandler.Credential.ProjectID
	//name := vNetworkID
	name := GetCBDefaultVNetName()
	cblogger.Infof("Name : [%s]", name)
	info, err := vNetworkHandler.Client.Networks.Get(projectID, name).Do()
	if err != nil {
		cblogger.Error(err)
		return irs.VNetworkInfo{}, err
	}

	networkInfo := irs.VNetworkInfo{
		Name: info.Name,
		Id:   strconv.FormatUint(info.Id, 10),
		KeyValueList: []irs.KeyValue{
			{"SubnetId", info.Name},
		},
	}

	return networkInfo, nil
}

func (vNetworkHandler *GCPVNetworkHandler) DeleteVNetwork(vNetworkID string) (bool, error) {
	projectID := vNetworkHandler.Credential.ProjectID
	//name := vNetworkID
	name := GetCBDefaultVNetName()
	cblogger.Infof("Name : [%s]", name)
	info, err := vNetworkHandler.Client.Networks.Delete(projectID, name).Do()
	cblogger.Info(info)
	if err != nil {
		cblogger.Error(err)
		return false, err
	}
	//fmt.Println(info)
	return true, nil
}
