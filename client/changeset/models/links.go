package models

type LinkedPropertyChanges struct {
	// The ID of the LinkedProperty in the schema. Can be empty or missing if the linkedProperty does not exist.
	// In this case, Create below should be non-empty
	ID string `json:"id,omitempty"`
	// SchemaFromID is the Pennsieve schema id of the model that acts as the parent, i.e., the "from", for these
	//linked properties. This is needed for both create and updates and deletes, so it should not be empty
	SchemaFromID string `json:"schema_from_id"`
	// If Create is non-empty, the link should be created in the model schema
	Create SchemaLinkedPropertiesCreate `json:"create,omitempty"`
}

// SchemaLinkedPropertiesCreate can be used as the payload to POST /models/datasets/<dataset id>/concepts/<model id>/linked/bulk
type SchemaLinkedPropertiesCreate []SchemaLinkedPropertyCreate

type SchemaLinkedPropertyCreate struct {
	// Name is the name of the linked property in the schema of the parent model
	Name string `json:"name"`
	// DisplayName is the display name of the linked property in the schema of the parent model
	DisplayName string `json:"displayName"`
	// SchemaTo is the model id of the child model of the linked property, i.e., the "to" model
	SchemaTo string `json:"schema_to"`
	// Position is the position of the linked property in the schema of the parent model. (?)
	Position int `json:"position"`
}

// InstanceLinkedPropertyCreatePayload can be used as the payload to
// POST /models/datasets/<dataset id>/concepts/<model id>/instances/<record id>/linked
// to create a new linked prop instance
type InstanceLinkedPropertyCreatePayload struct {
	Name                   string              `json:"name"`
	DisplayName            string              `json:"displayName"`
	SchemaLinkedPropertyID string              `json:"schemaLinkedPropertyId"`
	ToID                   PennsieveInstanceID `json:"to"`
}

// InstanceLinkedPropertyDelete contains the additional info needed to call
// DELETE /models/datasets/{dataset id}/concepts/{model Id}/instances/{record id}/linked/{link instance Id}
// No payload required
type InstanceLinkedPropertyDelete struct {
	FromID                   PennsieveInstanceID `json:"from_id"`
	InstanceLinkedPropertyID PennsieveInstanceID `json:"instance_linked_property_id"`
}
