package metadata

import changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"

const ContributorModelName = "contributor"
const ContributorDisplayName = "Contributor"

// Keys should correspond to json struct tags

const FirstNameKey = "first_name"
const LastNameKey = "last_name"
const MiddleInitialKey = "middle_initial"
const DegreeKey = "degree"
const ORCIDKey = "orcid"
const NodeIDKey = "node_id"

type Contributor struct {
	FirstName     string `json:"first_name"`
	MiddleInitial string `json:"middle_initial,omitempty"`
	LastName      string `json:"last_name"`
	Degree        string `json:"degree,omitempty"`
	ORCID         string `json:"orcid,omitempty"`
	NodeID        string `json:"node_id,omitempty"`
}

type SavedContributor struct {
	PennsieveID changesetmodels.PennsieveInstanceID `json:"-"`
	Contributor
}

func (s SavedContributor) GetPennsieveID() changesetmodels.PennsieveInstanceID {
	return s.PennsieveID
}
