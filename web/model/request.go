package model

type Request struct {
	Domain        string `json:"domain" binding:"required"`
	ChallengeType string `json:"challengeType" binding:"required"`
}
