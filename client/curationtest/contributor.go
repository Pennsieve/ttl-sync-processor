package curationtest

import "github.com/pennsieve/ttl-sync-processor/client/models/curation"
import "github.com/google/uuid"

type ContributorBuilder struct {
	contributor *curation.Contributor
}

func NewContributorBuilder() *ContributorBuilder {
	affiliation := uuid.NewString()
	name := uuid.NewString()
	firstName := uuid.NewString()
	id := uuid.NewString()
	lastName := uuid.NewString()
	return &ContributorBuilder{contributor: curation.NewContributor(id, firstName, lastName, name, affiliation)}
}

func (b *ContributorBuilder) WithMiddleName() *ContributorBuilder {
	b.contributor = b.contributor.WithMiddleName(uuid.NewString())
	return b
}

func (b *ContributorBuilder) WithDataRemoteUserID() *ContributorBuilder {
	b.contributor = b.contributor.WithDataRemoteUserID(uuid.NewString())
	return b
}

func (b *ContributorBuilder) WithRoles(roleCount int) *ContributorBuilder {
	if roleCount > 0 {
		roles := make([]string, roleCount)
		for i := 0; i < roleCount; i++ {
			roles[i] = uuid.NewString()
		}
		b.contributor = b.contributor.WithRoles(roles...)
	}
	return b
}

func (b *ContributorBuilder) WithAffiliation() *ContributorBuilder {
	b.contributor = b.contributor.WithAffiliation(uuid.NewString(), uuid.NewString(), uuid.NewString())
	return b
}

func (b *ContributorBuilder) WithORCID() *ContributorBuilder {
	b.contributor = b.contributor.WithORCID(uuid.NewString(), uuid.NewString())
	return b
}

func (b *ContributorBuilder) Build() curation.Contributor {
	return *b.contributor
}
