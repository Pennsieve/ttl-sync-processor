package sync

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/metadatatest"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestComputeChangeset(t *testing.T) {
	for scenario, test := range map[string]func(tt *testing.T){
		"handle edge cases without panic": emptyChangesets,
		"smoke test":                      smokeTest,
	} {
		t.Run(scenario, func(t *testing.T) {
			test(t)
		})
	}
}

func emptyChangesets(t *testing.T) {
	changes, err := ComputeChangeset(emptySchema, &metadata.SavedDatasetMetadata{}, &metadata.DatasetMetadata{})
	require.NoError(t, err)
	require.NotNil(t, changes)
	assert.Empty(t, changes.Models)
	assert.Empty(t, changes.LinkedProperties)
}

func smokeTest(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().
		WithModel(metadata.ContributorModelName, metadata.ContributorDisplayName).
		WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).
		WithModel(metadata.SampleModelName, metadata.SampleDisplayName).
		WithLinkedProperty(metadata.SampleSubjectLinkName, metadata.SampleSubjectLinkDisplayName).Build(),
	)

	contributor := metadatatest.NewContributorBuilder().WithNodeID().Build()
	subject := metadatatest.NewSubjectBuilder().Build()
	sample := metadatatest.NewSampleBuilder().Build()
	sampleSubject := metadata.SampleSubject{
		SampleID:  sample.ExternalID(),
		SubjectID: subject.ExternalID(),
	}

	changes, err := ComputeChangeset(schemaData,
		&metadata.SavedDatasetMetadata{},
		&metadata.DatasetMetadata{
			Contributors:   []metadata.Contributor{contributor},
			Subjects:       []metadata.Subject{subject},
			Samples:        []metadata.Sample{sample},
			SampleSubjects: []metadata.SampleSubject{sampleSubject},
		},
	)
	require.NoError(t, err)

	assert.Len(t, changes.Models, 3)
	for _, m := range changes.Models {
		require.NotNil(t, m.ID)
		assert.Len(t, m.Records.Create, 1)
	}

	assert.Len(t, changes.LinkedProperties, 1)
	sampleSubjectChanges := changes.LinkedProperties[0]
	assert.NotNil(t, sampleSubjectChanges.ID)
	assert.Len(t, sampleSubjectChanges.Instances.Create, 1)
}
