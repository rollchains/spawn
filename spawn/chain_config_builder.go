package spawn

func (s NewChainConfig) WithOrg(org string) NewChainConfig {
	s.GithubOrg = org
	return s
}

func (s NewChainConfig) WithProjectName(proj string) NewChainConfig {
	s.ProjectName = proj
	return s
}

func (s NewChainConfig) WithBech32Prefix(bech string) NewChainConfig {
	s.Bech32Prefix = bech
	return s
}

func (s NewChainConfig) WithDenom(denom string) NewChainConfig {
	s.Denom = denom
	return s
}

func (s NewChainConfig) WithHomeDir(home string) NewChainConfig {
	s.HomeDir = home
	return s
}

func (s NewChainConfig) WithBinDaemon(bin string) NewChainConfig {
	s.BinDaemon = bin
	return s
}
