package properties

type ServerConfig struct {
	Port int
}

type GitLabConfig struct {
	Host      string
	AuthToken string
}

func (config GitLabConfig) UrlProjects() string {
	return config.Host + "/api/v4/projects"
}

func (config GitLabConfig) UrlPipelines(id string) string {
	return config.Host + "/api/v4/projects/" + id + "/pipelines"
}

func (config GitLabConfig) UrlDeleteProject(id string) string {
	return config.Host + "/api/v4/projects/" + id
}

type ReqCreateRepo struct {
	Name, Path, Description, NamespaceId, ImportUrl string
}

type RespCreateRepo struct {
	Id                      int
	Name, Path, Description string
	PathWithNamespace       string `json:"path_with_namespace"`
	WebUrl                  string `json:"web_url"`
	SshUrlToRepo            string `json:"ssh_url_to_repo"`
	HttpUrlToRepo           string `json:"http_url_to_repo"`
	CreatorId               int    `json:"creator_id"`
}

type ReqDeleteRepo struct {
	Id string
}

type ReqListPipelines struct {
	ProjectId string
}

type RespListPipeline struct {
	Id        int
	Iid       int
	ProjectId int
	Status    string
}
