package oneandone

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/docker/machine/drivers"
	"github.com/docker/machine/log"
	"github.com/docker/machine/ssh"
	"github.com/docker/machine/state"
	oaocs "github.com/jlusiardi/oneandone-cloudserver-api"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	minCores = 1
	maxCores = 16
	minRam   = 1
	maxRam   = 128
	minSsd   = 20
	maxSsd   = 500
	stepSsd  = 20
)

const Endpoint string = ""

type Driver struct {
	Endpoint       string
	AccessToken    string
	VmId           string
	FirewallId     string
	MachineName    string
	CaCertPath     string
	PrivateKeyPath string
	StorePath      string
	IPAddress      string
	Cores          int
	Ram            int
	Ssd            int
}

func init() {
	drivers.Register("oneandone", &drivers.RegisteredDriver{
		New:            NewDriver,
		GetCreateFlags: GetCreateFlags,
	})
}

func GetCreateFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			EnvVar: "ONEANDONE_ACCESS_TOKEN",
			Name:   "oneandone-access-token",
			Usage:  "1&1 access token",
		},
		cli.StringFlag{
			EnvVar: "ONEANDONE_CORES",
			Name:   "oneandone-cores",
			Usage:  "number of cores for the Docker Host (" + strconv.Itoa(minCores) + "-" + strconv.Itoa(maxCores) + ")",
		},
		cli.StringFlag{
			EnvVar: "ONEANDONE_RAM",
			Name:   "oneandone-ram",
			Usage:  "size of RAM for the Docker Host in GB (" + strconv.Itoa(minRam) + "-" + strconv.Itoa(maxRam) + ")",
		},
		cli.StringFlag{
			EnvVar: "ONEANDONE_SSD",
			Name:   "oneandone-ssd",
			Usage:  "size of the SSD for the Docker Host in GB (" + strconv.Itoa(minSsd) + "-" + strconv.Itoa(maxSsd) + ", steps of " + strconv.Itoa(stepSsd) + ")",
		},
		cli.StringFlag{
			EnvVar: "ONEANDONE_ENDPOINT",
			Name:   "oneandone-endpoint",
			Usage:  "1&1 CloudServer rest api endpoint",
		},
	}
}

func NewDriver(machineName string, storePath string, caCert string, privateKey string) (drivers.Driver, error) {
	return &Driver{MachineName: machineName, StorePath: storePath, CaCertPath: caCert, PrivateKeyPath: privateKey}, nil
}

func (d *Driver) DriverName() string {
	return "oneandone"
}

func (d *Driver) AuthorizePort(ports []*drivers.Port) error {
	return nil
}

func (d *Driver) DeauthorizePort(ports []*drivers.Port) error {
	return nil
}

func (d *Driver) Create() error {
	log.Infof("Creating a new 1&1 CloudServer ... %v", d.FirewallId)

	appliance, err := d.getAPI().ServerApplianceFindNewest("Linux", "Ubuntu", "Minimal", 64, true)
	if err != nil {
		return err
	}
	log.Debugf("Auto-select appliance '%v' as base image", appliance.Name)
	firewall, err := d.getAPI().CreateFirewallPolicy(oaocs.FirewallPolicyCreateData{
		Name:        "[Docker Machine] " + d.MachineName,
		Description: "Firewall policy for docker machine " + d.MachineName,
		Rules: []oaocs.FirewallPolicyRulesCreateData{
			oaocs.FirewallPolicyRulesCreateData{
				Protocol: "TCP",
				PortFrom: oaocs.Int2Pointer(1),
				PortTo:   oaocs.Int2Pointer(65535),
				SourceIp: "0.0.0.0",
			},
		},
	})
	if err != nil {
		return err
	}
	log.Debugf("create firewall policy with id '%v'", firewall.Id)
	d.FirewallId = firewall.Id

	server, err := d.getAPI().CreateServer(oaocs.ServerCreateData{
		Name:             "[Docker Machine] " + d.MachineName,
		Description:      d.MachineName + " created by docker machine",
		ApplianceId:      appliance.Id,
		FirewallPolicyId: d.FirewallId,
		Hardware: oaocs.Hardware{
			CoresPerProcessor: 1,
			Vcores:            d.Cores,
			Ram:               d.Ram,
			Hdds: []oaocs.Hdd{
				oaocs.Hdd{
					IsMain: true,
					Size:   d.Ssd,
				},
			},
		},
		PowerOn: true,
	})

	if err != nil {
		d.Remove()
		return err
	}
	d.VmId = server.Id

	firewall.WaitForState("ACTIVE")
	server.WaitForState("POWERED_ON")

	server, _ = d.getAPI().GetServer(d.VmId)
	d.IPAddress = server.Ips[0].Ip

	// create and install SSH key
	log.Infof("Generating SSH key ...")
	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return err
	}
	d.WaitForServerReady(server)
	err = d.installSshKey(server.Password)
	if err != nil {
		d.Remove()
		return err
	}

	log.Infof("Sucessfully created a new 1&1 CloudServer with IP: '%s'", d.IPAddress)
	return nil
}

func (d *Driver) WaitForServerReady(server *oaocs.Server) error {
	log.Infof("Waiting for SSH to get ready ...")

	sshPort, _ := d.GetSSHPort()
	WaitForTcpPortToBeOpen(d.IPAddress, sshPort)

	log.Infof("Waiting for package manager to get ready ...")
	client, err := getSSHClient(d.GetSSHUsername(), d.IPAddress, sshPort, server.Password)
	if err != nil {
		return fmt.Errorf("Failed to establish an ssh session to the server")
	}
	result, _ := executeCmd(client, "ps -C aptitude >/dev/null && echo 1 || echo 0")
	for !strings.Contains(result, "0") {
		result, _ = executeCmd(client, "ps -C aptitude >/dev/null && echo 1 || echo 0")
		log.Debugf("Waiting for package manager to get ready. Retry in 5 sec ...")
		time.Sleep(5 * time.Second)
	}
	return nil
}

func (d *Driver) installSshKey(password string) error {
	fileBytes, err := ioutil.ReadFile(d.GetSSHKeyPath() + ".pub")
	if err != nil {
		return fmt.Errorf("Cannot read SSH public key: %v", err)
	}
	key := string(fileBytes)

	sshPort, _ := d.GetSSHPort()
	client, err := getSSHClient(d.GetSSHUsername(), d.IPAddress, sshPort, password)
	if err != nil {
		return fmt.Errorf("Cannot create SSH client to connect to server: %v", err)
	}
	_, err = executeCmd(client, "mkdir -p ~/.ssh; chmod 700 ~/.ssh; echo \""+key+"\" >> ~/.ssh/authorized_keys")
	if err != nil {
		return fmt.Errorf("Cannot install SSH public key on server: %v", err)
	}
	return nil
}

func (d *Driver) GetIP() (string, error) {
	return d.IPAddress, nil
}

func (d *Driver) GetMachineName() string {
	return d.MachineName
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.IPAddress, nil
}

func (d *Driver) GetSSHKeyPath() string {
	return filepath.Join(d.StorePath, "id_rsa")
}

func (d *Driver) GetSSHPort() (int, error) {
	return 22, nil
}

func (d *Driver) GetSSHUsername() string {
	return "root"
}

func (d *Driver) GetState() (state.State, error) {
	vm, err := d.getAPI().GetServer(d.VmId)
	if err != nil {
		return state.None, err
	}

	switch vm.Status.State {
	case "POWERING_ON":
		return state.Starting, nil
	case "REBOOTING":
	case "POWERED_ON":
		return state.Running, nil
	case "POWERED_OFF":
		return state.Stopped, nil
	case "POWERING_OFF":
		return state.Stopping, nil
	case "REMOVING":
	case "CONFIGURING":
	case "DEPLOYING":
		return state.Error, nil
	}

	return state.None, nil
}

func (d *Driver) PreCreateCheck() error {
	//server name available
	servers, serverErr := d.getAPI().GetServers()
	if serverErr != nil {
		return fmt.Errorf("Failed to validate that the server name is available: %v", serverErr)
	}
	for index, _ := range servers {
		if servers[index].Name == "[Docker Machine] "+d.MachineName {
			return fmt.Errorf("For the given name: '%s' already exists a 1&1 CloudServer", d.MachineName)
		}
	}

	//firewall policy name available
	fwp, fwpErr := d.getAPI().GetFirewallPolicies()
	if fwpErr != nil {
		return fmt.Errorf("Failed to validate that the firewall policy name is available: %v", fwpErr)
	}
	for index, _ := range fwp {
		if fwp[index].Name == "[Docker Machine] "+d.MachineName {
			return fmt.Errorf("For the given name: '%s' already exists a firewall policy", d.MachineName)
		}
	}
	return nil
}

func (d *Driver) GetURL() (string, error) {
	return fmt.Sprintf("tcp://%s:2376", d.IPAddress), nil
}

func (d *Driver) Remove() error {
	log.Infof("Removing the 1&1 CloudServer named '%s' ...", d.MachineName)

	// delete firewall (if still exists)
	firewall, err := d.getAPI().GetFirewallPolicy(d.FirewallId)
	if err == nil {
		firewall, err = firewall.Delete()
		if err != nil {
			log.Debugf("Deleting firewall caused error: %v", err)
		}
	} else {
		log.Debugf("Finding firewall caused error: %v", err)
	}

	server, err := d.getAPI().GetServer(d.VmId)
	if err == nil {
		server, err = server.Delete()
		if err != nil {
			log.Debugf("Deleting server caused error: %v", err)
		}
	} else {
		log.Debugf("Finding server caused error: %v", err)
	}

	if firewall != nil {
		firewall.WaitUntilDeleted()
	}
	if server != nil {
		server.WaitUntilDeleted()
	}
	return nil
}

func (d *Driver) Start() error {
	log.Infof("Starting the 1&1 CloudServer named '%s' ...", d.MachineName)
	server, getErr := d.getAPI().GetServer(d.VmId)
	if getErr != nil {
		return fmt.Errorf("Failed to start the 1&1 CloudServer: %v", getErr)
	}
	_, startErr := server.Start()
	if startErr != nil {
		return fmt.Errorf("Failed to start the 1&1 CloudServer: %v", startErr)
	}
	return nil
}

func (d *Driver) Stop() error {
	log.Infof("Stopping the 1&1 CloudServer named '%s' ...", d.MachineName)
	server, getErr := d.getAPI().GetServer(d.VmId)
	if getErr != nil {
		return fmt.Errorf("Failed to stop the 1&1 CloudServer: %v", getErr)
	}
	_, stopErr := server.Shutdown(false)
	if stopErr != nil {
		return fmt.Errorf("Failed to stop the 1&1 CloudServer: %v", stopErr)
	}
	return nil
}

func (d *Driver) Restart() error {
	log.Infof("Restarting the 1&1 CloudServer named '%s' ...", d.MachineName)
	server, getErr := d.getAPI().GetServer(d.VmId)
	if getErr != nil {
		return fmt.Errorf("Failed to restart the 1&1 CloudServer: %v", getErr)
	}
	_, restartErr := server.Reboot(false)
	if restartErr != nil {
		return fmt.Errorf("Failed to restart the 1&1 CloudServer: %v", restartErr)
	}
	return nil
}


func (d *Driver) Kill() error {
	log.Infof("Killing the 1&1 CloudServer named '%s' ...", d.MachineName)
	server, getErr := d.getAPI().GetServer(d.VmId)
	if getErr != nil {
		return fmt.Errorf("Failed to kill the 1&1 CloudServer: %v", getErr)
	}
	_, killErr := server.Shutdown(true)
	if killErr != nil {
		return fmt.Errorf("Failed to kill the 1&1 CloudServer: %v", killErr)
	}
	return nil
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.Endpoint = flags.String("oneandone-endpoint")
	if d.Endpoint == "" {
		return fmt.Errorf("oneandone driver requires the --oneandone-endpoint option")
	}
	d.AccessToken = flags.String("oneandone-access-token")
	if d.AccessToken == "" {
		return fmt.Errorf("oneandone driver requires the --oneandone-access-token option")
	}
	d.Cores = flags.Int("oneandone-cores")
	if d.Cores == 0 {
		log.Debugf("no number of cores specified, use %v core", minCores)
		d.Cores = minCores
	}
	if d.Cores < minCores || d.Cores > maxCores {
		return fmt.Errorf("oneandone driver requires the --oneandone-cores option to be an integer (" + strconv.Itoa(minCores) + "-" + strconv.Itoa(maxCores) + ")")
	}
	d.Ram = flags.Int("oneandone-ram")
	if d.Ram == 0 {
		log.Debugf("no amount of RAM specified, use %v GB", minRam)
		d.Ram = minRam
	}
	if d.Ram < minRam || d.Ram > maxRam {
		return fmt.Errorf("oneandone driver requires the --oneandone-ram option to be an integer (" + strconv.Itoa(minRam) + "-" + strconv.Itoa(maxRam) + ")")
	}
	d.Ssd = flags.Int("oneandone-ssd")
	if d.Ssd == 0 {
		log.Debugf("no amount of SSD specified, use %v GB", minSsd)
		d.Ssd = minSsd
	}
	if d.Ssd < minSsd || d.Ssd > maxSsd || (d.Ssd%stepSsd) != 0 {
		return fmt.Errorf("oneandone driver requires the --oneandone-ssd option to be an integer (" + strconv.Itoa(minSsd) + "-" + strconv.Itoa(maxSsd) + ", steps of " + strconv.Itoa(stepSsd) + ")")
	}
	return nil
}

func (d *Driver) getAPI() *oaocs.API {
	return oaocs.New(d.AccessToken, d.Endpoint)
}
