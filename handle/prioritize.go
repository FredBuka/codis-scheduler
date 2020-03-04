package handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	schedulerApi "k8s.io/kubernetes/pkg/scheduler/api"
)

// prioritizeFunc
func prioritizeFunc(args schedulerApi.ExtenderArgs) (*schedulerApi.HostPriorityList, error) {
	nodes := args.Nodes.Items
	var priorityList schedulerApi.HostPriorityList
	priorityList = make([]schedulerApi.HostPriority, len(nodes))
	for i, node := range nodes {
		priorityList[i] = schedulerApi.HostPriority{
			Host:  node.Name,
			Score: 0,
		}
	}
	fmt.Println(args.Pod.Name, "prioritize phase...")
	return &priorityList, nil
}

// PrioritizeHandler
func PrioritizeHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	log.Println("Prioritize request info", buf.String())
	var extenderArgs schedulerApi.ExtenderArgs
	var hostPriorityList *schedulerApi.HostPriorityList
	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		panic(err)
	}
	if list, err := prioritizeFunc(extenderArgs); err != nil {
		panic(err)
	} else {
		hostPriorityList = list
	}
	if resultBody, err := json.Marshal(hostPriorityList); err != nil {
		panic(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resultBody)
	}
}
