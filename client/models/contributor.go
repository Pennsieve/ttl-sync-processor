package models

type Contributor struct {
	Affiliation            *Affiliation `json:"affiliation,omitempty""`
	ContributorAffiliation string       `json:"contributor_affiliation"`
	ContributorName        string       `json:"contributor_name"`
	ContributorORCID       *ORCID       `json:"contributor_orcid,omitempty"`
	ContributorRole        []string     `json:"contributor_role,omitempty"`
	FirstName              string       `json:"first_name"`
	Id                     string       `json:"id"`
	LastName               string       `json:"last_name"`
}

func NewContributor(contributorAffiliation string, contributorName string, firstName string, id string, lastName string) *Contributor {
	return &Contributor{
		ContributorAffiliation: contributorAffiliation,
		ContributorName:        contributorName,
		FirstName:              firstName,
		Id:                     id,
		LastName:               lastName,
	}
}

func (c *Contributor) WithRoles(contributorRoles ...string) *Contributor {
	c.ContributorRole = append(c.ContributorRole, contributorRoles...)
	return c
}

func (c *Contributor) WithAffiliation(id string, label string, system string) *Contributor {
	affiliation := NewAffiliation(id, label, system)
	c.Affiliation = affiliation
	return c
}

func (c *Contributor) WithORCID(id string, label string) *Contributor {
	orcid := NewORCID(id, label)
	c.ContributorORCID = orcid
	return c
}

type Affiliation struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	System string `json:"system"`
	Type   string `json:"type"`
}

func NewAffiliation(id string, label string, system string) *Affiliation {
	return &Affiliation{Id: id, Label: label, System: system, Type: "identifier"}
}

type ORCID struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	System string `json:"system"`
	Type   string `json:"type"`
}

func NewORCID(id string, label string) *ORCID {
	return &ORCID{Id: id, Label: label, System: "Orcid", Type: "identifier"}
}
