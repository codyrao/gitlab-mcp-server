# GitLab MCP Server

A Model Context Protocol (MCP) server for GitLab that enables AI assistants to interact with GitLab APIs in a structured and secure way.

## Features

- **Multiple Transport Support**: Stdio, SSE (Server-Sent Events), and Streamable
- **Secure Operations**: Restricted operations to prevent accidental modifications to protected branches
- **GitLab Integration**: Full support for common GitLab operations
- **Docker Support**: Easy deployment with Docker

## Supported Operations

### Read Operations (Always Allowed)
- List projects
- Get project details
- List merge requests
- Get merge request details
- Get file content
- List commits
- List pipelines
- List issues

### Write Operations (With Restrictions)
- Create merge request note (comment on MR)
- Create merge request (source/target cannot be master/test)

### Restricted Operations
- Merge code (blocked)
- Push to master/test branches (blocked)
- Modify configurations (blocked)
- Delete operations (blocked)

## Installation

### From Source

```bash
git clone https://github.com/yourusername/gitlab-mcp-server.git
cd gitlab-mcp-server
go build -o gitlab-mcp-server ./cmd/server
```

### Using Docker

```bash
# Build the image
docker build -t gitlab-mcp-server .

# Run with environment variables
docker run -e GITLAB_TOKEN="your-token" \
           -e GITLAB_HOST="https://gitlab.com" \
           gitlab-mcp-server
```

### Pull from Docker Hub

```bash
docker pull codyrao/gitlab-mcp-server:latest

docker run -e GITLAB_TOKEN="your-token" \
           -e GITLAB_HOST="https://gitlab.com" \
           codyrao/gitlab-mcp-server:latest
```

## Configuration

### Configuration File (config.yaml)

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  transport: "stdio"

gitlab:
  host: "https://gitlab.com"
  token: "your-gitlab-token-here"

logging:
  level: "info"
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| GITLAB_TOKEN | GitLab personal access token | Required |
| GITLAB_HOST | GitLab instance URL | https://gitlab.com |
| GITLAB_MCP_TRANSPORT | Transport type (stdio, sse) | stdio |
| GITLAB_MCP_PORT | Port for SSE transport | 8080 |

## Usage

### STDIO Transport (Default)

For Claude Desktop or other MCP clients:

```json
{
  "mcpServers": {
    "gitlab": {
      "command": "docker run --rm -i -e GITLAB_TOKEN codyrao/gitlab-mcp-server:latest",
      "args": ["-transport", "stdio"]
    }
  }
}
```

Or with local binary:

```json
{
  "mcpServers": {
    "gitlab": {
      "command": "./gitlab-mcp-server",
      "args": ["-config", "config.yaml", "-transport", "stdio"]
    }
  }
}
```

### SSE Transport

```bash
# Using config file
./gitlab-mcp-server -config config.yaml -transport sse -port 8080

# Using environment variables
docker run -p 8080:8080 \
           -e GITLAB_TOKEN="your-token" \
           -e GITLAB_MCP_TRANSPORT="sse" \
           -e GITLAB_MCP_PORT="8080" \
           codyrao/gitlab-mcp-server:latest
```

Then configure your MCP client to connect to:
- SSE: `http://localhost:8080/connect`

## Available Tools

### gitlab_list_projects
List GitLab projects for the authenticated user.

**Parameters:**
- `membership` (boolean): Limit to projects where the authenticated user is a member
- `owned` (boolean): Limit to projects owned by the authenticated user
- `search` (string): Search projects by name
- `per_page` (integer): Number of results per page

### gitlab_get_project
Get details of a GitLab project.

**Parameters:**
- `project_id` (string): The project ID or path (e.g., 'group/project')

### gitlab_list_mrs
List merge requests in a GitLab project.

**Parameters:**
- `project_id` (string): The project ID or path
- `state` (string): Filter by state (opened, closed, merged, all)
- `per_page` (integer): Number of results per page

### gitlab_get_mr_details
Get detailed information about a merge request.

**Parameters:**
- `project_id` (string): The project ID or path
- `mr_iid` (integer): The merge request IID

### gitlab_create_MR_note
Create a comment on a merge request.

**Parameters:**
- `project_id` (string): The project ID or path
- `mr_iid` (integer): The merge request IID
- `body` (string): The comment text

### gitlab_get_file_content
Get the content of a file from a GitLab repository.

**Parameters:**
- `project_id` (string): The project ID or path
- `file_path` (string): Path to the file in the repository
- `ref` (string): The branch, tag, or commit reference

### gitlab_list_commits
List commits in a GitLab project.

**Parameters:**
- `project_id` (string): The project ID or path
- `ref_name` (string): The branch or tag name
- `per_page` (integer): Number of results per page

### gitlab_list_pipelines
List pipelines in a GitLab project.

**Parameters:**
- `project_id` (string): The project ID or path
- `per_page` (integer): Number of results per page

### gitlab_create_mr
Create a new merge request. Cannot create MR to or from master/test branches.

**Parameters:**
- `project_id` (string): The project ID or path
- `source_branch` (string): The source branch name
- `target_branch` (string): The target branch name (cannot be master or test)
- `title` (string): The merge request title
- `description` (string): The merge request description

### gitlab_list_issues
List issues in a GitLab project.

**Parameters:**
- `project_id` (string): The project ID or path
- `state` (string): Filter by state (opened, closed, all)
- `per_page` (integer): Number of results per page

## Security Restrictions

This MCP server implements strict security restrictions to prevent accidental modifications:

1. **Protected Branches**: Operations on `master`, `main`, `test`, `develop`, and `release` branches are blocked
2. **Merge Operations**: Merge requests cannot be merged through this server
3. **Delete Operations**: All delete operations are blocked
4. **Configuration Changes**: Project settings cannot be modified
5. **Branch Creation**: Creating branches with protected names is blocked

## GitLab Token Permissions

Your GitLab personal access token needs the following permissions:

- `api`: Full API access
- `read_api`: Read API access
- `read_repository`: Read repository access
- `write_repository`: Write repository access (for create MR operations)

## Development

### Running Locally

```bash
# Install dependencies
go mod download

# Run the server
go run ./cmd/server -config config.yaml -transport stdio
```

### Building for Docker

```bash
# Build the image
docker build -t gitlab-mcp-server .

# Push to Docker Hub
docker login -u codyrao
docker tag gitlab-mcp-server codyrao/gitlab-mcp-server:latest
docker push codyrao/gitlab-mcp-server:latest
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
