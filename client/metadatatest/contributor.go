package metadatatest

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"math/rand"
)

type ContributorBuilder struct {
	c *metadata.Contributor
}

func NewContributorBuilder() *ContributorBuilder {
	c := &metadata.Contributor{
		FirstName: uuid.NewString(),
		LastName:  uuid.NewString(),
	}
	return &ContributorBuilder{c: c}
}
func (b *ContributorBuilder) Build() metadata.Contributor {
	return *b.c
}

func (b *ContributorBuilder) WithMiddleInitial() *ContributorBuilder {
	someLetters := "ABCDEFGHIJKLMNOPQURSTUVXYZ"
	b.c.MiddleInitial = fmt.Sprintf("%c", someLetters[rand.Intn(len(someLetters))])
	return b
}

func (b *ContributorBuilder) WithORICID() *ContributorBuilder {
	b.c.ORCID = fmt.Sprintf("https://orcid.org/%s", uuid.NewString())
	return b
}

func (b *ContributorBuilder) WithNodeID() *ContributorBuilder {
	b.c.NodeID = fmt.Sprintf("N:user:%s", uuid.NewString())
	return b
}

func (b *ContributorBuilder) WithDegree() *ContributorBuilder {
	b.c.Degree = uuid.NewString()
	return b
}
