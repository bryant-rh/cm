package client

import (
	"github.com/bryant-rh/cm/internal/model"
	"time"
)

type SaRes struct {
	ID        int32     `json:"id"`
	SaID      string    `json:"sa_id"`
	SaName    string    `json:"sa_name"`
	SaToken   string    `json:"sa_token"`
	NameSpace string    `json:"namespace"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type saData struct {
	GlobalRes
	Data string `json:"data"`
}

type saList struct {
	GlobalRes
	Data []SaRes `json:"data"`
}

type nsList struct {
	GlobalRes
	Data []model.Namespace `json:"data"`
}

//Sa_List /sa/list
func (c *CMClient) Sa_List(project_name, cluster_name, sa_name string) (*saList, error) {
	res := &saList{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"cluster_name": cluster_name,
			"sa_name":      sa_name,
		}).
		SetResult(res).Get("/sa/list")
	return res, err
}

//Sa_ListNs /sa/listns
func (c *CMClient) Sa_ListNs(project_name, cluster_name, sa_name string) (*nsList, error) {
	res := &nsList{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"cluster_name": cluster_name,
			"sa_name":      sa_name,
		}).
		SetResult(res).Get("/sa/listns")
	return res, err
}

//Sa_GetToken /sa/gettoken
func (c *CMClient) Sa_GetToken(project_name, cluster_name, ns_name string) (*saData, error) {
	res := &saData{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"cluster_name": cluster_name,
			"ns_name":      ns_name,
		}).
		SetResult(res).Get("/sa/gettoken")
	return res, err
}

type CreateSaRequestBody struct {
	ProjectName string `json:"project_name,omitempty"`
	ClusterName string `json:"cluster_name,omitempty"`
	SaName      string `json:"sa_name,omitempty"`
	SaToken     string `json:"sa_token,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
}

//Sa_Create /sa/create
func (c *CMClient) Sa_Create(project_name, cluster_name, sa_name, sa_token, namespace string) (*saList, error) {
	res := &saList{}
	sa_body := &CreateSaRequestBody{ProjectName: project_name, ClusterName: cluster_name, SaName: sa_name, SaToken: sa_token, Namespace: namespace}

	_, err := c.R().
		SetBody(sa_body).
		SetResult(res).Post("/sa/create")
	return res, err
}

type UpdateSaRequestBody struct {
	SaToken string `json:"sa_token,omitempty"`
}

//Sa_Update /sa/update
func (c *CMClient) Sa_Update(project_name, cluster_name, sa_name, sa_token string) (*saData, error) {
	res := &saData{}
	sa_body := &UpdateSaRequestBody{SaToken: sa_token}

	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"cluster_name": cluster_name,
			"sa_name":      sa_name,
		}).
		SetBody(sa_body).
		SetResult(res).Put("/sa/update")
	return res, err
}

//Sa_Delete /sa/delete
func (c *CMClient) Sa_Delete(project_name, cluster_name, sa_name string) (*saData, error) {
	res := &saData{}

	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"cluster_name": cluster_name,
			"sa_name":      sa_name,
		}).
		SetResult(res).Delete("/sa/delete")
	return res, err
}

type AddSaNsRequestBody struct {
	ProjectName string `json:"project_name,omitempty"`
	ClusterName string `json:"cluster_name,omitempty"`
	SaName      string `json:"sa_name,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
}

//Sa_AddNs /sa/adns
func (c *CMClient) Sa_AddNs(project_name, cluster_name, sa_name, namespace string) (*saList, error) {
	res := &saList{}
	sa_body := &AddSaNsRequestBody{ProjectName: project_name, ClusterName: cluster_name, SaName: sa_name, Namespace: namespace}

	_, err := c.R().
		SetBody(sa_body).
		SetResult(res).Post("/sa/addns")
	return res, err
}

//Sa_DelNs /sa/delns
func (c *CMClient) Sa_DelNs(project_name, cluster_name, sa_name, namespace string) (*saData, error) {
	res := &saData{}

	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"cluster_name": cluster_name,
			"sa_name":      sa_name,
			"namespace":      namespace,
		}).
		SetResult(res).Delete("/sa/delns")
	return res, err
}
