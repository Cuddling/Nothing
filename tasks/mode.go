package tasks

type TaskMode int

const (
	ModeShopifySafe = iota
	ModeShopifyFast
)

// ModeToString Returns a stringified version of a tasks's mode
func ModeToString(mode TaskMode) string {
	switch mode {
	case ModeShopifySafe:
		return "Safe"
	case ModeShopifyFast:
		return "Fast"
	default:
		return "None"
	}
}
