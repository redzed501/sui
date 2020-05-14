package providers

import (
	"encoding/json"
	"fmt"
	"net/http"

	rules "github.com/containous/traefik/v2/pkg/rules"
	log "github.com/sirupsen/logrus"
	"github.com/willfantom/sui/config"
)

type TraefikProvider struct {
	URL    string
	User   string
	Pass   string
	Docker *DockerProvider
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

	ap.TypeConfig = &tp
	return ap, nil
}

func NewDockerTraefikProvider(cnf *config.TraefikConfig, dkr *AppProvider) (*AppProvider, error) {
	app, err := NewTraefikProvider(cnf)
	app.TypeConfig.(*TraefikProvider).Docker = dkr.TypeConfig.(*DockerProvider)
	return app, err
}

func (tp *TraefikProvider) GetApps(list map[string]*App) error {

	_, err := tp.fetchRouters()
	if err != nil {
		return err
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