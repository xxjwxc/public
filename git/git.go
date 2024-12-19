package git

import (
	"crypto/tls"
	"net/http"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// https://blog.csdn.net/zhang_yasong/article/details/138196565
type GitTools struct {
	gitLabClient *gitlab.Client
}

func NewGitLabController(url, token string) (*GitTools, error) {
	// 创建一个自定义的http.Client  忽略证书
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(url), gitlab.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return &GitTools{gitLabClient: git}, nil
}

// GetUser 获取用户信息
func (g *GitTools) GetUser() (*gitlab.User, error) {
	u, _, err := g.gitLabClient.Users.CurrentUser()
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetALLProjects 获取所有项目
func (g *GitTools) GetALLProjects() ([]*gitlab.Project, error) {
	lbo := &gitlab.ListProjectsOptions{ListOptions: gitlab.ListOptions{Page: 1, PerPage: 50}}
	var pro []*gitlab.Project
	for {
		p, _, err := g.gitLabClient.Projects.ListProjects(lbo)
		if err != nil {
			return nil, err
		}
		pro = append(pro, p...)
		if len(p) < 50 {
			break
		}
		lbo.ListOptions.Page++
	}
	return pro, nil
}

// GetProjectFromName 通过项目名获取项目
func (g *GitTools) GetProjectFromName(projectName string) (*gitlab.Project, error) {
	p, _, err := g.gitLabClient.Projects.GetProject(projectName, nil)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetProjectFromID 通过项目ID获取项目
func (g *GitTools) GetProjectFromID(pid int) (*gitlab.Project, error) {
	p, _, err := g.gitLabClient.Projects.GetProject(pid, nil)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetBranch 获取分支
func (g *GitTools) GetBranch(pid int, branchName string) (*gitlab.Branch, error) {
	p, _, err := g.gitLabClient.Branches.GetBranch(pid, branchName)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetTag 获取tag
func (g *GitTools) GetTag(pid int, tagName string) (*gitlab.Tag, error) {
	p, _, err := g.gitLabClient.Tags.GetTag(pid, tagName)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetCommits 获取commits
func (g *GitTools) GetCommits(pid int) ([]*gitlab.Commit, error) {
	lco := &gitlab.ListCommitsOptions{ListOptions: gitlab.ListOptions{Page: 1, PerPage: 50}}
	var pro []*gitlab.Commit
	for {
		p, _, err := g.gitLabClient.Commits.ListCommits(pid, lco)
		if err != nil {
			return nil, err
		}
		pro = append(pro, p...)
		if len(p) < 50 {
			break
		}
		lco.ListOptions.Page++
	}
	return pro, nil
}

// GetTopCommit 获取commit
func (g *GitTools) GetTopCommit(pid int, branchName string) (*gitlab.Commit, error) {
	p, _, err := g.gitLabClient.Commits.GetCommit(pid, branchName, nil)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetTopCommit 获取commit
func (g *GitTools) GetAllCommit(pid int, branchName string) ([]*gitlab.Commit, error) {
	lco := &gitlab.ListCommitsOptions{RefName: gitlab.Ptr(branchName), ListOptions: gitlab.ListOptions{Page: 1, PerPage: 50}}
	var pro []*gitlab.Commit
	for {
		p, _, err := g.gitLabClient.Commits.ListCommits(pid, lco)
		if err != nil {
			return nil, err
		}
		pro = append(pro, p...)
		if len(p) < 50 {
			break
		}
		lco.ListOptions.Page++
	}
	return pro, nil
}

// GetALLBranches 获取所有分支
func (g *GitTools) GetALLBranches(pid int) ([]*gitlab.Branch, error) {
	lbo := &gitlab.ListBranchesOptions{ListOptions: gitlab.ListOptions{Page: 1, PerPage: 50}}
	var pro []*gitlab.Branch
	for {
		p, _, err := g.gitLabClient.Branches.ListBranches(pid, lbo)
		if err != nil {
			return nil, err
		}
		pro = append(pro, p...)
		if len(p) < 50 {
			break
		}
		lbo.ListOptions.Page++
	}
	return pro, nil
}

// GetALLTags 获取所有tag
func (g *GitTools) GetALLTags(pid int) ([]*gitlab.Tag, error) {
	lto := &gitlab.ListTagsOptions{ListOptions: gitlab.ListOptions{Page: 1, PerPage: 50}}
	var pro []*gitlab.Tag
	for {
		p, _, err := g.gitLabClient.Tags.ListTags(pid, lto)
		if err != nil {
			return nil, err
		}
		pro = append(pro, p...)
		if len(p) < 50 {
			break
		}
		lto.ListOptions.Page++
	}
	return pro, nil
}
