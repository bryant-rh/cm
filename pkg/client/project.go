package client

import (
	"github.com/bryant-rh/cm/internal/model"
)

type projectData struct {
	GlobalRes
	Data string `json:"data"`
}

type projectList struct {
	GlobalRes
	Data []model.Project `json:"data"`
}

type projectOne struct {
	GlobalRes
	Data model.Project `json:"data"`
}

//Project_List /project/list
func (c *CMClient) Project_List(project_name string) (*projectList, error) {
	res := &projectList{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
		}).
		SetResult(res).Get("/project/list")
	return res, err
}

//Project_GetId /project/get_id
func (c *CMClient) Project_GetId(project_name string) (*projectData, error) {
	res := &projectData{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
		}).
		SetResult(res).Get("/project/get_id")
	return res, err
}

type ProjectRequestBody struct {
	ProjectName string `json:"project_name,omitempty"`
	Description string `json:"description,omitempty"`
}

//Project_Create /project/create
func (c *CMClient) Project_Create(project_name, description string) (*projectOne, error) {
	res := &projectOne{}
	project_body := &ProjectRequestBody{ProjectName: project_name, Description: description}

	_, err := c.R().
		SetBody(project_body).
		SetResult(res).Post("/project/create")
	return res, err
}



//Project_Update /project/update
func (c *CMClient) Project_Update(project_id, new_project_name, new_description string) (*projectData, error) {
	res := &projectData{}
	project_body := &ProjectRequestBody{ProjectName: new_project_name, Description: new_description}

	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_id": project_id,
		}).
		SetBody(project_body).
		SetResult(res).Put("/project/update")
	return res, err
}

//Project_Delete /project/delete
func (c *CMClient) Project_Delete(project_name string) (*projectData, error) {
	res := &projectData{}

	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
		}).
		SetResult(res).Delete("/project/delete")
	return res, err
}
