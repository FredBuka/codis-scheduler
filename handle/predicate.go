package handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"k8s.io/api/core/v1"
	schedulerApi "k8s.io/kubernetes/pkg/scheduler/api"
)

// predicateFunc
func predicateFunc(args schedulerApi.ExtenderArgs) *schedulerApi.ExtenderFilterResult {
	canSchedule := make([]v1.Node, 0, len(args.Nodes.Items))
	canNotSchedule := make(map[string]string)
	for _, node := range args.Nodes.Items {
		result, err := cli.isAvailableSchedulerNode(args.Pod.Namespace, args.Pod.Name, node.Name)
		if err != nil {
			canNotSchedule[node.Name] = err.Error()
			continue
		}

		if result {
			canSchedule = append(canSchedule, node)
			fmt.Println("node", args.Pod.Namespace, node.Name, node.Namespace)
		}

	}
	result := schedulerApi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: canSchedule,
		},
		FailedNodes: canNotSchedule,
		Error:       "",
	}

	fmt.Println(args.Pod.Name, "predicate phase...")
	return &result
}

// PredicateHandler
func PredicateHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	log.Println("Predicate request info", buf.String())

	var extenderArgs schedulerApi.ExtenderArgs
	var extenderFilterResult *schedulerApi.ExtenderFilterResult

	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		extenderFilterResult = &schedulerApi.ExtenderFilterResult{
			Nodes:       nil,
			FailedNodes: nil,
			Error:       err.Error(),
		}
	} else {
		extenderFilterResult = predicateFunc(extenderArgs)
	}
	if resultBody, err := json.Marshal(extenderFilterResult); err != nil {
		panic(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resultBody)
	}
}
