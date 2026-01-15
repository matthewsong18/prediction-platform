package bets

type BetStatus int

const (
	Pending BetStatus = iota
	Won
	Lost
)

func (bs BetStatus) String() string {
	switch bs {
	case Pending:
		return "PENDING"
	case Won:
		return "WON"
	case Lost:
		return "LOST"
	default:
		return "UNKNOWN"
	}
}

type bet struct {
	PollID              string
	UserID              string
	SelectedOptionIndex int
	BetStatus           BetStatus
}

type Bet interface {
	GetBetKey() BetKey
	GetSelectedOptionIndex() int
	GetBetStatus() BetStatus
}

func (b *bet) GetBetKey() BetKey           { return BetKey{b.PollID, b.UserID} }
func (b *bet) GetSelectedOptionIndex() int { return b.SelectedOptionIndex }
func (b *bet) GetBetStatus() BetStatus     { return b.BetStatus }

type BetKey struct {
	PollID string
	UserID string
}

type UserStats struct {
	UserID       string
	Wins         int
	Losses       int
	Total        int
	WinLossRatio float64
}
