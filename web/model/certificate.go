package model

// Certificate database model that logs stores information about active certificates
type Certificate struct {
	Domain        string `json:"domain" binding:"required"`
	ChallengeType string `json:"challengeType" binding:"required"`
}
