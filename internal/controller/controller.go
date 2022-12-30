package controller

import (
	"fmt"
	"mi-c2/internal/logging"
	"mi-c2/internal/model"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	clusters     map[string]*model.Cluster
	clusterMutex sync.Mutex
	log          logging.Logger
	client       *resty.Client
	projects     *model.AvailableProjects
)

func Init() {
	clusters = make(map[string]*model.Cluster)
	log = logging.Log.With().Str("module", "controller").Logger()
	client = resty.New()

	projects = &model.AvailableProjects{}

	client.R().
		SetResult(projects).
		SetHeader("Accept", "application/json").
		Get("https://warriorhq.archiveteam.org/projects.json")

	go func() {
		for {
			client.R().
				SetResult(projects).
				SetHeader("Accept", "application/json").
				Get("https://warriorhq.archiveteam.org/projects.json")
			log.Info().Msg("updated active projects")
			time.Sleep(5 * time.Minute)
		}
	}()
}

func UpdateClusterStatus(cluster *model.Cluster) {
	clusterMutex.Lock()
	defer clusterMutex.Unlock()

	clusters[cluster.Name] = cluster
	log.Info().Str("cluster", cluster.Name).Msg("cluster updated")
}

func RequestClusterChange(cause, cluster, project string, workers, concurrency int) {
	clusterMutex.Lock()
	defer clusterMutex.Unlock()

	clusters[cluster].Connection.Socket.WriteJSON(model.Cluster{
		Project:     project,
		Workers:     workers,
		Concurrency: concurrency,
		Cause:       cause,
	})
}

func RemoveCluster(name string) {
	clusterMutex.Lock()
	defer clusterMutex.Unlock()

	delete(clusters, name)
	log.Info().Str("cluster", name).Msg("cluster deleted")
}

func GetCluster(name string) (*model.Cluster, error) {
	clusterMutex.Lock()
	defer clusterMutex.Unlock()

	cluster, ok := clusters[name]
	if !ok {
		return nil, fmt.Errorf("cluster does not exist or has not coneccted recently")
	}

	return cluster, nil
}

func GetClusters() []*model.Cluster {
	clusterMutex.Lock()
	defer clusterMutex.Unlock()

	res := []*model.Cluster{}
	for _, cluster := range clusters {
		res = append(res, cluster)
	}

	return res
}

func GetProjects() *model.AvailableProjects {
	return projects
}
