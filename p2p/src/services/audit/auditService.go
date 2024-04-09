package audit

import "time"

func CreateAuditForUser(iPAddress string) *Audit {
	currentTime := time.Now()
	return CreateAudit(
		WithIPAddress(iPAddress),
		WithRegisterDate(currentTime.Format(time.RFC1123)),
		WithLastOnlineDate(currentTime.Format(time.RFC1123)),
		WithLocationData("LAN"),
		WithNumberOfActions(1))
}

func CreateAuditForRoom() *Audit {
	currentTime := time.Now()
	return CreateAudit(
		WithCreateDate(currentTime.Format(time.RFC1123)),
		WithLocationData("LAN"),
		WithNumberOfActions(1))
}

func CreateAuditForMessage() *Audit {
	currentTime := time.Now()
	return CreateAudit(
		WithCreateDate(currentTime.Format(time.RFC1123)),
		WithLocationData("LAN"))
}
