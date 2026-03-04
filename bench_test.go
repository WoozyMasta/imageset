// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"io"
	"path/filepath"
	"testing"
)

var (
	// benchmarkSinkInt prevents benchmark code elimination.
	benchmarkSinkInt int
)

// BenchmarkParseBytes measures parse throughput from in-memory bytes.
func BenchmarkParseBytes(b *testing.B) {
	for _, fixture := range fixtureCases {
		fixture := fixture
		data := readFixtureBytes(b, fixture.file)

		b.Run(fixture.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				doc, err := ParseBytes(data)
				if err != nil {
					b.Fatalf("ParseBytes(%s): %v", fixture.file, err)
				}

				benchmarkSinkInt = len(doc.Images) + len(doc.Groups)
			}
		})
	}
}

// BenchmarkParseFile measures parse throughput from file path.
func BenchmarkParseFile(b *testing.B) {
	for _, fixture := range fixtureCases {
		fixture := fixture
		path := filepath.Join("testdata", fixture.file)

		b.Run(fixture.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				doc, err := ParseFile(path)
				if err != nil {
					b.Fatalf("ParseFile(%s): %v", fixture.file, err)
				}

				benchmarkSinkInt = len(doc.Images) + len(doc.Groups)
			}
		})
	}
}

// BenchmarkValidate measures validation on already parsed documents.
func BenchmarkValidate(b *testing.B) {
	for _, fixture := range fixtureCases {
		fixture := fixture
		doc, err := ParseFile(filepath.Join("testdata", fixture.file))
		if err != nil {
			b.Fatalf("ParseFile(%s): %v", fixture.file, err)
		}

		b.Run(fixture.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				if err := Validate(doc); err != nil {
					b.Fatalf("Validate(%s): %v", fixture.file, err)
				}
			}
		})
	}
}

// BenchmarkFormat measures canonical text formatting throughput.
func BenchmarkFormat(b *testing.B) {
	opts := &FormatOptions{UseCamelCaseNames: false, Indent: "  "}
	for _, fixture := range fixtureCases {
		fixture := fixture
		doc, err := ParseFile(filepath.Join("testdata", fixture.file))
		if err != nil {
			b.Fatalf("ParseFile(%s): %v", fixture.file, err)
		}

		b.Run(fixture.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				out, err := Format(doc, opts)
				if err != nil {
					b.Fatalf("Format(%s): %v", fixture.file, err)
				}

				benchmarkSinkInt = len(out)
			}
		})
	}
}

// BenchmarkWrite measures write throughput to io.Discard.
func BenchmarkWrite(b *testing.B) {
	opts := &FormatOptions{UseCamelCaseNames: false, Indent: "  "}
	for _, fixture := range fixtureCases {
		fixture := fixture
		doc, err := ParseFile(filepath.Join("testdata", fixture.file))
		if err != nil {
			b.Fatalf("ParseFile(%s): %v", fixture.file, err)
		}

		b.Run(fixture.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				if err := Write(io.Discard, doc, opts); err != nil {
					b.Fatalf("Write(%s): %v", fixture.file, err)
				}
			}
		})
	}
}

// BenchmarkPipeline measures parse+validate+format in one loop.
func BenchmarkPipeline(b *testing.B) {
	opts := &FormatOptions{UseCamelCaseNames: false, Indent: "  "}
	for _, fixture := range fixtureCases {
		fixture := fixture
		data := readFixtureBytes(b, fixture.file)

		b.Run(fixture.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				doc, err := ParseBytes(data)
				if err != nil {
					b.Fatalf("ParseBytes(%s): %v", fixture.file, err)
				}
				if err := Validate(doc); err != nil {
					b.Fatalf("Validate(%s): %v", fixture.file, err)
				}

				out, err := Format(doc, opts)
				if err != nil {
					b.Fatalf("Format(%s): %v", fixture.file, err)
				}

				benchmarkSinkInt = len(out)
			}
		})
	}
}
