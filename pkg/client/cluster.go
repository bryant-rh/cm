package client

import (
	"github.com/bryant-rh/cm/internal/model"
)

type clusterData struct {
	GlobalRes
	Data string `json:"data"`
}

type clusterList struct {
	GlobalRes
	Data []model.Cluster `json:"data"`
}

type clusterOne struct {
	GlobalRes
	Data model.Cluster `json:"data"`
}

//Cluster_List /cluster/list
func (c *CMClient) Cluster_List(project_name, cluster_name string) (*clusterList, error) {
	res := &clusterList{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"cluster_name": cluster_name,
		}).
		SetResult(res).Get("/cluster/list")
	return res, err
}

//Cluster_GetId /cluster/get_id
func (c *CMClient) Cluster_GetId(project_name, cluster_name string) (*clusterData, error) {
	res := &clusterData{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"cluster_name": cluster_name,
		}).
		SetResult(res).Get("/cluster/get_id")
	return res, err
}

//Cluster_Label /cluster/label
func (c *CMClient) Cluster_label(project_name, label_key, label_value string) (*clusterList, error) {
	res := &clusterList{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"label_key": label_key,
			"label_value": label_value,
		}).
		SetResult(res).Get("/cluster/label")
	return res, err
}


type CreateClusterRequestBody struct {
	ProjectName string `json:"project_name,omitempty"`
	ClusterName string `json:"cluster_name,omitempty"`
	Description string `json:"description,omitempty"`
}

//Cluster_Create /cluster/create
func (c *CMClient) Cluster_Create(project_name, cluster_name, description string) (*clusterOne, error) {
	res := &clusterOne{}
	cluster_body := &CreateClusterRequestBody{ProjectName: project_name, ClusterName: cluster_name, Description: description}

	_, err := c.R().
		SetBody(cluster_body).
		SetResult(res).Post("/cluster/create")
	return res, err
}

type UpdateClusterRequestBody struct {
	ClusterName string `json:"cluster_name,omitempty"`
	Description string `json:"description,omitempty"`
}

//Cluster_Update /cluster/update
func (c *CMClient) Cluster_Update(cluster_id, project_name, new_cluster_name, new_description string) (*clusterData, error) {
	res := &clusterData{}
	cluster_body := &UpdateClusterRequestBody{ClusterName: new_cluster_name, Description: new_description}

	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"cluster_id":   cluster_id,
		}).
		SetBody(cluster_body).
		SetResult(res).Put("/cluster/update")
	return res, err
}

//Cluster_Delete /cluster/delete
func (c *CMClient) Cluster_Delete(project_name, cluster_name string) (*clusterData, error) {
	res := &clusterData{}

	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"cluster_name": cluster_name,
		}).
		SetResult(res).Delete("/cluster/delete")
	return res, err
}
