package mcp

import (
	"context"
	"fmt"

	"github.com/codyrao/gitlab-mcp-server/internal/security"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	gitlabapi "github.com/xanzy/go-gitlab"
)

type ListProjectsInput struct {
	Membership *bool   `json:"membership" jsonschema:"Limit to projects where the authenticated user is a member"`
	Owned      *bool   `json:"owned" jsonschema:"Limit to projects owned by the authenticated user"`
	Search     *string `json:"search" jsonschema:"Search projects by name"`
	PerPage    *int    `json:"per_page" jsonschema:"Number of results per page"`
}

type ListProjectsOutput struct {
	Projects []ProjectInfo `json:"projects"`
}

type ProjectInfo struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Path              string `json:"path"`
	PathWithNamespace string `json:"path_with_namespace"`
	Description       string `json:"description"`
	DefaultBranch     string `json:"default_branch"`
	Visibility        string `json:"visibility"`
	CreatedAt         string `json:"created_at"`
	LastActivityAt    string `json:"last_activity_at"`
}

func (s *Server) handleListProjects(ctx context.Context, req *mcpsdk.CallToolRequest, input ListProjectsInput) (
	*mcpsdk.CallToolResult,
	ListProjectsOutput,
	error,
) {
	result := s.validator.ValidateOperation(security.Operation{
		Type: security.OperationRead,
	})
	if !result.Allowed {
		return nil, ListProjectsOutput{}, fmt.Errorf(result.Reason)
	}

	opts := &gitlabapi.ListProjectsOptions{}
	if input.Membership != nil {
		opts.Membership = input.Membership
	}
	if input.Owned != nil {
		opts.Owned = input.Owned
	}
	if input.Search != nil {
		opts.Search = input.Search
	}
	if input.PerPage != nil {
		opts.PerPage = *input.PerPage
	}

	projects, err := s.client.ListProjects(opts)
	if err != nil {
		return nil, ListProjectsOutput{}, err
	}

	projectList := make([]ProjectInfo, len(projects))
	for i, p := range projects {
		projectList[i] = ProjectInfo{
			ID:                p.ID,
			Name:              p.Name,
			Path:              p.Path,
			PathWithNamespace: p.PathWithNamespace,
			Description:       p.Description,
			DefaultBranch:     p.DefaultBranch,
			Visibility:        string(p.Visibility),
			CreatedAt:         p.CreatedAt.Format("2006-01-02T15:04:05Z"),
			LastActivityAt:    p.LastActivityAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return nil, ListProjectsOutput{Projects: projectList}, nil
}

type GetProjectInput struct {
	ProjectID string `json:"project_id" jsonschema:"The project ID or path (e.g., 'group/project' or project ID number)"`
}

type GetProjectOutput struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Path              string   `json:"path"`
	PathWithNamespace string   `json:"path_with_namespace"`
	Description       string   `json:"description"`
	DefaultBranch     string   `json:"default_branch"`
	Visibility        string   `json:"visibility"`
	CreatedAt         string   `json:"created_at"`
	LastActivityAt    string   `json:"last_activity_at"`
	Archived          bool     `json:"archived"`
	TagList           []string `json:"tag_list"`
	StarCount         int      `json:"star_count"`
	ForksCount        int      `json:"forks_count"`
}

func (s *Server) handleGetProject(ctx context.Context, req *mcpsdk.CallToolRequest, input GetProjectInput) (
	*mcpsdk.CallToolResult,
	GetProjectOutput,
	error,
) {
	result := s.validator.ValidateOperation(security.Operation{
		Type: security.OperationRead,
	})
	if !result.Allowed {
		return nil, GetProjectOutput{}, fmt.Errorf(result.Reason)
	}

	if input.ProjectID == "" {
		return nil, GetProjectOutput{}, fmt.Errorf("project_id is required")
	}

	project, err := s.client.GetProject(input.ProjectID)
	if err != nil {
		return nil, GetProjectOutput{}, err
	}

	return nil, GetProjectOutput{
		ID:                project.ID,
		Name:              project.Name,
		Path:              project.Path,
		PathWithNamespace: project.PathWithNamespace,
		Description:       project.Description,
		DefaultBranch:     project.DefaultBranch,
		Visibility:        string(project.Visibility),
		CreatedAt:         project.CreatedAt.Format("2006-01-02T15:04:05Z"),
		LastActivityAt:    project.LastActivityAt.Format("2006-01-02T15:04:05Z"),
		Archived:          project.Archived,
		TagList:           project.TagList,
		StarCount:         project.StarCount,
		ForksCount:        project.ForksCount,
	}, nil
}

type ListMRsInput struct {
	ProjectID string  `json:"project_id" jsonschema:"The project ID or path"`
	State     *string `json:"state" jsonschema:"Filter by state (opened, closed, merged, all)"`
	PerPage   *int    `json:"per_page" jsonschema:"Number of results per page"`
}

type ListMRsOutput struct {
	MRs []MRInfo `json:"mrs"`
}

type MRInfo struct {
	IID          int    `json:"iid"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	State        string `json:"state"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	Author       string `json:"author"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func (s *Server) handleListMRs(ctx context.Context, req *mcpsdk.CallToolRequest, input ListMRsInput) (
	*mcpsdk.CallToolResult,
	ListMRsOutput,
	error,
) {
	result := s.validator.ValidateOperation(security.Operation{
		Type: security.OperationRead,
	})
	if !result.Allowed {
		return nil, ListMRsOutput{}, fmt.Errorf(result.Reason)
	}

	opts := &gitlabapi.ListProjectMergeRequestsOptions{}
	if input.State != nil {
		opts.State = input.State
	}
	if input.PerPage != nil {
		opts.PerPage = *input.PerPage
	}

	mrs, err := s.client.ListMergeRequests(input.ProjectID, opts)
	if err != nil {
		return nil, ListMRsOutput{}, err
	}

	mrList := make([]MRInfo, len(mrs))
	for i, mr := range mrs {
		author := ""
		if mr.Author != nil {
			author = mr.Author.Username
		}
		mrList[i] = MRInfo{
			IID:          mr.IID,
			Title:        mr.Title,
			Description:  mr.Description,
			State:        string(mr.State),
			SourceBranch: mr.SourceBranch,
			TargetBranch: mr.TargetBranch,
			Author:       author,
			CreatedAt:    mr.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:    mr.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return nil, ListMRsOutput{MRs: mrList}, nil
}

type GetMRDetailsInput struct {
	ProjectID string `json:"project_id" jsonschema:"The project ID or path"`
	MRIID     int    `json:"mr_iid" jsonschema:"The merge request IID"`
}

type GetMRDetailsOutput struct {
	IID          int    `json:"iid"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	State        string `json:"state"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	Author       string `json:"author"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	WebURL       string `json:"web_url"`
}

func (s *Server) handleGetMRDetails(ctx context.Context, req *mcpsdk.CallToolRequest, input GetMRDetailsInput) (
	*mcpsdk.CallToolResult,
	GetMRDetailsOutput,
	error,
) {
	result := s.validator.ValidateOperation(security.Operation{
		Type: security.OperationRead,
	})
	if !result.Allowed {
		return nil, GetMRDetailsOutput{}, fmt.Errorf(result.Reason)
	}

	mr, err := s.client.GetMergeRequest(input.ProjectID, input.MRIID)
	if err != nil {
		return nil, GetMRDetailsOutput{}, err
	}

	author := ""
	if mr.Author != nil {
		author = mr.Author.Username
	}

	return nil, GetMRDetailsOutput{
		IID:          mr.IID,
		Title:        mr.Title,
		Description:  mr.Description,
		State:        string(mr.State),
		SourceBranch: mr.SourceBranch,
		TargetBranch: mr.TargetBranch,
		Author:       author,
		CreatedAt:    mr.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    mr.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		WebURL:       mr.WebURL,
	}, nil
}

type CreateMRNoteInput struct {
	ProjectID string `json:"project_id" jsonschema:"The project ID or path"`
	MRIID     int    `json:"mr_iid" jsonschema:"The merge request IID"`
	Body      string `json:"body" jsonschema:"The comment text"`
}

type CreateMRNoteOutput struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	Author    string `json:"author"`
}

func (s *Server) handleCreateMRNote(ctx context.Context, req *mcpsdk.CallToolRequest, input CreateMRNoteInput) (
	*mcpsdk.CallToolResult,
	CreateMRNoteOutput,
	error,
) {
	result := s.validator.ValidateOperation(security.Operation{
		Type: security.OperationCreate,
	})
	if !result.Allowed {
		return nil, CreateMRNoteOutput{}, fmt.Errorf(result.Reason)
	}

	note, err := s.client.CreateMergeRequestNote(input.ProjectID, input.MRIID, &gitlabapi.CreateMergeRequestNoteOptions{
		Body: &input.Body,
	})
	if err != nil {
		return nil, CreateMRNoteOutput{}, err
	}

	author := ""
	if note.Author.ID != 0 {
		author = note.Author.Username
	}

	return nil, CreateMRNoteOutput{
		ID:        note.ID,
		Body:      note.Body,
		CreatedAt: note.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Author:    author,
	}, nil
}

type GetFileContentInput struct {
	ProjectID string `json:"project_id" jsonschema:"The project ID or path"`
	FilePath  string `json:"file_path" jsonschema:"Path to the file in the repository"`
	Ref       string `json:"ref" jsonschema:"The branch, tag, or commit reference"`
}

type GetFileContentOutput struct {
	Content string `json:"content"`
	Size    int64  `json:"size"`
}

func (s *Server) handleGetFileContent(ctx context.Context, req *mcpsdk.CallToolRequest, input GetFileContentInput) (
	*mcpsdk.CallToolResult,
	GetFileContentOutput,
	error,
) {
	result := s.validator.ValidateOperation(security.Operation{
		Type: security.OperationRead,
	})
	if !result.Allowed {
		return nil, GetFileContentOutput{}, fmt.Errorf(result.Reason)
	}

	file, err := s.client.GetFile(input.ProjectID, input.FilePath, input.Ref)
	if err != nil {
		return nil, GetFileContentOutput{}, err
	}

	return nil, GetFileContentOutput{
		Content: string(file),
		Size:    int64(len(file)),
	}, nil
}

type ListCommitsInput struct {
	ProjectID string  `json:"project_id" jsonschema:"The project ID or path"`
	RefName   *string `json:"ref_name" jsonschema:"The branch or tag name"`
	PerPage   *int    `json:"per_page" jsonschema:"Number of results per page"`
}

type ListCommitsOutput struct {
	Commits []CommitInfo `json:"commits"`
}

type CommitInfo struct {
	ID        string `json:"id"`
	ShortID   string `json:"short_id"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
}

func (s *Server) handleListCommits(ctx context.Context, req *mcpsdk.CallToolRequest, input ListCommitsInput) (
	*mcpsdk.CallToolResult,
	ListCommitsOutput,
	error,
) {
	result := s.validator.ValidateOperation(security.Operation{
		Type: security.OperationRead,
	})
	if !result.Allowed {
		return nil, ListCommitsOutput{}, fmt.Errorf(result.Reason)
	}

	opts := &gitlabapi.ListCommitsOptions{}
	if input.RefName != nil {
		opts.RefName = input.RefName
	}
	if input.PerPage != nil {
		opts.PerPage = *input.PerPage
	}

	commits, err := s.client.ListCommits(input.ProjectID, opts)
	if err != nil {
		return nil, ListCommitsOutput{}, err
	}

	commitList := make([]CommitInfo, len(commits))
	for i, c := range commits {
		commitList[i] = CommitInfo{
			ID:        c.ID,
			ShortID:   c.ShortID,
			Title:     c.Title,
			Message:   c.Message,
			Author:    c.AuthorName,
			CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return nil, ListCommitsOutput{Commits: commitList}, nil
}

type ListPipelinesInput struct {
	ProjectID string `json:"project_id" jsonschema:"The project ID or path"`
	PerPage   *int   `json:"per_page" jsonschema:"Number of results per page"`
}

type ListPipelinesOutput struct {
	Pipelines []PipelineInfo `json:"pipelines"`
}

type PipelineInfo struct {
	ID        int    `json:"id"`
	Ref       string `json:"ref"`
	SHA       string `json:"sha"`
	Status    string `json:"status"`
	Source    string `json:"source"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (s *Server) handleListPipelines(ctx context.Context, req *mcpsdk.CallToolRequest, input ListPipelinesInput) (
	*mcpsdk.CallToolResult,
	ListPipelinesOutput,
	error,
) {
	result := s.validator.ValidateOperation(security.Operation{
		Type: security.OperationRead,
	})
	if !result.Allowed {
		return nil, ListPipelinesOutput{}, fmt.Errorf(result.Reason)
	}

	opts := &gitlabapi.ListProjectPipelinesOptions{}
	if input.PerPage != nil {
		opts.PerPage = *input.PerPage
	}

	pipelines, err := s.client.ListPipelines(input.ProjectID, opts)
	if err != nil {
		return nil, ListPipelinesOutput{}, err
	}

	pipelineList := make([]PipelineInfo, len(pipelines))
	for i, p := range pipelines {
		pipelineList[i] = PipelineInfo{
			ID:        p.ID,
			Ref:       p.Ref,
			SHA:       p.SHA,
			Status:    string(p.Status),
			Source:    p.Source,
			CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return nil, ListPipelinesOutput{Pipelines: pipelineList}, nil
}

type CreateMRInput struct {
	ProjectID    string `json:"project_id" jsonschema:"The project ID or path"`
	SourceBranch string `json:"source_branch" jsonschema:"The source branch name"`
	TargetBranch string `json:"target_branch" jsonschema:"The target branch name (cannot be master or test)"`
	Title        string `json:"title" jsonschema:"The merge request title"`
	Description  string `json:"description" jsonschema:"The merge request description"`
}

type CreateMROutput struct {
	IID          int    `json:"iid"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	State        string `json:"state"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	WebURL       string `json:"web_url"`
}

func (s *Server) handleCreateMR(ctx context.Context, req *mcpsdk.CallToolRequest, input CreateMRInput) (
	*mcpsdk.CallToolResult,
	CreateMROutput,
	error,
) {
	validationResult := s.validator.CanCreateMR(input.SourceBranch, input.TargetBranch)
	if !validationResult.Allowed {
		return nil, CreateMROutput{}, fmt.Errorf(validationResult.Reason)
	}

	mr, err := s.client.CreateMergeRequest(input.ProjectID, &gitlabapi.CreateMergeRequestOptions{
		SourceBranch: &input.SourceBranch,
		TargetBranch: &input.TargetBranch,
		Title:        &input.Title,
		Description:  &input.Description,
	})
	if err != nil {
		return nil, CreateMROutput{}, err
	}

	return nil, CreateMROutput{
		IID:          mr.IID,
		Title:        mr.Title,
		Description:  mr.Description,
		State:        string(mr.State),
		SourceBranch: mr.SourceBranch,
		TargetBranch: mr.TargetBranch,
		WebURL:       mr.WebURL,
	}, nil
}

type ListIssuesInput struct {
	ProjectID string  `json:"project_id" jsonschema:"The project ID or path"`
	State     *string `json:"state" jsonschema:"Filter by state (opened, closed, all)"`
	PerPage   *int    `json:"per_page" jsonschema:"Number of results per page"`
}

type ListIssuesOutput struct {
	Issues []IssueInfo `json:"issues"`
}

type IssueInfo struct {
	IID         int    `json:"iid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
	Author      string `json:"author"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (s *Server) handleListIssues(ctx context.Context, req *mcpsdk.CallToolRequest, input ListIssuesInput) (
	*mcpsdk.CallToolResult,
	ListIssuesOutput,
	error,
) {
	result := s.validator.ValidateOperation(security.Operation{
		Type: security.OperationRead,
	})
	if !result.Allowed {
		return nil, ListIssuesOutput{}, fmt.Errorf(result.Reason)
	}

	opts := &gitlabapi.ListIssuesOptions{}
	if input.State != nil {
		opts.State = input.State
	}
	if input.PerPage != nil {
		opts.PerPage = *input.PerPage
	}

	issues, err := s.client.ListIssues(input.ProjectID, opts)
	if err != nil {
		return nil, ListIssuesOutput{}, err
	}

	issueList := make([]IssueInfo, len(issues))
	for i, issue := range issues {
		author := ""
		if issue.Author != nil {
			author = issue.Author.Username
		}
		issueList[i] = IssueInfo{
			IID:         issue.IID,
			Title:       issue.Title,
			Description: issue.Description,
			State:       string(issue.State),
			Author:      author,
			CreatedAt:   issue.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   issue.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return nil, ListIssuesOutput{Issues: issueList}, nil
}
