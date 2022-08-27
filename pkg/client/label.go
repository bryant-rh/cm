package client

type labelData struct {
	GlobalRes
	Data string `json:"data"`
}

type CluseterLabel struct {
	ProjectName string `json:"project_name"`
	ClusterName string `json:"cluster_name"`
	Labels      string `json:"Labels"`
}

type labelList struct {
	GlobalRes
	Data []CluseterLabel `json:"data"`
}

//Label_List /label/list
func (c *CMClient) Label_List(project_name, cluster_name string) (*labelList, error) {
	res := &labelList{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"cluster_name": cluster_name,
		}).
		SetResult(res).Get("/label/list")
	return res, err
}

type CreateLabelRequestBody struct {
	ProjectName string `json:"project_name,omitempty"`
	ClusterName string `json:"cluster_name,omitempty"`
	LabelKey    string `json:"label_key,omitempty"`
	LabelValue  string `json:"label_value,omitempty"`
}

//Label_Create /label/create
func (c *CMClient) Label_Create(project_name, cluster_name, label_key, label_value string) (*labelList, error) {
	res := &labelList{}
	label_body := &CreateLabelRequestBody{ProjectName: project_name, ClusterName: cluster_name, LabelKey: label_key, LabelValue: label_value}

	_, err := c.R().
		SetBody(label_body).
		SetResult(res).Post("/label/create")
	return res, err
}

//Label_Delete /label/delete
func (c *CMClient) Label_Delete(project_name, cluster_name, label_key, label_value string) (*labelData, error) {
	res := &labelData{}

	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"cluster_name": cluster_name,
			"label_key":    label_key,
			"label_value":  label_value,
		}).
		SetResult(res).Delete("/label/delete")
	return res, err
}
