package domain

var (
	StackoverflowType = "stackoverflow"
	GithubType        = "github"
)

type Link struct {
	ID      int64
	URL     string
	Type    string
	Tags    []string
	Filters []string
}
