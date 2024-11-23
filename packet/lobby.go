package packet

const LobbyInfoSize = 1191

// Lobby Info Packet
type LobbyInfo struct {
	AiControlled uint8  `json:"m_aiControlled"` // Whether the vehicle is AI (1) or Human (0) controlled
	TeamId       uint8  `json:"m_teamId"`       // Team id - see appendix (255 if no team currently selected)
	Nationality  uint8  `json:"m_national"`     // Nationality of the driver
	Name         string `json:"m_name"`         // Name of participant in UTF-8 format – null terminated  Will be truncated with ... (U+2026) if too long
	CarNumber    uint8  `json:"m_carNumber"`    // Car number of the player
	ReadyStatus  uint8  `json:"m_readyStatus"`  // 0 = not ready, 1 = ready, 2 = spectating
}

var _ Data = (*LobbyInfoData)(nil)

type LobbyInfoData struct {
	Header     Header
	LobbyInfos [22]LobbyInfo
}

func (l LobbyInfoData) GetHeader() Header {
	return l.Header
}

func (l LobbyInfoData) Id() Id {
	return LobbyInfoId
}
