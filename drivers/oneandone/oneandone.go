package oneandone

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/machine/drivers"
	"github.com/docker/machine/provider"
	"github.com/docker/machine/state"
	gossh "golang.org/x/crypto/ssh"
	"path/filepath"
	"strconv"
	oaocs "github.com/jlusiardi/oneandone-cloudserver-api"
)

const (
	minCores          = 1
	maxCores          = 16
	minRam            = 1
	maxRam            = 128
	minSsd            = 20
	maxSsd            = 500
	stepSsd           = 20
)

const Endpoint string = ""

type Driver struct {
	AccessToken    string
	VmId           string
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
	log.Infof("Creating a new 1&1 CloudServer ...")
	api := oaocs.New(d.AccessToken, Endpoint)
		
	in := oaocs.ServerCreateData{}
	in.Name = d.MachineName
	in.Description = d.MachineName + " created by docker machine"
	in.ApplianceId = "C14988A9ABC34EA64CD5AAC0D33ABCAF"
	hw := oaocs.Hardware{}
	in.Hardware = hw
	hw.CoresPerProcessor = 1
	hw.Vcores = d.Cores
	hw.Ram = d.Ram
	hdd := oaocs.Hdd{}
	hdd.IsMain = true
	hdd.Size = d.Ssd
	hw.Hdds = []oaocs.Hdd{hdd}
	/*server := */api.CreateServer(in)

	return nil
}

func (d *Driver) GetIP() (string, error) {
	return d.IPAddress, nil
}

func (d *Driver) GetMachineName() string {
	return d.MachineName
}

func (d *Driver) GetProviderType() provider.ProviderType {
	return provider.Remote
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
	return state.None, nil
}

func (d *Driver) PreCreateCheck() error {
	return nil
}

func (d *Driver) GetURL() (string, error) {
	return fmt.Sprintf("tcp://%s:2376", d.IPAddress), nil
}

func (d *Driver) Kill() error {
	log.Infof("Killing the 1&1 CloudServer named '%s' ...", d.MachineName)
	return nil
}

func (d *Driver) Remove() error {
	log.Infof("Removing the 1&1 CloudServer named '%s' ...", d.MachineName)
	return nil
}

func (d *Driver) Start() error {
	log.Infof("Starting the 1&1 CloudServer named '%s' ...", d.MachineName)
	return nil
}

func (d *Driver) Stop() error {
	log.Infof("Stopping the 1&1 CloudServer named '%s' ...", d.MachineName)
	return nil
}

func (d *Driver) Restart() error {
	log.Infof("Restarting the 1&1 CloudServer named '%s' ...", d.MachineName)
	return nil
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.AccessToken = flags.String("oneandone-access-token")
	if d.AccessToken == "" {
		return fmt.Errorf("oneandone driver requires the --oneandone-access-token option")
	}
	d.Cores = flags.Int("oneandone-cores")
	if d.Cores < minCores || d.Cores > maxCores {
		return fmt.Errorf("oneandone driver requires the --oneandone-cores option to be an integer (" + strconv.Itoa(minCores) + "-" + strconv.Itoa(maxCores) + ")")
	}
	d.Ram = flags.Int("oneandone-ram")
	if d.Ram < minRam || d.Ram > maxRam {
		return fmt.Errorf("oneandone driver requires the --oneandone-ram option to be an integer (" + strconv.Itoa(minRam) + "-" + strconv.Itoa(maxRam) + ")")
	}
	d.Ssd = flags.Int("oneandone-ssd")
	if d.Ssd < minSsd || d.Ssd > maxSsd || (d.Ssd%stepSsd) != 0 {
		return fmt.Errorf("oneandone driver requires the --oneandone-ssd option to be an integer (" + strconv.Itoa(minSsd) + "-" + strconv.Itoa(maxSsd) + ", steps of " + strconv.Itoa(stepSsd) + ")")
	}
	return nil
}

func getClient(user string, ip string, password string) (*gossh.Client, error) {
	sshConfig := &gossh.ClientConfig{
		User: user,
		Auth: []gossh.AuthMethod{gossh.Password(password)},
	}
	client, err := gossh.Dial("tcp", ip+":22", sshConfig)
	if err != nil {
		return nil, err
	}
	return client, err
}

func executeCmd(client *gossh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		log.Info(err)
		return "", err
	}
	var b bytes.Buffer
	session.Stdout = &b
	err = session.Run(cmd)
	if err != nil {
		return "", err
	}
	return b.String(), err
}
