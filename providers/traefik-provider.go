package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/containous/traefik/v2/pkg/rules"
	log "github.com/sirupsen/logrus"
)

////----- Constants --->
const (
	traefikRoutersPath string = "api/http/routers"
)

////----- end

////----- Models --->
type Traefik struct {
	URL     string
	User    string
	Pass    string
	Dockers map[string]*Docker
	Ignore  []string
}
type TraefikConfig struct {
	URL     string   `json:"url"`
	Dockers []string `json:"dockers"`
	Ignored []string `json:"ignored"`
	User    string   `json:"user"`
	Pass    string   `json:"pass"`
}
type TraefikRouter struct {
	Name    string `json:"service"`
	RuleStr string `json:"rule"`
	Status  string `json:"status"`
	Enabled bool
}
type TraefikVersion struct {
	Version string `json:"Version"`
	Codename string `json:"Codename"`
}
////----- end


////----- Common App Provider Functions --->
func newTraefik(name string) (*Traefik, error) {
	log.Debugf("creating new traefik provider")
	config, err := loadTraefikConfig(name)
	if err != nil {
		return nil, err
	}
	traefik := Traefik{
		User: config.User,
		Pass: config.Pass,
		URL: config.URL,
		Ignore: config.Ignored,
		Dockers: make(map[string]*Docker),
	}
	if !traefik.TestConnection(true) {
		return nil, fmt.Errorf("could not create traefik connection | %s", name)
	}
	for _, dkr := range config.Dockers {
		traefik.Dockers[dkr], err = newDocker(dkr)
		if err != nil {
			return nil, fmt.Errorf("could not add dockers to traefik provider | %s", name)
		}
	}
	return &traefik, nil
}

func toTraefik(cnf interface{}) (*Traefik, error) {
	var traefik *Traefik
	traefik, valid := cnf.(*Traefik)
	if !valid {
		return nil, fmt.Errorf("conversion to traefik type not valid")
	}
	return traefik, nil
}

// GetApps fetches all the routers visible fromthe provided traefik instance
// It will however not provide any apps back where the enabled flag was false
// if provided with a docker connection
// Returns a map of App formatted routers
func (tr *Traefik) GetApps() map[string]*App {

	routers, err := tr.getRouterList()
	if err != nil {
		return nil
	}

	apps := make(map[string]*App)
	for _, router := range routers {
		app := newApp()
		name := router.Name
		for _, igname := range tr.Ignore {
			if strings.ToLower(name) == strings.ToLower(igname) {
				app.Enabled = false
			}
		}
		app.Icon = getDefaultIcon(name)
		if router.RuleStr != "" {
			domainStrs, err := rules.ParseDomains(router.RuleStr)
			if err != nil {
				log.Errorf("failed to parse traefik router domain")
				continue
			}
			if len(domainStrs) > 0 {
				app.URL = domainStrs[0]
			}
		}
		for _, dkr := range tr.Dockers {
			newName, upName := dkr.UpgradeApp(name, app)
			if upName {
				name = newName
			}
		}
		if app.Enabled {
			apps[name] = app
		}
	}

	return apps
}

// TestConnection checks to see if traefik can be communicated with via
// the given url. This is done by checking the version
// Returns true if communication was possible
func (tr *Traefik) TestConnection(output bool) bool {
	version, err := tr.getTraefikVersion()
	if err != nil {
		return false
	}
	log.Debugf("traefik version checked | %s", version.Version)
	return true
}

func loadTraefikConfig(name string) (*TraefikConfig, error) {
	configPath := fmt.Sprintf("%s/%s.json", getFileConfigRoot(), name)
	trCnf := newTraefikConfig()
	cnfFile, err := os.Open(configPath)
	if err != nil {
		log.Errorf("no config file found | %s", configPath)
		return nil, err
	}
	defer cnfFile.Close()
	err = json.NewDecoder(cnfFile).Decode(trCnf)
	if err != nil {
		log.Errorf("config file could not be parsed | %s", configPath)
		return nil, err
	}
	if len(trCnf.URL) > 0 {
		if trCnf.URL[len(trCnf.URL)-1:] == "/" {
			trCnf.URL = trCnf.URL[:len(trCnf.URL)-1]
		}
	}
	return trCnf, nil
}
////----- end

////----- Provider Specific ---> eof
func newTraefikConfig() *TraefikConfig {
	return &TraefikConfig{}
}

func (tr *Traefik) requestFromTraefik(path string) (*http.Response, error) {
	//TODO: support non localhost hosts....
	path = fmt.Sprintf("%s/%s", tr.URL, path)
	req, err := http.NewRequest("GET", path, nil)
	if tr.User != "" {
		req.SetBasicAuth(tr.User, tr.Pass)
	}
	client := &http.Client{}
	response, err := client.Do(req)
	return response, err
}

func (tr *Traefik) getRouterList() ([]*TraefikRouter, error) {
	response, err := tr.requestFromTraefik(traefikRoutersPath)
	if err != nil || response.StatusCode != 200 {
		log.Errorf("failed to fetch traefik router list")
		return nil, err
	}
	var routerList []*TraefikRouter
	err = json.NewDecoder(response.Body).Decode(&routerList)
	if err != nil {
		log.Errorf("traefik router list could not be parsed")
		return nil, err
	}
	return routerList, nil
}

func (tr *Traefik) getTraefikVersion() (*TraefikVersion, error) {
	response, err := tr.requestFromTraefik("api/version")
	if err != nil || response.StatusCode != 200 {
		log.Errorf("failed to fetch traefik version")
		return nil, err
	}
	var versionInfo *TraefikVersion
	err = json.NewDecoder(response.Body).Decode(&versionInfo)
	if err != nil {
		log.Errorf("traefik version info could not be parsed\n")
		return nil, err
	}
	return versionInfo, nil
}