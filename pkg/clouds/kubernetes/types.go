package kubernetes

import "strings"

type APIVerb string

var APIVerbList []APIVerb = []APIVerb{
	API_VERB_CREATE,
	API_VERB_GET,
	API_VERB_LIST,
	API_VERB_WATCH,
	API_VERB_UPDATE,
	API_VERB_PATCH,
	API_VERB_DELETE,
	API_VERB_DELETECOLLECTION,
	API_VERB_IMPERSONATE,
	API_VERB_BIND,
	API_VERB_APPROVE,
	API_VERB_ESCALATE,
}

const (
	API_VERB_CREATE           APIVerb = "create"
	API_VERB_GET              APIVerb = "get"
	API_VERB_LIST             APIVerb = "list"
	API_VERB_WATCH            APIVerb = "watch"
	API_VERB_UPDATE           APIVerb = "update"
	API_VERB_PATCH            APIVerb = "patch"
	API_VERB_DELETE           APIVerb = "delete"
	API_VERB_DELETECOLLECTION APIVerb = "deletecollection"
	API_VERB_IMPERSONATE      APIVerb = "impersonate"
	API_VERB_BIND             APIVerb = "bind"
	API_VERB_APPROVE          APIVerb = "approve"
	API_VERB_ESCALATE         APIVerb = "escalate"
)

type Resource struct {
	GroupName    string
	GroupVersion string
	Name         string
	Namespaced   bool
	SubResource  string
}

func (r *Resource) String() string {
	rb := &strings.Builder{}

	rb.WriteString(r.Name)

	if r.GroupVersion != "" {
		rb.WriteRune('/')
		rb.WriteString(r.GroupVersion)
	}

	if r.SubResource != "" {
		rb.WriteRune('/')
		rb.WriteString(r.SubResource)
	}

	return rb.String()
}
