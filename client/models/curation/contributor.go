package curation

type Contributor struct {
	Affiliation            *Affiliation `json:"affiliation,omitempty"`
	ContributorAffiliation string       `json:"contributor_affiliation"`
	ContributorName        string       `json:"contributor_name"`
	ContributorORCID       *ORCID       `json:"contributor_orcid,omitempty"`
	ContributorRole        []string     `json:"contributor_role,omitempty"`
	FirstName              string       `json:"first_name"`
	ID                     string       `json:"id"`
	LastName               string       `json:"last_name"`
	MiddleName             string       `json:"middle_name,omitempty"`
}

func NewContributor(id string, firstName string, lastName string, contributorName string, contributorAffiliation string) *Contributor {
	return &Contributor{
		ContributorAffiliation: contributorAffiliation,
		ContributorName:        contributorName,
		FirstName:              firstName,
		ID:                     id,
		LastName:               lastName,
	}
}

func (c *Contributor) WithMiddleName(middleName string) *Contributor {
	c.MiddleName = middleName
	return c
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

type Affiliation embeddedIdentifier

func NewAffiliation(id string, label string, system string) *Affiliation {
	affIdentifier := newDescriptionlessIdentifier(id, label, system)
	return (*Affiliation)(&affIdentifier)
}

type ORCID embeddedIdentifier

func NewORCID(id string, label string) *ORCID {
	orcidIdentifier := newDescriptionlessIdentifier(id, label, "Orcid")
	return (*ORCID)(&orcidIdentifier)
}
