package model

type Request struct {
	Domain        string `json:"domain" binding:"required"`
	ChallengeType string `json:"challengeType" binding:"required"`
}

func (h Request) GetByID(id string) (*Request, error) {
	return &Request{
		Domain:        "example.com",
		ChallengeType: "test",
	}, nil
}
