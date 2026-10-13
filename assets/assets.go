package assets

import (
	"embed"
)

//go:embed pieces/*
var PiecesFS embed.FS

//go:embed fonts/*
var FontsFS embed.FS
