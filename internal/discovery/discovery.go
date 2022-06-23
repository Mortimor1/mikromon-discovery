package discovery

import (
	"bytes"
	"encoding/json"
	"github.com/Mortimor1/mikromon-discovery/internal/config"
	"github.com/Mortimor1/mikromon-discovery/pkg/logging"
	"github.com/Mortimor1/mikromon-discovery/pkg/utils"
	"io/ioutil"
	"net/http"
	"time"
)

type Discovery struct {
	logger *logging.Logger
	cfg    *config.Config
}

func (d *Discovery) Run(cfg *config.Config) {
	d.logger = logging.GetLogger()
	d.cfg = cfg

	subnets := d.loadSubnets()

	if subnets == nil {
		d.logger.Fatalf("Error: Service core not response. %s", d.cfg.Service.Core)
	}

	subnetsForScan := make([]string, 0)

	for _, subnet := range *subnets {
		if subnet.State == true {
			subnetsForScan = append(subnetsForScan, subnet.Subnet)
		}
	}

	if len(subnetsForScan) == 0 {
		return
	}

	hosts, err := scan(subnetsForScan)
	if err != nil {
		d.logger.Error(err)
	}

	for _, host := range hosts {
		device := d.createDevice(host)
		if device != nil {
			d.logger.Infof("Created new device with address [%s]", device.IpAddress)
		}
	}
	d.logger.Info("Discovery finish")
}

func (d *Discovery) loadSubnets() *[]Subnet {
	c := http.Client{Timeout: time.Duration(5) * time.Second}
	for _, url := range d.cfg.Service.Core {
		resp, err := c.Get(url + "/subnets")
		if err != nil {
			d.logger.Errorf("Error %s", err)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			d.logger.Errorf("Error %s", err)
			continue
		}

		var subnets []Subnet
		err = json.Unmarshal(body, &subnets)

		if err != nil {
			d.logger.Errorf("Error %s", err)
			continue
		}

		resp.Body.Close()

		return &subnets
	}
	return nil
}

func (d *Discovery) createDevice(host int64) *Device {
	device := Device{
		IpAddress: host,
		Name:      "(" + utils.IntToIp(host).String() + ")",
		State:     true,
		Status:    "OK",
	}

	jsonData, err := json.Marshal(device)

	if err != nil {
		d.logger.Errorf("Error %s", err)
	}

	c := http.Client{Timeout: time.Duration(5) * time.Second}
	for _, url := range d.cfg.Service.Core {
		resp, err := c.Post(url+"/devices", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			d.logger.Errorf("Error: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			return &device
		}
	}
	return nil
}
