package constants

const (
	HardhatDefaultWallet1  string = "f39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	HardhatDefaultWallet2  string = "70997970C51812dc3A010C7d01b50e0d17dc79C8"
	HardhatDefaultWallet3  string = "3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"
	HardhatDefaultWallet4  string = "90F79bf6EB2c4f870365E785982E1f101E93b906"
	HardhatDefaultWallet5  string = "15d34AAf54267DB7D7c367839AAf71A00a2C6A65"
	HardhatDefaultWallet6  string = "9965507D1a55bcC2695C58ba16FB37d819B0A4dc"
	HardhatDefaultWallet7  string = "976EA74026E726554dB657fA54763abd0C3a0aa9"
	HardhatDefaultWallet8  string = "14dC79964da2C08b23698B3D3cc7Ca32193d9955"
	HardhatDefaultWallet9  string = "23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f"
	HardhatDefaultWallet10 string = "a0Ee7A142d267C1f36714E4a8F75612F20a79720"
	HardhatDefaultWallet11 string = "Bcd4042DE499D14e55001CcbB24a551F3b954096"
	HardhatDefaultWallet12 string = "71bE63f3384f5fb98995898A86B02Fb2426c5788"
	HardhatDefaultWallet13 string = "FABB0ac9d68B0B445fB7357272Ff202C5651694a"
	HardhatDefaultWallet14 string = "1CBd3b2770909D4e10f157cABC84C7264073C9Ec"
	HardhatDefaultWallet15 string = "dF3e18d64BC6A983f673Ab319CCaE4f1a57C7097"
	HardhatDefaultWallet16 string = "cd3B766CCDd6AE721141F452C550Ca635964ce71"
	HardhatDefaultWallet17 string = "2546BcD3c84621e976D8185a91A922aE77ECEc30"
	HardhatDefaultWallet18 string = "bDA5747bFD65F08deb54cb465eB87D40e51B197E"
	HardhatDefaultWallet19 string = "dD2FD4581271e230360230F9337D5c0430Bf44C0"
	HardhatDefaultWallet20 string = "8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199"
)

// GetKovanWallet retrieves the wallet details for our Kovan test account
func GetKovanWallet() string {
	return getSecret("kovanWallet")
}

// Can utilize IAM roles / some other way of safely authenticating with a secrets management tool to retrieve secret
// keys and wallet addresses
func getSecret(name string) string {
	// Grab secret from somewhere
	return ""
}
