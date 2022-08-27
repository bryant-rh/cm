package client

var (
	ReplicasDefault        = 2
	DataTunnelCountDefault = 50
	RetryIntervalDefault   = "5s"
)

type ProxyRequestBody struct {
	Image           string `json:"image"`
	Cluster         string `json:"cluster"`
	GatewayUrl      string `json:"gatewayUrl"`
	Namespace       string `json:"namespace,omitempty"`
	Replicas        int    `json:"replicas,omitempty"`
	RetryInterval   string `json:"retryInterval,omitempty"`
	DataTunnelCount int    `json:"dataTunnelCount,omitempty"`
}

//Proxy_Create /console/proxies
func (c *CMClient) Proxy_Create(image, cluster, gatewayUrl, namespace string) (string, error) {
	project_body := &ProxyRequestBody{
		Image:           image,
		Cluster:         cluster,
		GatewayUrl:      gatewayUrl,
		Namespace:       namespace,
		Replicas:        ReplicasDefault,
		RetryInterval:   RetryIntervalDefault,
		DataTunnelCount: DataTunnelCountDefault,
	}

	resp, err := c.R().
		SetBody(project_body).
		Post("/console/proxies")
	return resp.String(), err
}
