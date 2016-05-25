package controller

// PeerOffer struct: use when offering to start P2P
type PeerOffer struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
	From string `json:"from"`
}

// PeerCandidate struct: use after offering
type PeerCandidate struct {
	Candidate PeerCandidateChild `json:"candidate"`
	From      string             `json:"from"`
}

// PeerCandidateChild struct
type PeerCandidateChild struct {
	Candidate     string `json:"candidate"`
	SdpMLineIndex int    `json:"sdpMLineIndex"`
	SdpMid        string `json:"sdpMid"`
}

// PeerPool struct: connected client ids
type PeerPool struct {
	Type string   `json:"type"`
	Ids  []string `json:"ids"`
}
