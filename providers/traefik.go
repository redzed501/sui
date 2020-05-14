package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	rules "github.com/containous/traefik/v2/pkg/rules"
	log "github.com/sirupsen/logrus"
	"github.com/willfantom/sui/config"
)

type TraefikProvider struct {
	URL    string
	User   string
	Pass   string
	Docker *DockerProvider
	Ignore []string
}

type TraefikRouter struct {
	Name    string `json:"service"`
	RuleStr string `json:"rule"`
	Domain  string
	Status  string `json:"status"`
}

func NewTraefikProvider(cnf *config.TraefikConfig) (*AppProvider, error) {
	ap := newAppProvider()
	ap.PType = Traefik

	var tp TraefikProvider
	tp.URL = cnf.URL
	tp.User = cnf.User
	tp.Pass = cnf.Pass
	tp.Docker = nil
	tp.Ignore = cnf.IgnoredList

	ap.TypeConfig = &tp
	return ap, nil
}

func NewDockerTraefikProvider(cnf *config.TraefikConfig, dkr *AppProvider) (*AppProvider, error) {
	app, err := NewTraefikProvider(cnf)
	app.TypeConfig.(*TraefikProvider).Docker = dkr.TypeConfig.(*DockerProvider)
	return app, err
}

func (tp *TraefikProvider) GetApps(list map[string]*App) error {

	routers, err := tp.fetchRouters()
	if err != nil {
		return err
	}

	for _, router := range routers {

		app := newApp()
		name := router.Name
		if tp.Docker != nil {
			ci, err := tp.Docker.GetLocalContainerInfo(router.Name)
			if err == nil {
				dName, exist := ci.Labels[nameFromLabel]
				if exist {
					name = dName
				}
			}
		}
		for _, ignoreName := range tp.Ignore {
			if strings.ToLower(name) == ignoreName {
				continue
			}
		}
		app.URL = router.Domain
		app.Icon = "application"
		app.Protected = false
		if tp.Docker != nil {
			ci, err := tp.Docker.GetLocalContainerInfo(router.Name)
			if err == nil {
				dURL, exist := ci.Labels[urlFromLabel]
				if exist {
					app.URL = dURL
				}
				dIcon, exist := ci.Labels[iconFromLabel]
				if exist {
					app.Icon = dIcon
				}
				dProtec, exist := ci.Labels[iconFromLabel]
				if exist {
					dProtecBool, err := strconv.ParseBool(dProtec)
					if err == nil {
						app.Protected = dProtecBool
					}
				}
			}
		}
		list[name] = app
	}

	return nil
}

func (tp *TraefikProvider) fetchRouters() ([]*TraefikRouter, error) {
	client := &http.Client{}
	path := fmt.Sprintf("%s/api/http/routers", tp.URL)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(tp.User, tp.Pass)
	response, err := client.Do(req)
	if err != nil {
		log.Errorf("failed to fetch traefik router info")
		return nil, err
	}
	var routerList []*TraefikRouter
	err = json.NewDecoder(response.Body).Decode(&routerList)
	if err != nil {
		log.Errorf("traefik router list could not be parsed")
		return nil, err
	}

	for idx, router := range routerList {
		if router.RuleStr != "" {
			domainStrs, err := rules.ParseDomains(router.RuleStr)
			if err != nil {
				log.Errorf("failed to parse traefik router domain")
				return nil, err
			}
			if len(domainStrs) > 0 {
				routerList[idx].Domain = domainStrs[0]
			}
			
		}
	}

	return routerList, nil
}