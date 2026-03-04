// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

// Size represents width and height in pixels.
type Size struct {
	Width  int `json:"width" yaml:"width"`   // Width in pixels.
	Height int `json:"height" yaml:"height"` // Height in pixels.
}

// Point represents a 2D pixel position.
type Point struct {
	X int `json:"x" yaml:"x"` // Horizontal coordinate.
	Y int `json:"y" yaml:"y"` // Vertical coordinate.
}

// Texture defines a texture reference used by .imageset.
type Texture struct {
	Path string `json:"path" yaml:"path"`                     // Texture path.
	Mpix int    `json:"mpix,omitempty" yaml:"mpix,omitempty"` // Pixels per meter.
}

// Image defines one sprite entry in .imageset.
type Image struct {
	Name  string `json:"name" yaml:"name"`                       // Image name.
	Pos   Point  `json:"pos" yaml:"pos"`                         // Top-left position.
	Size  Size   `json:"size" yaml:"size"`                       // Sprite size.
	Flags Flags  `json:"flags,omitempty" yaml:"flags,omitempty"` // Tile flags bitset.
}

// Group defines a named collection of images.
type Group struct {
	Name   string  `json:"name" yaml:"name"`                         // Group name.
	Images []Image `json:"images,omitempty" yaml:"images,omitempty"` // Group images.
}

// Document is the root .imageset model.
type Document struct {
	Name     string    `json:"name,omitempty" yaml:"name,omitempty"`         // Set name.
	Textures []Texture `json:"textures,omitempty" yaml:"textures,omitempty"` // Textures list.
	Images   []Image   `json:"images,omitempty" yaml:"images,omitempty"`     // Root images.
	Groups   []Group   `json:"groups,omitempty" yaml:"groups,omitempty"`     // Groups list.
	RefSize  Size      `json:"ref_size" yaml:"ref_size"`                     // Atlas size.
}
