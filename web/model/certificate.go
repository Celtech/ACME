package model

type Certificate struct {
	Domain        string `json:"domain" binding:"required"`
	ChallengeType string `json:"challengeType" binding:"required"`
}
