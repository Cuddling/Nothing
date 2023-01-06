package shopify

func GetGateway(shopId int) int {
	switch shopId {
	// A-Ma-Maniere
	case 6269065:
		return 26102467
	// Slam Jam
	case 57677054136:
		return 66119860408
	// Oneness Boutique
	case 1875180:
		return 3919159
	// Sneaker Politics
	case 2147974:
		return 73944301756
	default:
		return 0
	}
}
