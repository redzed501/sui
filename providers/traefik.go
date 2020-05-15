package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	rules "github.com/containous/traefik/v2/pkg/rules"
	log "github.com/sirupsen/logrus"
	"github.com/willfantom/sui/config"
)

type TraefikProvider struct {
	URL    string
	User   string
	Pass   string
	Dockers map[string]*DockerProvider
	Ignore []string
}

type TraefikRouter struct {
	Name    string `json:"service"`
	RuleStr string `json:"rule"`
	Domain  string
	Status  string `json:"status"`
	Ignored bool
}

func NewTraefikProvider(cnf *config.TraefikConfig) (*AppProvider, error) {
	var err error
	ap := newAppProvider()
	ap.PType = Traefik

	var tp TraefikProvider
	tp.URL = cnf.URL
	tp.User = cnf.User
	tp.Pass = cnf.Pass
	tp.Ignore = cnf.IgnoredList
	tp.Dockers = make(map[string]*DockerProvider)

	for name, dp := range cnf.DockerConfigs {
		tp.Dockers[name], err = NewDockerProviderLite(dp)
		if err != nil {
			log.Errorf("failed to add docker provider (%s) to traefik provider", name)
		}
	}

	ap.TypeConfig = &tp
	return ap, nil
}

func (tp *TraefikProvider) GetApps(list map[string]*App) error {
	//Get routers from traefik api
	routers, err := tp.fetchRouters()
	if err != nil {
		return err
	}
	for _, router := range routers {
		if router.Ignored {
			continue
		}
		app := newApp()
		name := strings.ToLower(router.Name)
		app.URL = router.Domain
		app.Icon = getDefaultIcon(name)
		app.Protected = false	
		list[name] = app
	}
	//Update with docker label info
	// service name in traefik must match container's name
	for _, dp := range tp.Dockers {
		cil, err := dp.GetContainerList()
		if err != nil {
			continue
		}
		for _, ci := range cil {
			nameMatch := ci.Names[0][1:]
			_, match := list[nameMatch]
			if match {
				updated := list[nameMatch].UpdateFromDockerLabels(ci)
				if updated {
					log.Debugf("updated app (%s) with docker labels", nameMatch)
				}
			}
		}
	}


	/*
		if tp.Docker != nil {
			ci, err := tp.Docker.GetLocalContainerInfo(router.Name)
			if err == nil {
				dName, exist := ci.Labels[nameFromLabel]
				if exist {
					name = dName
				}
			}
		}
		toIgnore := false
		for _, ignoreName := range tp.Ignore {
			if strings.ToLower(name) == ignoreName {
				toIgnore = true
			}
		}
		if toIgnore {
			continue
		}
		app.URL = router.Domain
		app.Icon = getDefaultIcon(name)
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
	}*/

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
		for _, ignore := range tp.Ignore {
			if strings.ToLower(ignore) == strings.ToLower(router.Name) {
				router.Ignored = true
			}
		}
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