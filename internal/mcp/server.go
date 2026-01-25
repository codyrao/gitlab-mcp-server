package mcp

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/codyrao/gitlab-mcp-server/internal/config"
	"github.com/codyrao/gitlab-mcp-server/internal/gitlab"
	"github.com/codyrao/gitlab-mcp-server/internal/security"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	config    *config.Config
	client    *gitlab.GitLabClient
	validator *security.Validator
	mcpServer *mcpsdk.Server
}

func NewServer(cfg *config.Config, client *gitlab.GitLabClient) *Server {
	s := &Server{
		config:    cfg,
		client:    client,
		validator: security.NewValidator(),
	}
	s.createMCPServer()
	return s
}

func (s *Server) createMCPServer() {
	s.mcpServer = mcpsdk.NewServer(&mcpsdk.Implementation{
		Name:    "gitlab-mcp-server",
		Version: "1.0.0",
	}, nil)

	s.registerTools()
}

func (s *Server) registerTools() {
	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "gitlab_list_projects",
		Description: "List GitLab projects for authenticated user",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"membership": map[string]any{
					"type":        "boolean",
					"description": "Limit to projects where the authenticated user is a member",
				},
				"owned": map[string]any{
					"type":        "boolean",
					"description": "Limit to projects owned by the authenticated user",
				},
				"search": map[string]any{
					"type":        "string",
					"description": "Search projects by name",
				},
				"per_page": map[string]any{
					"type":        "integer",
					"description": "Number of results per page",
				},
			},
		},
	}, s.handleListProjects)

	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "gitlab_get_project",
		Description: "Get details of a GitLab project",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID or path (e.g., 'group/project' or project ID number)",
				},
			},
			"required": []string{"project_id"},
		},
	}, s.handleGetProject)

	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "gitlab_list_mrs",
		Description: "List merge requests in a GitLab project",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID or path",
				},
				"state": map[string]any{
					"type":        "string",
					"description": "Filter by state (opened, closed, merged, all)",
				},
				"per_page": map[string]any{
					"type":        "integer",
					"description": "Number of results per page",
				},
			},
			"required": []string{"project_id"},
		},
	}, s.handleListMRs)

	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "gitlab_get_mr_details",
		Description: "Get detailed information about a merge request",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID or path",
				},
				"mr_iid": map[string]any{
					"type":        "integer",
					"description": "The merge request IID",
				},
			},
			"required": []string{"project_id", "mr_iid"},
		},
	}, s.handleGetMRDetails)

	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "gitlab_create_MR_note",
		Description: "Create a comment on a merge request",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID or path",
				},
				"mr_iid": map[string]any{
					"type":        "integer",
					"description": "The merge request IID",
				},
				"body": map[string]any{
					"type":        "string",
					"description": "The comment text",
				},
			},
			"required": []string{"project_id", "mr_iid", "body"},
		},
	}, s.handleCreateMRNote)

	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "gitlab_get_file_content",
		Description: "Get the content of a file from a GitLab repository",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID or path",
				},
				"file_path": map[string]any{
					"type":        "string",
					"description": "Path to the file in the repository",
				},
				"ref": map[string]any{
					"type":        "string",
					"description": "The branch, tag, or commit reference",
				},
			},
			"required": []string{"project_id", "file_path", "ref"},
		},
	}, s.handleGetFileContent)

	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "gitlab_list_commits",
		Description: "List commits in a GitLab project",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID or path",
				},
				"ref_name": map[string]any{
					"type":        "string",
					"description": "The branch or tag name",
				},
				"per_page": map[string]any{
					"type":        "integer",
					"description": "Number of results per page",
				},
			},
			"required": []string{"project_id"},
		},
	}, s.handleListCommits)

	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "gitlab_list_pipelines",
		Description: "List pipelines in a GitLab project",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID or path",
				},
				"per_page": map[string]any{
					"type":        "integer",
					"description": "Number of results per page",
				},
			},
			"required": []string{"project_id"},
		},
	}, s.handleListPipelines)

	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "gitlab_create_mr",
		Description: "Create a new merge request. Cannot create MR to or from master/test branches.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID or path",
				},
				"source_branch": map[string]any{
					"type":        "string",
					"description": "The source branch name",
				},
				"target_branch": map[string]any{
					"type":        "string",
					"description": "The target branch name (cannot be master or test)",
				},
				"title": map[string]any{
					"type":        "string",
					"description": "The merge request title",
				},
				"description": map[string]any{
					"type":        "string",
					"description": "The merge request description",
				},
			},
			"required": []string{"project_id", "source_branch", "target_branch", "title"},
		},
	}, s.handleCreateMR)

	mcpsdk.AddTool(s.mcpServer, &mcpsdk.Tool{
		Name:        "gitlab_list_issues",
		Description: "List issues in a GitLab project",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID or path",
				},
				"state": map[string]any{
					"type":        "string",
					"description": "Filter by state (opened, closed, all)",
				},
				"per_page": map[string]any{
					"type":        "integer",
					"description": "Number of results per page",
				},
			},
			"required": []string{"project_id"},
		},
	}, s.handleListIssues)
}

func (s *Server) Start(ctx context.Context, transportType string) error {
	switch transportType {
	case "stdio":
		transport := &mcpsdk.StdioTransport{}
		log.Printf("Starting MCP server with stdio transport")
		return s.mcpServer.Run(ctx, transport)
	case "sse":
		handler := mcpsdk.NewSSEHandler(func(req *http.Request) *mcpsdk.Server {
			return s.mcpServer
		}, nil)

		server := &http.Server{
			Addr:    s.config.Server.Host + ":" + strconv.Itoa(s.config.Server.Port),
			Handler: handler,
		}

		log.Printf("Starting MCP server with SSE transport on %s", server.Addr)
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("HTTP server error: %v", err)
			}
		}()

		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		server.Shutdown(shutdownCtx)
		return nil
	default:
		log.Fatalf("Unsupported transport type: %s", transportType)
		return nil
	}
}

func (s *Server) GetMCPServer() *mcpsdk.Server {
	return s.mcpServer
}
