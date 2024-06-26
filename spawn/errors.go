package spawn

import "errors"

var (
	ErrCfgEmptyOrg         = errors.New("github organization name cannot be empty")
	ErrCfgEmptyProject     = errors.New("project name cannot be empty")
	ErrCfgProjSpecialChars = errors.New("project name cannot contain special characters")
	ErrCfgBinTooShort      = errors.New("bin daemon name is too short")
	ErrCfgDenomTooShort    = errors.New("token denom is too short")
	ErrCfgHomeDirTooShort  = errors.New("home directory is too short")
	ErrCfgEmptyBech32      = errors.New("bech32 prefix cannot be empty")
	ErrCfgBech32Alpha      = errors.New("bech32 prefix must only contain alphabetical characters")
)
