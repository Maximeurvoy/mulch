package server

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/c2h5oh/datasize"
)

// VMConfig stores needed parameters for a new VM
type VMConfig struct {
	Name        string
	Hostname    string
	Timezone    string
	AppUser     string
	SeedImage   string
	InitUpgrade bool
	DiskSize    uint64
	RAMSize     uint64
	CPUCount    int
	// + prepare scripts
	Prepare []*VMConfigScript
	// + save scripts
	// + restore scripts
}

// VMConfigScript is a script for prepare, save and restore steps
type VMConfigScript struct {
	ScriptFile string
	As         string
}

type tomlVMConfig struct {
	Name        string
	Hostname    string
	Timezone    string
	AppUser     string            `toml:"app_user"`
	SeedImage   string            `toml:"seed_image"`
	InitUpgrade bool              `toml:"init_upgrade"`
	DiskSize    datasize.ByteSize `toml:"disk_size"`
	RAMSize     datasize.ByteSize `toml:"ram_size"`
	CPUCount    int               `toml:"cpu_count"`
	Prepare     []tomlVMConfigScript
}

type tomlVMConfigScript struct {
	ScriptFile string `toml:"script_file"`
	As         string
}

// NewVMConfigFromTomlReader cretes a new VMConfig instance from
// a io.Reader containing VM configuration description
func NewVMConfigFromTomlReader(configIn io.Reader) (*VMConfig, error) {
	vmConfig := &VMConfig{}

	// defaults (if not in the file)
	tConfig := &tomlVMConfig{
		Hostname:    "localhost.localdomain",
		Timezone:    "Europe/Paris",
		AppUser:     "app",
		InitUpgrade: true,
		CPUCount:    1,
	}

	if _, err := toml.DecodeReader(configIn, tConfig); err != nil {
		return nil, err
	}

	if tConfig.Name == "" || !IsValidTokenName(tConfig.Name) {
		return nil, fmt.Errorf("invalid VM name '%s'", tConfig.Name)
	}
	vmConfig.Name = tConfig.Name

	vmConfig.Hostname = tConfig.Hostname
	vmConfig.Timezone = tConfig.Timezone

	if tConfig.AppUser == "" {
		return nil, fmt.Errorf("invalid app_user name '%s'", tConfig.AppUser)
	}
	vmConfig.AppUser = tConfig.AppUser

	// TODO: check the seed image exists
	if tConfig.SeedImage == "" {
		return nil, fmt.Errorf("invalid seed image '%s'", tConfig.SeedImage)
	}
	vmConfig.SeedImage = tConfig.SeedImage

	vmConfig.InitUpgrade = tConfig.InitUpgrade

	if tConfig.DiskSize < 1*datasize.MB {
		return nil, fmt.Errorf("looks like a too small disk (%s)", tConfig.DiskSize)
	}
	vmConfig.DiskSize = tConfig.DiskSize.Bytes()

	if tConfig.RAMSize < 1*datasize.MB {
		return nil, fmt.Errorf("looks like a too small RAM amount (%s)", tConfig.RAMSize)
	}
	vmConfig.RAMSize = tConfig.RAMSize.Bytes()

	if tConfig.CPUCount < 1 {
		return nil, fmt.Errorf("need a least one CPU")
	}
	vmConfig.CPUCount = tConfig.CPUCount

	for _, tScript := range tConfig.Prepare {
		script := &VMConfigScript{}

		if !IsValidTokenName(tScript.As) {
			return nil, fmt.Errorf("'%s' is not a valid user name", tScript.As)
		}
		script.As = tScript.As

		// test readability
		file, err := os.Open(tScript.ScriptFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		script.ScriptFile = tScript.ScriptFile

		vmConfig.Prepare = append(vmConfig.Prepare, script)
	}

	return vmConfig, nil
}
