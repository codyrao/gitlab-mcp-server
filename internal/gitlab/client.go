package gitlab

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type GitLabClient struct {
	client *gitlab.Client
	host   string
}

func NewGitLabClient(host, token string) (*GitLabClient, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(host))
	if err != nil {
		return nil, fmt.Errorf("failed to create gitlab client: %w", err)
	}

	return &GitLabClient{
		client: client,
		host:   host,
	}, nil
}

func (c *GitLabClient) GetClient() *gitlab.Client {
	return c.client
}

func (c *GitLabClient) GetHost() string {
	return c.host
}

func (c *GitLabClient) ListProjects(opts *gitlab.ListProjectsOptions) ([]*gitlab.Project, error) {
	projects, _, err := c.client.Projects.ListProjects(opts)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (c *GitLabClient) GetProject(projectID interface{}) (*gitlab.Project, error) {
	project, _, err := c.client.Projects.GetProject(projectID, &gitlab.GetProjectOptions{})
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (c *GitLabClient) ListMergeRequests(projectID interface{}, opts *gitlab.ListProjectMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
	mrs, _, err := c.client.MergeRequests.ListProjectMergeRequests(projectID, opts)
	if err != nil {
		return nil, err
	}
	return mrs, nil
}

func (c *GitLabClient) GetMergeRequest(projectID interface{}, mrIID int) (*gitlab.MergeRequest, error) {
	mr, _, err := c.client.MergeRequests.GetMergeRequest(projectID, mrIID, &gitlab.GetMergeRequestsOptions{})
	if err != nil {
		return nil, err
	}
	return mr, nil
}

func (c *GitLabClient) ListMergeRequestNotes(projectID interface{}, mrIID int, opts *gitlab.ListMergeRequestNotesOptions) ([]*gitlab.Note, error) {
	notes, _, err := c.client.Notes.ListMergeRequestNotes(projectID, mrIID, opts)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (c *GitLabClient) CreateMergeRequestNote(projectID interface{}, mrIID int, opts *gitlab.CreateMergeRequestNoteOptions) (*gitlab.Note, error) {
	note, _, err := c.client.Notes.CreateMergeRequestNote(projectID, mrIID, opts)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (c *GitLabClient) GetFile(projectID, filePath, ref string) ([]byte, error) {
	file, _, err := c.client.RepositoryFiles.GetRawFile(projectID, filePath, &gitlab.GetRawFileOptions{Ref: &ref})
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (c *GitLabClient) GetFileInfo(projectID, filePath, ref string) (*gitlab.File, error) {
	file, _, err := c.client.RepositoryFiles.GetFile(projectID, filePath, &gitlab.GetFileOptions{Ref: &ref})
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (c *GitLabClient) ListCommits(projectID interface{}, opts *gitlab.ListCommitsOptions) ([]*gitlab.Commit, error) {
	commits, _, err := c.client.Commits.ListCommits(projectID, opts)
	if err != nil {
		return nil, err
	}
	return commits, nil
}

func (c *GitLabClient) ListPipelines(projectID interface{}, opts *gitlab.ListProjectPipelinesOptions) ([]*gitlab.PipelineInfo, error) {
	pipelines, _, err := c.client.Pipelines.ListProjectPipelines(projectID, opts)
	if err != nil {
		return nil, err
	}
	return pipelines, nil
}

func (c *GitLabClient) GetPipeline(projectID interface{}, pipelineID int) (*gitlab.Pipeline, error) {
	pipeline, _, err := c.client.Pipelines.GetPipeline(projectID, pipelineID)
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

func (c *GitLabClient) CreateMergeRequest(projectID interface{}, opts *gitlab.CreateMergeRequestOptions) (*gitlab.MergeRequest, error) {
	mr, _, err := c.client.MergeRequests.CreateMergeRequest(projectID, opts)
	if err != nil {
		return nil, err
	}
	return mr, nil
}

func (c *GitLabClient) ListProjectsOwned(opts *gitlab.ListProjectsOptions) ([]*gitlab.Project, error) {
	projects, _, err := c.client.Projects.ListProjects(&gitlab.ListProjectsOptions{
		Membership: gitlab.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (c *GitLabClient) SearchProjects(query string, opts *gitlab.ListProjectsOptions) ([]*gitlab.Project, error) {
	projects, _, err := c.client.Projects.ListProjects(&gitlab.ListProjectsOptions{
		Search: &query,
	})
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (c *GitLabClient) ListIssues(projectID interface{}, opts *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
	issues, _, err := c.client.Issues.ListIssues(&gitlab.ListIssuesOptions{})
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func (c *GitLabClient) GetIssue(projectID interface{}, issueID int) (*gitlab.Issue, error) {
	issue, _, err := c.client.Issues.GetIssue(projectID, issueID)
	if err != nil {
		return nil, err
	}
	return issue, nil
}

func (c *GitLabClient) CreateIssue(projectID interface{}, opts *gitlab.CreateIssueOptions) (*gitlab.Issue, error) {
	issue, _, err := c.client.Issues.CreateIssue(projectID, opts)
	if err != nil {
		return nil, err
	}
	return issue, nil
}

func (c *GitLabClient) ListGroupsOwned(opts *gitlab.ListGroupsOptions) ([]*gitlab.Group, error) {
	groups, _, err := c.client.Groups.ListGroups(opts)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (c *GitLabClient) GetGroup(groupID string) (*gitlab.Group, error) {
	group, _, err := c.client.Groups.GetGroup(groupID, &gitlab.GetGroupOptions{})
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (c *GitLabClient) ListGroupProjects(groupID string, opts *gitlab.ListGroupProjectsOptions) ([]*gitlab.Project, error) {
	projects, _, err := c.client.Groups.ListGroupProjects(groupID, opts)
	if err != nil {
		return nil, err
	}
	return projects, nil
}
