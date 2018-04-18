package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/machine/drivers/amazonec2"

	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/auth"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/drivers/rpc"
	"github.com/docker/machine/libmachine/engine"
	"github.com/docker/machine/libmachine/host"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/swarm"
)

func main() {

	tmpfile, _ := ioutil.TempFile("", "")
	log.SetOutWriter(tmpfile)
	log.SetErrWriter(tmpfile)
	defer tmpfile.Close()

	api, h, err := construct()
	if err != nil {
		panic(fmt.Errorf("failed to construct API client and host config: %v", err))
	}

	// This script is expected to be executed by "go run main.go"
	os.Args = append(os.Args, "create")

	switch os.Args[1] {
	case "rm":
		if err := remove(api, h); err != nil {
			panic(fmt.Errorf("failed to remove machine: %v", err))
		}
	case "create":
		if err := create(api, h); err != nil {
			panic(fmt.Errorf("failed to create machine: %v", err))
		}
	default:
		panic(fmt.Errorf("unsupported command: %v", os.Args[1]))
	}

}

func construct() (*libmachine.Client, *host.Host, error) {

	name := "foobar"
	driver := "amazonec2"

	d := amazonec2.NewDriver("", "")
	d.MachineName = name
	d.StorePath = mcndirs.GetBaseDir()
	d.Region = "ap-northeast-1"
	d.PrivateIPOnly = true

	raw, err := json.Marshal(d)
	if err != nil {
		return nil, nil, err
	}

	api := libmachine.NewClient(mcndirs.GetBaseDir(), mcndirs.GetMachineCertDir())
	h, err := api.NewHost(driver, raw)
	if err != nil {
		return nil, nil, err
	}

	certdir := mcndirs.GetMachineCertDir()
	h.HostOptions = &host.Options{
		AuthOptions: &auth.Options{
			CertDir:          certdir,
			CaCertPath:       filepath.Join(certdir, "ca.pem"),
			CaPrivateKeyPath: filepath.Join(certdir, "ca-key.pem"),
			ClientCertPath:   filepath.Join(certdir, "cert.pem"),
			ClientKeyPath:    filepath.Join(certdir, "key.pem"),
			ServerCertPath:   filepath.Join(mcndirs.GetMachineDir(), name, "server.pem"),
			ServerKeyPath:    filepath.Join(mcndirs.GetMachineDir(), name, "server-key.pem"),
			StorePath:        filepath.Join(mcndirs.GetMachineDir(), name),
		},
		EngineOptions: &engine.Options{
			TLSVerify:  true,
			InstallURL: drivers.DefaultEngineInstallURL,
		},
		SwarmOptions: &swarm.Options{
			IsSwarm:   false,
			Master:    false,
			Discovery: "",
			Host:      "tcp://0.0.0.0:3376",
		},
	}

	return api, h, nil
}

func create(api *libmachine.Client, h *host.Host) error {

	exists, err := api.Exists(h.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("a machine with name %s already exists", h.Name)
	}

	machineflags := h.Driver.GetCreateFlags()
	driveropts := rpcdriver.RPCFlags{
		Values: make(map[string]interface{}),
	}
	for _, f := range machineflags {
		driveropts.Values[f.String()] = f.Default()
		if f.Default() == nil {
			driveropts.Values[f.String()] = false
		}
	}

	// FIXME: Hardcoded
	driveropts.Values["amazonec2-region"] = "ap-northeast-1"

	if err := h.Driver.SetConfigFromFlags(driveropts); err != nil {
		return err
	}

	if err := api.Create(h); err != nil {
		return err
	}

	if err := api.Save(h); err != nil {
		return err
	}

	return nil
}

func remove(api *libmachine.Client, h *host.Host) error {

	exists, err := api.Exists(h.Name)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("a machine with name %s doesn't exist", h.Name)
	}

	h, err = api.Load(h.Name)
	if err != nil {
		return fmt.Errorf("failed to load machine configs: %v", err)
	}

	if err := h.Driver.Remove(); err != nil {
		return fmt.Errorf("failed to remove remote machine: %v", err)
	}

	if err := api.Remove(h.Name); err != nil {
		return fmt.Errorf("failed to remove local machine: %v", err)
	}

	return nil
}
