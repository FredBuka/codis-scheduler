package handle

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8SClient struct {
	ClientSet *kubernetes.Clientset
}

var cli = new(K8SClient)
var replicas = kingpin.Flag("server-replicas", "num of per codis-server group").Default("2").Int()

func init() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	if config, err := rest.InClusterConfig(); err != nil {
		panic(err)
	} else if ClientSet, err := kubernetes.NewForConfig(config); err != nil {
		panic(err)
	} else {
		cli.ClientSet = ClientSet
	}
}

// getNodeNameByPodName return nodeName according to podName
func (c *K8SClient) getNodeNameByPodName(namespace, podName string) (error, string) {
	nodeName := ""
	resp, err := c.ClientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return err, nodeName
	}

	for _, item := range resp.Items {
		if item.Name == podName {
			nodeName = item.Spec.NodeName
		}
	}
	return nil, nodeName
}

// isAvailableSchedulerNode check if the current node can be scheduled.
// if the peer pods is scheduled on current node
func (c *K8SClient) isAvailableSchedulerNode(namespace, noSchedulerPodName, curNode string) (bool, error) {
	podNumber := strings.Split(noSchedulerPodName, "-")
	if len(podNumber) == 0 {
		fmt.Println("podName is:", noSchedulerPodName)
		return true, nil
	}

	podSeqInt, _ := strconv.Atoi(podNumber[len(podNumber)-1])
	prefix := strings.Join(podNumber[:len(podNumber)-1], "-")

	peerPodSeqArr := c.getPeerPodNameSerialNumber(podSeqInt)
	for _, serial := range peerPodSeqArr {
		peerPod := fmt.Sprintf("%v-%v", prefix, serial)
		err, peerNode := c.getNodeNameByPodName(namespace, peerPod)
		if err != nil {
			fmt.Println("Get node of peer pod fail", peerPod, err)
			continue
		}

		if peerNode == curNode {
			return false, fmt.Errorf("peer pod: %s Already scheduled on this node", peerPod)
		}
	}

	return true, nil
}

// getPeerPodNameSerialNumber return serial number of other members in the same group
func (c *K8SClient) getPeerPodNameSerialNumber(curSeqNum int) (nums []int) {
	serveReplicas := *replicas
	gid := curSeqNum/serveReplicas + 1
	for i := 0; i < serveReplicas; i++ {
		peer := (gid-1)*serveReplicas + i
		if peer != curSeqNum {
			nums = append(nums, peer)
		}
	}
	return
}
