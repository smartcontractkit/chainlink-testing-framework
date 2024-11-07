package simple_node_set

// UpdateNodeConfigs updates nodes configuration TOML files
// this API is discouraged, however, you can use it if nodes require restart or configuration updates, temporarily!
func UpdateNodeConfigs(in *Input, cfg string) {
	in.NodeSpecs[0].Node.UserConfigOverrides = in.NodeSpecs[0].Node.UserConfigOverrides + cfg
	in.Out = nil
}
