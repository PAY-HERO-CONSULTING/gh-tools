package custom_types

type UUIDIdentifier struct {
	ID string `json:"id"`
}

func (i UUIDIdentifier) IsNew() bool {
	return i.ID == ""
}

func (i *UUIDIdentifier) UUID() {

	if i.ID == "" {
		i.ID = GenerateUUID()
	}
}
