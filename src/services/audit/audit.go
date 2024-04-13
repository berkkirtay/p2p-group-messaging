// Copyright (c) 2024 Berk Kirtay

package audit

type Audit struct {
	IPAddress       string `json:"ipAddress,omitempty" bson:"ipAddress,omitempty"`
	RegisterDate    string `json:"registerDate,omitempty" bson:"registerDate,omitempty"`
	CreateDate      string `json:"createDate,omitempty" bson:"createDate,omitempty"`
	LastOnlineDate  string `json:"lastOnlineDate,omitempty" bson:"lastOnlineDate,omitempty"`
	LocationData    string `json:"locationData,omitempty" bson:"locationData,omitempty"`
	NumberOfActions int64  `json:"numberOfActions,omitempty" bson:"numberOfActions,omitempty"`
}

type AuditOption func(Audit) Audit

func WithIPAddress(iPAddress string) AuditOption {
	return func(audit Audit) Audit {
		audit.IPAddress = iPAddress
		return audit
	}
}

func WithRegisterDate(registerDate string) AuditOption {
	return func(audit Audit) Audit {
		audit.RegisterDate = registerDate
		return audit
	}
}

func WithCreateDate(createDate string) AuditOption {
	return func(audit Audit) Audit {
		audit.CreateDate = createDate
		return audit
	}
}

func WithLastOnlineDate(lastOnlineDate string) AuditOption {
	return func(audit Audit) Audit {
		audit.LastOnlineDate = lastOnlineDate
		return audit
	}
}

func WithLocationData(locationData string) AuditOption {
	return func(audit Audit) Audit {
		audit.LocationData = locationData
		return audit
	}
}

func WithNumberOfActions(numberOfActions int64) AuditOption {
	return func(audit Audit) Audit {
		audit.NumberOfActions = numberOfActions
		return audit
	}
}

func CreateDefaultAudit() Audit {
	return Audit{}
}

func CreateAudit(options ...AuditOption) *Audit {
	audit := CreateDefaultAudit()

	for _, option := range options {
		audit = option(audit)
	}

	return &audit
}
