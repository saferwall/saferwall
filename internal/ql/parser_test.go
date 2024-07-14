package ql

import (
	"testing"
)

func TestParser(t *testing.T) {
	/*
			t.Run("TestParseFileSizeQuery", func(t *testing.T) {
				tests := []struct {
					input         string
					expectedQuery SingleQuery
				}{
					{
						input: `size:100kb+`,
						expectedQuery: SingleQuery{
							token.Token{
								File,
								"size",
							},
							GREATERTHAN,
							"100000",
						},
					}, {
						input: `size:1mb-`,
						expectedQuery: SingleQuery{
							Token{
								SIZE,
								"size",
							},
							LESSERTHAN,
							"1000000",
						},
					}, {
						input: `size:250kb`,
						expectedQuery: SingleQuery{
							Token{
								SIZE,
								"size",
							},
							EQUAL,
							"250000",
						},
					}, {
						input: `size:500+`,
						expectedQuery: SingleQuery{
							Token{
								SIZE,
								"size",
							},
							GREATERTHAN,
							"500",
						},
					}, {
						input: `size:200`,
						expectedQuery: SingleQuery{
							Token{
								SIZE,
								"size",
							},
							EQUAL,
							"200",
						},
					}, {
						input: `size:200kb-`,
						expectedQuery: SingleQuery{
							Token{
								SIZE,
								"size",
							},
							LESSERTHAN,
							"200000",
						},
					},
				}

				for i, tt := range tests {
					l := NewLexer(tt.input)
					p := NewParser(l)
					query, err := p.parseFileSizeQuery()
					if err != nil {
						t.Fatalf("test[%d] failed with error :%q", i, err)
					}
					if query.Token != tt.expectedQuery.Token {
						t.Fatalf("test[%d] expected token=%q got token=%q", i, tt.expectedQuery.Token, query.Token)
					} else if query.Value != tt.expectedQuery.Value {
						t.Fatalf("test[%d] expected value=%q got value=%q", i, tt.expectedQuery.Value, query.Value)
					} else if query.Op != tt.expectedQuery.Op {
						t.Fatalf("test[%d]=%s expected op=%q got op=%q", i, tt.input, tt.expectedQuery.Op, query.Op)
					}

				}
			})
			t.Run("TestParseFileTypeQuery", func(t *testing.T) {
				tests := []struct {
					input         string
					expectedQuery SingleQuery
				}{
					{
						input: `type:pe`,
						expectedQuery: SingleQuery{
							Token{
								TYPE,
								"type",
							},
							EQUAL,
							"pe",
						},
					}, {
						input: `type:pedll`,
						expectedQuery: SingleQuery{
							Token{
								TYPE,
								"type",
							},
							EQUAL,
							"pedll",
						},
					}, {
						input: `type:elf`,
						expectedQuery: SingleQuery{
							Token{
								TYPE,
								"type",
							},
							EQUAL,
							"elf",
						},
					},
				}

				for i, tt := range tests {
					l := NewLexer(tt.input)
					p := NewParser(l)
					query, err := p.parseFileTypeQuery()
					if err != nil {
						t.Fatalf("test[%d] failed with error :%q", i, err)
					}
					if query.Token != tt.expectedQuery.Token {
						t.Fatalf("test[%d] expected token=%q got token=%q", i, tt.expectedQuery.Token, query.Token)
					} else if query.Value != tt.expectedQuery.Value {
						t.Fatalf("test[%d] expected value=%q got value=%q", i, tt.expectedQuery.Value, query.Value)
					} else if query.Op != tt.expectedQuery.Op {
						t.Fatalf("test[%d] expected op=%q got op=%q", i, tt.expectedQuery.Op, query.Op)
					}
				}
			})
			t.Run("TestParseDatetimeQuery", func(t *testing.T) {
				tests := []struct {
					input         string
					expectedQuery SingleQuery
				}{
					{
						input: `ls:"2009-01-01T19:59:22"`,
						expectedQuery: SingleQuery{
							Token{
								LS,
								"ls",
							},
							EQUAL,
							"2009-01-01T19:59:22",
						},
					}, {
						input: `fs:"2012-08-2116:00:00"`,
						expectedQuery: SingleQuery{
							Token{
								FS,
								"fs",
							},
							EQUAL,
							"2012-08-2116:00:00",
						},
					}, {
						input: `fs:"2012-08-2116:59:22"`,
						expectedQuery: SingleQuery{
							Token{
								FS,
								"fs",
							},
							EQUAL,
							"2012-08-2116:59:22",
						},
					},
				}

				for i, tt := range tests {
					l := NewLexer(tt.input)
					p := NewParser(l)
					query, err := p.parseDatetimeQuery()
					if err != nil {
						t.Fatalf("test[%d] failed with error :%q", i, err)
					}
					if query.Token != tt.expectedQuery.Token {
						t.Fatalf("test[%d] expected token=%q got token=%q", i, tt.expectedQuery.Token, query.Token)
					} else if query.Value != tt.expectedQuery.Value {
						t.Fatalf("test[%d] expected value=%q got value=%q", i, tt.expectedQuery.Value, query.Value)
					} else if query.Op != tt.expectedQuery.Op {
						t.Fatalf("test[%d] expected op=%q got op=%q", i, tt.expectedQuery.Op, query.Op)
					}
				}
			})
			t.Run("TestParsePositivesQuery", func(t *testing.T) {
				tests := []struct {
					input         string
					expectedQuery SingleQuery
				}{
					{
						input: `positives:5`,
						expectedQuery: SingleQuery{
							Token{
								POSITIVES,
								"positives",
							},
							EQUAL,
							"5",
						},
					}, {
						input: `positives:10+`,
						expectedQuery: SingleQuery{
							Token{
								POSITIVES,
								"positives",
							},
							GREATERTHAN,
							"10",
						},
					},
				}

				for i, tt := range tests {
					l := NewLexer(tt.input)
					p := NewParser(l)
					query, err := p.parsePositivesQuery()
					if err != nil {
						t.Fatalf("test[%d] failed with error :%q", i, err)
					}
					if query.Token != tt.expectedQuery.Token {
						t.Fatalf("test[%d] expected token=%q got token=%q", i, tt.expectedQuery.Token, query.Token)
					} else if query.Value != tt.expectedQuery.Value {
						t.Fatalf("test[%d] expected value=%q got value=%q", i, tt.expectedQuery.Value, query.Value)
					} else if query.Op != tt.expectedQuery.Op {
						t.Fatalf("test[%d] expected op=%q got op=%q", i, tt.expectedQuery.Op, query.Op)
					}
				}
			})
			t.Run("TestParseSimpleQuery", func(t *testing.T) {
				tests := []struct {
					input         string
					expectedQuery SingleQuery
				}{
					{
						input: `type:peexe`,
						expectedQuery: SingleQuery{
							Token{
								TYPE,
								"type",
							},
							EQUAL,
							"peexe",
						},
					}, {
						input: `size:200kb+`,
						expectedQuery: SingleQuery{
							Token{
								SIZE,
								"size",
							},
							GREATERTHAN,
							"200000",
						},
					}, {
						input: `size:1mb-`,
						expectedQuery: SingleQuery{
							Token{
								SIZE,
								"size",
							},
							LESSERTHAN,
							"1000000",
						},
					}, {
						input: `size:2mb-,`,
						expectedQuery: SingleQuery{
							Token{
								SIZE,
								"size",
							},
							LESSERTHAN,
							"2000000",
						},
					}, {
						input: `ls:"2009-01-01T19:59:22"`,
						expectedQuery: SingleQuery{
							Token{
								LS,
								"ls",
							},
							EQUAL,
							"2009-01-01T19:59:22",
						},
					}, {
						input: `fs:"2012-08-2116:00:00"`,
						expectedQuery: SingleQuery{
							Token{
								FS,
								"fs",
							},
							EQUAL,
							"2012-08-2116:00:00",
						},
					}, {
						input: `fs:"2012-08-2116:59:22"`,
						expectedQuery: SingleQuery{
							Token{
								FS,
								"fs",
							},
							EQUAL,
							"2012-08-2116:59:22",
						},
					}, {
						input: `positives:5`,
						expectedQuery: SingleQuery{
							Token{
								POSITIVES,
								"positives",
							},
							EQUAL,
							"5",
						},
					}, {
						input: `positives:10+`,
						expectedQuery: SingleQuery{
							Token{
								POSITIVES,
								"positives",
							},
							GREATERTHAN,
							"10",
						},
					}, {
						input: `section:".k40s"`,
						expectedQuery: SingleQuery{
							Token{
								SECTION,
								"section",
							},
							EQUAL,
							".k40s",
						},
					}, {
						input: `imports:"crypt32.dll"`,
						expectedQuery: SingleQuery{
							Token{
								IMPORTS,
								"imports",
							},
							EQUAL,
							"crypt32.dll",
						},
					}, {
						input: `exports:"rtEncryptFolder"`,
						expectedQuery: SingleQuery{
							Token{
								EXPORTS,
								"exports",
							},
							EQUAL,
							"rtEncryptFolder",
						},
					},
				}

				for i, tt := range tests {
					l := NewLexer(tt.input)
					p := NewParser(l)
					queries, err := p.Parse()
					if err != nil {
						t.Fatalf("test[%d] failed with error=%q :", i, err)
					}
					for _, q := range queries.Statements {
						query := q.(SingleQuery)
						if query.Token != tt.expectedQuery.Token {
							t.Fatalf("test[%d] expected token=%q got token=%q", i, tt.expectedQuery.Token, query.Token)
						} else if query.Value != tt.expectedQuery.Value {
							t.Fatalf("test[%d] expected value=%q got value=%q", i, tt.expectedQuery.Value, query.Value)
						} else if query.Op != tt.expectedQuery.Op {
							t.Fatalf("test[%d] expected op=%q got op=%q", i, tt.expectedQuery.Op, query.Op)
						}
					}
				}
			})
			t.Run("TestParseMultiQuery", func(t *testing.T) {
				tests := []struct {
					input         string
					expectedQuery []SingleQuery
				}{
					{
						input: `type:peexe,size:200kb+`,
						expectedQuery: []SingleQuery{
							{
								Token{
									TYPE,
									"type",
								},
								EQUAL,
								"peexe",
							},
							{
								Token{
									SIZE,
									"size",
								},
								GREATERTHAN,
								"200000",
							},
						},
					}, {
						input: `type:elf,positives:30+,size:200kb+`,
						expectedQuery: []SingleQuery{
							{
								Token{
									TYPE,
									"type",
								},
								EQUAL,
								"elf",
							}, {
								Token{
									POSITIVES,
									"positives",
								},
								GREATERTHAN,
								"30",
							}, {
								Token{
									SIZE,
									"size",
								},
								GREATERTHAN,
								"200000",
							},
						},
					},
				}

				for i, tt := range tests {
					l := NewLexer(tt.input)
					p := NewParser(l)
					queries, err := p.Parse()
					if err != nil {
						t.Fatalf("test[%d] failed with error=%q :", i, err)
					}
					t.Log(queries)
					for k, q := range queries.Statements {
						query := q.(SingleQuery)
						t.Log(query)
						if query != tt.expectedQuery[k] {
							t.Fatalf("test[%d] failed expected=%q got=%q ", k, tt.expectedQuery[k], query)
						}
					}
				}
			})
		}

		func TestParserFailure(t *testing.T) {
			t.Run("TestParseFileSizeQueryWithError", func(t *testing.T) {
				tests := []struct {
					input         string
					err           error
					expectedQuery SingleQuery
				}{
					{
						input:         `size,:100`,
						expectedQuery: emptySingleQuery,
						err:           newSyntaxError(SIZE, ASSIGN, COMMA),
					}, {
						input:         `size:a100`,
						expectedQuery: emptySingleQuery,
						err:           newSyntaxError(SIZE, INT, IDENT),
					}, {
						input:         `size:250,kb`,
						expectedQuery: emptySingleQuery,
						err:           newSyntaxError(SIZE, "one of [+,-,kb,mb]", COMMA),
					},
				}

				for i, tt := range tests {
					l := NewLexer(tt.input)
					p := NewParser(l)
					query, err := p.parseFileSizeQuery()
					if err == nil {
						t.Fatalf("test[%d] expected failure with error :%q got nil", i, tt.err)
					}
					if err.Error() != tt.err.Error() {
						t.Fatalf("test[%d] expected error to be :%q got %q", i, tt.err, err)
					}
					if query != emptySingleQuery {
						t.Fatalf("test[%d] expected empty query :%q got %q", i, emptySingleQuery, query)
					}
				}
			})
			t.Run("TestParseFileTypeQueryWithError", func(t *testing.T) {
				tests := []struct {
					input         string
					err           error
					expectedQuery SingleQuery
				}{
					{
						input:         `type,:pe`,
						expectedQuery: emptySingleQuery,
						err:           newSyntaxError(TYPE, ASSIGN, COMMA),
					}, {
						input:         `type:12`,
						expectedQuery: emptySingleQuery,
						err:           newSyntaxError(TYPE, IDENT, INT),
					}, {
						input:         `type:,pe`,
						expectedQuery: emptySingleQuery,
						err:           newSyntaxError(TYPE, IDENT, COMMA),
					},
				}

				for i, tt := range tests {
					l := NewLexer(tt.input)
					p := NewParser(l)
					query, err := p.parseFileTypeQuery()
					if err == nil {
						t.Fatalf("test[%d] expected failure with error :%q got nil", i, tt.err)
					}
					if err.Error() != tt.err.Error() {
						t.Fatalf("test[%d] expected error to be :%q got %q", i, tt.err, err)
					}
					if query != emptySingleQuery {
						t.Fatalf("test[%d] expected empty query :%q got %q", i, emptySingleQuery, query)
					}
				}
			})
			t.Run("TestParseDatetimeQueryWithError", func(t *testing.T) {
				tests := []struct {
					input         string
					err           error
					expectedQuery SingleQuery
				}{
					{
						input:         `fs,:"2012-08-2116:59:22"`,
						expectedQuery: emptySingleQuery,
						err:           newSyntaxError("DATE", ASSIGN, COMMA),
					}, {
						input:         `ls:123"2012-08-2116:59:22"`,
						expectedQuery: emptySingleQuery,
						err:           newSyntaxError("DATE", STRING, INT),
					},
				}

				for i, tt := range tests {
					l := NewLexer(tt.input)
					p := NewParser(l)
					query, err := p.parseDatetimeQuery()
					if err == nil {
						t.Fatalf("test[%d] expected failure with error :%q got nil", i, tt.err)
					}
					if err.Error() != tt.err.Error() {
						t.Fatalf("test[%d] expected error to be :%q got %q", i, tt.err, err)
					}
					if query != emptySingleQuery {
						t.Fatalf("test[%d] expected empty query :%q got %q", i, emptySingleQuery, query)
					}
				}
			})
			t.Run("TestParseQueryWithError", func(t *testing.T) {
				tests := []struct {
					input         string
					err           error
					expectedQuery SingleQuery
				}{
					{
						input:         `ano,:pe`,
						expectedQuery: emptySingleQuery,
						err:           errBadQuerySyntax,
					},
				}

				for i, tt := range tests {
					l := NewLexer(tt.input)
					p := NewParser(l)
					_, err := p.Parse()
					if err == nil {
						t.Fatalf("test[%d] expected failure with error :%q got nil", i, tt.err)
					}
					if err.Error() != tt.err.Error() {
						t.Fatalf("test[%d] expected error to be :%q got %q", i, tt.err, err)
					}
				}
			})
			t.Run("TestParseSectionQueryWithError", func(t *testing.T) {
				tests := []struct {
					input         string
					err           error
					expectedQuery SingleQuery
				}{
					{
						input:         `section,:".k40s"`,
						expectedQuery: emptySingleQuery,
						err:           newSyntaxError(SECTION, ASSIGN, COMMA),
					}, {
						input:         `section:23".k40s"`,
						expectedQuery: emptySingleQuery,
						err:           newSyntaxError(SECTION, STRING, INT),
					},
				}

				for i, tt := range tests {
					l := NewLexer(tt.input)
					p := NewParser(l)
					query, err := p.parseSectionQuery()
					if err == nil {
						t.Fatalf("test[%d] expected failure with error :%q got nil", i, tt.err)
					}
					if err.Error() != tt.err.Error() {
						t.Fatalf("test[%d] expected error to be :%q got %q", i, tt.err, err)
					}
					if query != emptySingleQuery {
						t.Fatalf("test[%d] expected empty query :%q got %q", i, emptySingleQuery, query)
					}
				}
			})
			t.Run("TestParsePositivesQueryWithError", func(t *testing.T) {
				tests := []struct {
					input         string
					err           error
					expectedQuery SingleQuery
				}{
					{
						input:         `positives,:100`,
						expectedQuery: emptySingleQuery,
						err:           newSyntaxError(POSITIVES, ASSIGN, COMMA),
					}, {
						input:         `positives:a100`,
						expectedQuery: emptySingleQuery,
						err:           newSyntaxError(POSITIVES, INT, IDENT),
					},
				}

				for i, tt := range tests {
					l := NewLexer(tt.input)
					p := NewParser(l)
					query, err := p.parsePositivesQuery()
					if err == nil {
						t.Fatalf("test[%d] expected failure with error :%q got nil", i, tt.err)
					}
					if err.Error() != tt.err.Error() {
						t.Fatalf("test[%d] expected error to be :%q got %q", i, tt.err, err)
					}
					if query != emptySingleQuery {
						t.Fatalf("test[%d] expected empty query :%q got %q", i, emptySingleQuery, query)
					}
				}
			})
		}
		func BenchmarkParser(b *testing.B) {
			tests := []struct {
				input         string
				expectedQuery []SingleQuery
			}{
				{
					input: `type:elf,positives:30+,size:200kb+`,
					expectedQuery: []SingleQuery{
						{
							Token{
								TYPE,
								"type",
							},
							EQUAL,
							"elf",
						}, {
							Token{
								POSITIVES,
								"positives",
							},
							GREATERTHAN,
							"30",
						}, {
							Token{
								SIZE,
								"size",
							},
							GREATERTHAN,
							"200000",
						},
					},
				},
			}

			for _, tt := range tests {
				l := NewLexer(tt.input)
				p := NewParser(l)
				p.parseFileSizeQuery()

			}
	*/
}
