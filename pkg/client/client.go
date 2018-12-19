package prometheusclient

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"gitlab.300.cn/paas-k8s/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringclient "gitlab.300.cn/paas-k8s/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	apiVersionMonitoring = "monitoring.coreos.com/v1"
	kindServiceMonitor   = "ServiceMonitor"

	intevalEndpoints = "30s"
)

type Client struct {
	KubeClient     kubernetes.Interface
	MonClientV1    monitoringclient.MonitoringV1Interface
	HTTPClient     *http.Client
	MasterHost     string
	DefaultTimeout time.Duration
}

// New setups a test framework and returns it.
func New(kubeconfig string) (*Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "build config from flags failed")
	}

	cli, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "creating new kube-client failed")
	}

	httpc := cli.CoreV1().RESTClient().(*rest.RESTClient).Client
	if err != nil {
		return nil, errors.Wrap(err, "creating http-client failed")
	}

	mClientV1, err := monitoringclient.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "creating v1 monitoring client failed")
	}

	f := &Client{
		MasterHost:     config.Host,
		KubeClient:     cli,
		MonClientV1:    mClientV1,
		HTTPClient:     httpc,
		DefaultTimeout: time.Minute,
	}

	return f, nil
}

func (c *Client) GetServiceMonitor(namespace, name string) (result *v1.ServiceMonitor, err error) {
	return c.MonClientV1.ServiceMonitors(namespace).Get(name, metav1.GetOptions{})
}

func (c *Client) DeleteServiceMonitor(namespace, name string) error {
	return c.MonClientV1.ServiceMonitors(namespace).Delete(name, &metav1.DeleteOptions{})
}

func (c *Client) CreateServiceMonitor(namespace, name string) (result *v1.ServiceMonitor, err error) {
	endpoints := make([]v1.Endpoint, 0)
	ep := v1.Endpoint{
		// TODO
		Port:        fmt.Sprintf("%v-18888k", name),
		Interval:    intevalEndpoints,
		HonorLabels: true,
	}
	endpoints = append(endpoints, ep)

	matchExpressions := make([]metav1.LabelSelectorRequirement, 0)
	me := metav1.LabelSelectorRequirement{
		Key:    "matchNames",
		Values: []string{namespace},
	}
	matchExpressions = append(matchExpressions, me)

	sm := &v1.ServiceMonitor{
		TypeMeta: metav1.TypeMeta{
			Kind:       kindServiceMonitor,
			APIVersion: apiVersionMonitoring,
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels: map[string]string{
				// TODO
				"prometheus": "k8s",
				"kubeapp":    name,
			},
		},
		Spec: v1.ServiceMonitorSpec{
			JobLabel: "kubeapp",
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"kubeapp": name,
				},
			},
			NamespaceSelector: v1.NamespaceSelector{
				MatchNames: []string{namespace},
			},
			Endpoints: endpoints,
		},
	}

	return c.MonClientV1.ServiceMonitors(namespace).Create(sm)
}
