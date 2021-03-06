package query

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

var viewLoadTests = []struct {
	Name          string
	Encoding      cmd.Encoding
	NoHeader      bool
	From          parser.FromClause
	UseInternalId bool
	Stdin         string
	Filter        *Filter
	Result        *View
	Error         string
}{
	{
		Name: "Dual View",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Dual{}},
			},
		},
		Result: &View{
			Header: []HeaderField{{}},
			RecordSet: []Record{
				{
					NewCell(value.NewNull()),
				},
			},
			Filter: &Filter{
				Variables:    []VariableMap{{}},
				TempViews:    []ViewMap{{}},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases:      AliasNodes{{}},
			},
		},
	},
	{
		Name: "Dual View With Omitted FromClause",
		From: parser.FromClause{},
		Result: &View{
			Header: []HeaderField{{}},
			RecordSet: []Record{
				{
					NewCell(value.NewNull()),
				},
			},
			Filter: &Filter{
				Variables:    []VariableMap{{}},
				TempViews:    []ViewMap{{}},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases:      AliasNodes{{}},
			},
		},
	},
	{
		Name: "Load File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table1.csv"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table1.csv",
				Delimiter: ',',
			},
			Filter: &Filter{
				Variables:    []VariableMap{{}},
				TempViews:    []ViewMap{{}},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Load with Parentheses",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Parentheses{
					Expr: parser.Table{
						Object: parser.Identifier{Literal: "table1.csv"},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table1.csv",
				Delimiter: ',',
			},
			Filter: &Filter{
				Variables:    []VariableMap{{}},
				TempViews:    []ViewMap{{}},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Load From Stdin",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{Stdin: "stdin"}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin: "column1,column2\n1,\"str1\"",
		Result: &View{
			Header: NewHeader("t", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "stdin",
				Delimiter: ',',
			},
			Filter: &Filter{
				Variables: []VariableMap{{}},
				TempViews: []ViewMap{
					{
						"STDIN": nil,
					},
				},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"T": "STDIN",
					},
				},
			},
		},
	},
	{
		Name: "Load From Stdin With Internal Id",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{Stdin: "stdin"}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		UseInternalId: true,
		Stdin:         "column1,column2\n1,\"str1\"",
		Result: &View{
			Header: NewHeaderWithId("t", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(0, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "stdin",
				Delimiter: ',',
			},
			Filter: &Filter{
				Variables: []VariableMap{{}},
				TempViews: []ViewMap{
					{
						"STDIN": nil,
					},
				},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"T": "STDIN",
					},
				},
			},
			UseInternalId: true,
		},
	},
	{
		Name:  "Load From Stdin With Omitted FromClause",
		From:  parser.FromClause{},
		Stdin: "column1,column2\n1,\"str1\"",
		Result: &View{
			Header: NewHeader("stdin", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "stdin",
				Delimiter: ',',
			},
			Filter: &Filter{
				Variables: []VariableMap{{}},
				TempViews: []ViewMap{
					{
						"STDIN": nil,
					},
				},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"STDIN": "STDIN",
					},
				},
			},
		},
	},
	{
		Name: "Load From Stdin Broken CSV Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{Stdin: "stdin"}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin: "column1,column2\n1\"str1\"",
		Error: "[L:- C:-] csv parse error in file stdin: line 1, column 8: wrong number of fields in line",
	},
	{
		Name: "Load From Stdin Duplicate Table Name Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}, Alias: parser.Identifier{Literal: "t"}},
				parser.Table{Object: parser.Stdin{Stdin: "stdin"}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin: "column1,column2\n1,\"str1\"",
		Error: "[L:- C:-] table name t is a duplicate",
	},
	{
		Name: "Stdin Empty Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Stdin{Stdin: "stdin"},
				},
			},
		},
		Error: "[L:- C:-] stdin is empty",
	},
	{
		Name: "Load File Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "notexist"},
				},
			},
		},
		Error: "[L:- C:-] file notexist does not exist",
	},
	{
		Name: "Load From File Duplicate Table Name Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}, Alias: parser.Identifier{Literal: "t"}},
				parser.Table{Object: parser.Identifier{Literal: "table2"}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Error: "[L:- C:-] table name t is a duplicate",
	},
	{
		Name: "Load From File Inline Table",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "it"}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Filter: &Filter{
			Variables: VariableScopes{{}},
			TempViews: TemporaryViewScopes{{}},
			Cursors:   CursorScopes{{}},
			InlineTables: InlineTableNodes{
				InlineTableMap{
					"IT": &View{
						Header: NewHeader("it", []string{"c1", "c2", "num"}),
						RecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("str1"),
								value.NewInteger(1),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("str2"),
								value.NewInteger(1),
							}),
							NewRecord([]value.Primary{
								value.NewString("3"),
								value.NewString("str3"),
								value.NewInteger(1),
							}),
						},
					},
				},
			},
			Aliases: AliasNodes{{}},
		},
		Result: &View{
			Header: NewHeader("t", []string{"c1", "c2", "num"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewInteger(1),
				}),
			},
			Filter: &Filter{
				Variables: VariableScopes{{}},
				TempViews: TemporaryViewScopes{{}},
				Cursors:   CursorScopes{{}},
				InlineTables: InlineTableNodes{
					{},
					InlineTableMap{
						"IT": &View{
							Header: NewHeader("it", []string{"c1", "c2", "num"}),
							RecordSet: []Record{
								NewRecord([]value.Primary{
									value.NewString("1"),
									value.NewString("str1"),
									value.NewInteger(1),
								}),
								NewRecord([]value.Primary{
									value.NewString("2"),
									value.NewString("str2"),
									value.NewInteger(1),
								}),
								NewRecord([]value.Primary{
									value.NewString("3"),
									value.NewString("str3"),
									value.NewInteger(1),
								}),
							},
						},
					},
				},
				Aliases: AliasNodes{
					{
						"T": "",
					},
					{},
				},
			},
		},
	},
	{
		Name: "Load From File Inline Table Duplicate Table Name Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}, Alias: parser.Identifier{Literal: "t"}},
				parser.Table{Object: parser.Identifier{Literal: "it"}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Filter: &Filter{
			Variables: VariableScopes{{}},
			TempViews: TemporaryViewScopes{{}},
			Cursors:   CursorScopes{{}},
			InlineTables: InlineTableNodes{
				InlineTableMap{
					"IT": &View{
						Header: NewHeader("it", []string{"c1", "c2", "num"}),
						RecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("str1"),
								value.NewInteger(1),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("str2"),
								value.NewInteger(1),
							}),
							NewRecord([]value.Primary{
								value.NewString("3"),
								value.NewString("str3"),
								value.NewInteger(1),
							}),
						},
					},
				},
			},
			Aliases: AliasNodes{{}},
		},
		Error: "[L:- C:-] table name t is a duplicate",
	},
	{
		Name:     "Load SJIS File",
		Encoding: cmd.SJIS,
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table_sjis"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table_sjis", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("日本語"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str"),
				}),
			},
			Filter: &Filter{
				Variables:    []VariableMap{{}},
				TempViews:    []ViewMap{{}},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"TABLE_SJIS": strings.ToUpper(GetTestFilePath("table_sjis.csv")),
					},
				},
			},
		},
	},
	{
		Name:     "Load No Header File",
		NoHeader: true,
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table_noheader"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table_noheader", []string{"c1", "c2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			Filter: &Filter{
				Variables:    []VariableMap{{}},
				TempViews:    []ViewMap{{}},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"TABLE_NOHEADER": strings.ToUpper(GetTestFilePath("table_noheader.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Load Multiple File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table1"},
				},
				parser.Table{
					Object: parser.Identifier{Literal: "table2"},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("4"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("4"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("4"),
					value.NewString("str44"),
				}),
			},
			Filter: &Filter{
				Variables:    []VariableMap{{}},
				TempViews:    []ViewMap{{}},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
						"TABLE2": strings.ToUpper(GetTestFilePath("table2.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Cross Join",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
						},
						JoinType: parser.Token{Token: parser.CROSS, Literal: "cross"},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("4"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("4"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("4"),
					value.NewString("str44"),
				}),
			},
			Filter: &Filter{
				Variables:    []VariableMap{{}},
				TempViews:    []ViewMap{{}},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
						"TABLE2": strings.ToUpper(GetTestFilePath("table2.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Inner Join",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
						},
						Condition: parser.JoinCondition{
							On: parser.Comparison{
								LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "column1"}},
								RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column3"}},
								Operator: "=",
							},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
			},
			Filter: &Filter{
				Variables:    []VariableMap{{}},
				TempViews:    []ViewMap{{}},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
						"TABLE2": strings.ToUpper(GetTestFilePath("table2.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Inner Join Using Condition",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table1b"},
						},
						Condition: parser.JoinCondition{
							Using: []parser.QueryExpression{
								parser.Identifier{Literal: "column1"},
							},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Column: "column1", IsFromTable: true, IsJoinColumn: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table1b", Column: "column2b", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("str1b"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("str2b"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("str3b"),
				}),
			},
			Filter: &Filter{
				Variables:    []VariableMap{{}},
				TempViews:    []ViewMap{{}},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"TABLE1":  strings.ToUpper(GetTestFilePath("table1.csv")),
						"TABLE1B": strings.ToUpper(GetTestFilePath("table1b.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Outer Join",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
						},
						Direction: parser.Token{Token: parser.LEFT, Literal: "left"},
						Condition: parser.JoinCondition{
							On: parser.Comparison{
								LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "column1"}},
								RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column3"}},
								Operator: "=",
							},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewNull(),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
			},
			Filter: &Filter{
				Variables:    []VariableMap{{}},
				TempViews:    []ViewMap{{}},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
						"TABLE2": strings.ToUpper(GetTestFilePath("table2.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Outer Join Natural",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table1b"},
						},
						Direction: parser.Token{Token: parser.RIGHT, Literal: "right"},
						Natural:   parser.Token{Token: parser.NATURAL, Literal: "natural"},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Column: "column1", IsFromTable: true, IsJoinColumn: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table1b", Column: "column2b", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("str1b"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("str2b"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("str3b"),
				}),
				NewRecord([]value.Primary{
					value.NewString("4"),
					value.NewNull(),
					value.NewString("str4b"),
				}),
			},
			Filter: &Filter{
				Variables:    []VariableMap{{}},
				TempViews:    []ViewMap{{}},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"TABLE1":  strings.ToUpper(GetTestFilePath("table1.csv")),
						"TABLE1B": strings.ToUpper(GetTestFilePath("table1b.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Join Left Side Table File Not Exist Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "notexist"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
						},
						JoinType: parser.Token{Token: parser.CROSS, Literal: "cross"},
					},
				},
			},
		},
		Error: "[L:- C:-] file notexist does not exist",
	},
	{
		Name: "Join Right Side Table File Not Exist Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "notexist"},
						},
						JoinType: parser.Token{Token: parser.CROSS, Literal: "cross"},
					},
				},
			},
		},
		Error: "[L:- C:-] file notexist does not exist",
	},
	{
		Name: "Load Subquery",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Subquery{
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Select: "select",
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
										parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
									},
								},
								FromClause: parser.FromClause{
									Tables: []parser.QueryExpression{
										parser.Table{Object: parser.Identifier{Literal: "table1"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "alias"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("alias", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			Filter: &Filter{
				Variables:    []VariableMap{{}},
				TempViews:    []ViewMap{{}},
				Cursors:      []CursorMap{{}},
				InlineTables: InlineTableNodes{{}},
				Aliases: AliasNodes{
					{
						"ALIAS": "",
					},
				},
			},
		},
	},
	{
		Name: "Load Subquery Duplicate Table Name Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}, Alias: parser.Identifier{Literal: "t"}},
				parser.Table{
					Object: parser.Subquery{
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Select: "select",
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
										parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
									},
								},
								FromClause: parser.FromClause{
									Tables: []parser.QueryExpression{
										parser.Table{Object: parser.Identifier{Literal: "table1"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "[L:- C:-] table name t is a duplicate",
	},
	{
		Name: "Load CSV Parse Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table_broken.csv"},
				},
			},
		},
		Error: fmt.Sprintf("[L:- C:-] csv parse error in file %s: line 3, column 7: wrong number of fields in line", GetTestFilePath("table_broken.csv")),
	},
	{
		Name: "Inner Join Join Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
						},
						Condition: parser.JoinCondition{
							On: parser.Comparison{
								LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "notexist"}},
								RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column3"}},
								Operator: "=",
							},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] field table1.notexist does not exist",
	},
	{
		Name: "Outer Join Join Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
						},
						Direction: parser.Token{Token: parser.LEFT, Literal: "left"},
						Condition: parser.JoinCondition{
							On: parser.Comparison{
								LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "column1"}},
								RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "notexist"}},
								Operator: "=",
							},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] field table2.notexist does not exist",
	},
	{
		Name: "Inner Join Using Condition Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table1b"},
						},
						Condition: parser.JoinCondition{
							Using: []parser.QueryExpression{
								parser.Identifier{Literal: "notexist"},
							},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestView_Load(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range viewLoadTests {
		ViewCache.Clean()

		tf.Delimiter = cmd.UNDEF
		tf.NoHeader = v.NoHeader
		if v.Encoding != "" {
			tf.Encoding = v.Encoding
		} else {
			tf.Encoding = cmd.UTF8
		}

		var oldStdin *os.File
		if 0 < len(v.Stdin) {
			oldStdin = os.Stdin
			r, w, _ := os.Pipe()
			w.WriteString(v.Stdin)
			w.Close()
			os.Stdin = r
		}

		view := NewView()
		if v.Filter == nil {
			v.Filter = NewEmptyFilter()
		}
		view.UseInternalId = v.UseInternalId

		err := view.Load(v.From, v.Filter.CreateNode())

		if 0 < len(v.Stdin) {
			os.Stdin = oldStdin
		}

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}

		if v.Result.FileInfo != nil {
			if filepath.Base(view.FileInfo.Path) != filepath.Base(v.Result.FileInfo.Path) {
				t.Errorf("%s: filepath = %q, want %q", v.Name, filepath.Base(view.FileInfo.Path), filepath.Base(v.Result.FileInfo.Path))
			}
			if view.FileInfo.Delimiter != v.Result.FileInfo.Delimiter {
				t.Errorf("%s: delimiter = %q, want %q", v.Name, view.FileInfo.Delimiter, v.Result.FileInfo.Delimiter)
			}
		}
		view.FileInfo = nil
		v.Result.FileInfo = nil

		if !reflect.DeepEqual(view.Filter.Aliases, v.Result.Filter.Aliases) {
			t.Errorf("%s: alias list = %q, want %q", v.Name, view.Filter.Aliases, v.Result.Filter.Aliases)
		}
		for i, tviews := range v.Result.Filter.TempViews {
			resultKeys := make([]string, len(tviews))
			for key := range tviews {
				resultKeys = append(resultKeys, key)
			}
			viewKeys := make([]string, len(view.Filter.TempViews[i]))
			for key := range view.Filter.TempViews[i] {
				viewKeys = append(viewKeys, key)
			}
			if !reflect.DeepEqual(resultKeys, viewKeys) {
				t.Errorf("%s: temp view list = %v, want %v", v.Name, view.Filter.TempViews, v.Result.Filter.TempViews)
			}
		}

		view.Filter = NewEmptyFilter()
		v.Result.Filter = NewEmptyFilter()

		if !reflect.DeepEqual(view, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, view, v.Result)
		}
	}
}

func TestNewViewFromGroupedRecord(t *testing.T) {
	fr := FilterRecord{
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2), value.NewInteger(3)}),
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2), value.NewInteger(3)}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str2"), value.NewString("str3")}),
				},
			},
		},
		RecordIndex: 0,
	}
	expect := &View{
		Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
		RecordSet: []Record{
			{NewCell(value.NewInteger(1)), NewCell(value.NewInteger(1)), NewCell(value.NewString("str1"))},
			{NewCell(value.NewInteger(2)), NewCell(value.NewInteger(2)), NewCell(value.NewString("str2"))},
			{NewCell(value.NewInteger(3)), NewCell(value.NewInteger(3)), NewCell(value.NewString("str3"))},
		},
	}

	result := NewViewFromGroupedRecord(fr)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, want %v", result, expect)
	}
}

var viewWhereTests = []struct {
	Name   string
	CPU    int
	View   *View
	Where  parser.WhereClause
	Result RecordSet
	Error  string
}{
	{
		Name: "Where",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Where: parser.WhereClause{
			Filter: parser.Comparison{
				LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				RHS:      parser.NewIntegerValueFromString("2"),
				Operator: "=",
			},
		},
		Result: RecordSet{
			NewRecordWithId(2, []value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	},
	{
		Name: "Where in Multi Threading",
		CPU:  3,
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Where: parser.WhereClause{
			Filter: parser.Comparison{
				LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				RHS:      parser.NewIntegerValueFromString("2"),
				Operator: "=",
			},
		},
		Result: RecordSet{
			NewRecordWithId(2, []value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	},
	{
		Name: "Where Filter Error",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Where: parser.WhereClause{
			Filter: parser.Comparison{
				LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				RHS:      parser.NewIntegerValueFromString("2"),
				Operator: "=",
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestView_Where(t *testing.T) {
	flags := cmd.GetFlags()

	for _, v := range viewWhereTests {
		flags.CPU = 1
		if v.CPU != 0 {
			flags.CPU = v.CPU
		}

		err := v.View.Where(v.Where)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if !reflect.DeepEqual(v.View.RecordSet, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, v.View.RecordSet, v.Result)
		}
	}
}

var viewGroupByTests = []struct {
	Name       string
	View       *View
	GroupBy    parser.GroupByClause
	Result     *View
	IsGrouped  bool
	GroupItems []string
	Error      string
}{
	{
		Name: "Group By",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("group1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("group1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		GroupBy: parser.GroupByClause{
			Items: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: INTERNAL_ID_COLUMN,
				},
				{
					View:        "table1",
					Column:      "column1",
					Number:      1,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					Number:      2,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					Number:      3,
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(3)}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str3")}),
					NewGroupCell([]value.Primary{value.NewString("group1"), value.NewString("group1")}),
				},
				{
					NewGroupCell([]value.Primary{value.NewInteger(2), value.NewInteger(4)}),
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
					NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
				},
			},
			Filter:    NewEmptyFilter(),
			isGrouped: true,
		},
	},
	{
		Name: "Group By With ColumnNumber",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("group1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("group1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		GroupBy: parser.GroupByClause{
			Items: []parser.QueryExpression{
				parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(3)},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: INTERNAL_ID_COLUMN,
				},
				{
					View:        "table1",
					Column:      "column1",
					Number:      1,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					Number:      2,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					Number:      3,
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(3)}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str3")}),
					NewGroupCell([]value.Primary{value.NewString("group1"), value.NewString("group1")}),
				},
				{
					NewGroupCell([]value.Primary{value.NewInteger(2), value.NewInteger(4)}),
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
					NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
				},
			},
			Filter:    NewEmptyFilter(),
			isGrouped: true,
		},
	},
	{
		Name: "Group By Evaluation Error",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("group1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("group1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		GroupBy: parser.GroupByClause{
			Items: []parser.QueryExpression{
				parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(0)},
			},
		},
		Error: "[L:- C:-] field table1.0 does not exist",
	},
	{
		Name: "Group By Empty Record",
		View: &View{
			Header:    NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{},
			Filter:    NewEmptyFilter(),
		},
		GroupBy: parser.GroupByClause{
			Items: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: INTERNAL_ID_COLUMN,
				},
				{
					View:        "table1",
					Column:      "column1",
					Number:      1,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					Number:      2,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					Number:      3,
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			RecordSet: []Record{},
			Filter:    NewEmptyFilter(),
			isGrouped: true,
		},
	},
	{
		Name: "Group By Empty Record with No Condition",
		View: &View{
			Header:    NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{},
			Filter:    NewEmptyFilter(),
		},
		GroupBy: parser.GroupByClause{
			Items: nil,
		},
		Result: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: INTERNAL_ID_COLUMN,
				},
				{
					View:        "table1",
					Column:      "column1",
					Number:      1,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					Number:      2,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					Number:      3,
					IsFromTable: true,
				},
			},
			RecordSet: []Record{},
			Filter:    NewEmptyFilter(),
			isGrouped: true,
		},
	},
}

func TestView_GroupBy(t *testing.T) {
	for _, v := range viewGroupByTests {
		err := v.View.GroupBy(v.GroupBy)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if !reflect.DeepEqual(v.View, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, v.View, v.Result)
		}
	}
}

var viewHavingTests = []struct {
	Name   string
	View   *View
	Having parser.HavingClause
	Result RecordSet
	Error  string
}{
	{
		Name: "Having",
		View: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: INTERNAL_ID_COLUMN,
				},
				{
					View:        "table1",
					Column:      "column1",
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			isGrouped: true,
			RecordSet: RecordSet{
				{
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str3")}),
					NewGroupCell([]value.Primary{value.NewString("group1"), value.NewString("group1")}),
				},
				{
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
					NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
				},
			},
			Filter: NewEmptyFilter(),
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.AggregateFunction{
					Name:     "sum",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				RHS:      parser.NewIntegerValueFromString("5"),
				Operator: ">",
			},
		},
		Result: RecordSet{
			{
				NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
				NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
				NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
				NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
			},
		},
	},
	{
		Name: "Having Filter Error",
		View: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: INTERNAL_ID_COLUMN,
				},
				{
					View:        "table1",
					Column:      "column1",
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			isGrouped: true,
			RecordSet: RecordSet{
				{
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str3")}),
					NewGroupCell([]value.Primary{value.NewString("group1"), value.NewString("group1")}),
				},
				{
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
					NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
				},
			},
			Filter: NewEmptyFilter(),
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.AggregateFunction{
					Name:     "sum",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
				RHS:      parser.NewIntegerValueFromString("5"),
				Operator: ">",
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Having Not Grouped",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: RecordSet{
				NewRecordWithId(1, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.AggregateFunction{
					Name:     "sum",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				RHS:      parser.NewIntegerValueFromString("5"),
				Operator: ">",
			},
		},
		Result: RecordSet{
			{
				NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2)}),
				NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
				NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
				NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
			},
		},
	},
	{
		Name: "Having All RecordSet Filter Error",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: RecordSet{
				NewRecordWithId(1, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.AggregateFunction{
					Name:     "sum",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
				RHS:      parser.NewIntegerValueFromString("5"),
				Operator: ">",
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestView_Having(t *testing.T) {
	for _, v := range viewHavingTests {
		err := v.View.Having(v.Having)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if !reflect.DeepEqual(v.View.RecordSet, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, v.View.RecordSet, v.Result)
		}
	}
}

var viewSelectTests = []struct {
	Name   string
	View   *View
	Select parser.SelectClause
	Result *View
	Error  string
}{
	{
		Name: "Select",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: INTERNAL_ID_COLUMN},
				{View: "table2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(3),
					value.NewString("1"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}, Alias: parser.Identifier{Literal: "c2"}},
				parser.Field{Object: parser.AllColumns{}},
				parser.Field{Object: parser.NewIntegerValueFromString("1"), Alias: parser.Identifier{Literal: "a"}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}, Alias: parser.Identifier{Literal: "c2a"}},
				parser.Field{Object: parser.ColumnNumber{View: parser.Identifier{Literal: "table2"}, Number: value.NewInteger(1)}, Alias: parser.Identifier{Literal: "t21"}},
				parser.Field{Object: parser.ColumnNumber{View: parser.Identifier{Literal: "table2"}, Number: value.NewInteger(1)}, Alias: parser.Identifier{Literal: "t21a"}},
				parser.Field{Object: parser.PrimitiveType{
					Literal: "2012-01-01",
					Value:   value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Aliases: []string{"c2", "c2a"}, Number: 2, IsFromTable: true},
				{View: "table2", Column: INTERNAL_ID_COLUMN},
				{View: "table2", Column: "column3", Aliases: []string{"t21", "t21a"}, Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
				{Column: "1", Aliases: []string{"a"}},
				{Column: "2012-01-01T00:00:00Z"},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewString("3"),
					value.NewString("str33"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(3),
					value.NewString("1"),
					value.NewString("str44"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
			},
			Filter:       NewEmptyFilter(),
			selectFields: []int{2, 1, 2, 4, 5, 6, 2, 4, 4, 7},
		},
	},
	{
		Name: "Select Distinct",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{View: "table2", Column: INTERNAL_ID_COLUMN},
				{View: "table2", Column: "column3", IsFromTable: true},
				{View: "table2", Column: "column4", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(3),
					value.NewString("4"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Select: parser.SelectClause{
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.NewIntegerValueFromString("1"), Alias: parser.Identifier{Literal: "a"}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", IsFromTable: true},
				{Column: "1", Aliases: []string{"a"}},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewInteger(1),
				}),
			},
			Filter:       NewEmptyFilter(),
			selectFields: []int{0, 1},
		},
	},
	{
		Name: "Select Aggregate Function",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{
					Object: parser.AggregateFunction{
						Name:     "sum",
						Distinct: parser.Token{},
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{Column: "sum(column1)"},
			},
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2)}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("2")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str2")}),
					NewCell(value.NewInteger(3)),
				},
			},
			Filter:       NewEmptyFilter(),
			selectFields: []int{3},
		},
	},
	{
		Name: "Select Aggregate Function Not Group Key Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AggregateFunction{
						Name:     "sum",
						Distinct: parser.Token{},
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] field column2 is not a group key",
	},
	{
		Name: "Select Aggregate Function All RecordSet Lazy Evaluation",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.NewIntegerValueFromString("1")},
				parser.Field{
					Object: parser.Arithmetic{
						LHS: parser.AggregateFunction{
							Name:     "sum",
							Distinct: parser.Token{},
							Args: []parser.QueryExpression{
								parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
							},
						},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: '+',
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{Column: "1"},
				{Column: "sum(column1) + 1"},
			},
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2)}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("2")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str2")}),
					NewCell(value.NewInteger(1)),
					NewCell(value.NewInteger(4)),
				},
			},
			Filter:       NewEmptyFilter(),
			selectFields: []int{3, 4},
		},
	},
	{
		Name: "Select Analytic Function",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "row_number",
						Over: "over",
						AnalyticClause: parser.AnalyticClause{
							PartitionClause: parser.PartitionClause{
								PartitionBy: "partition by",
								Values: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								OrderBy: "order by",
								Items: []parser.QueryExpression{
									parser.OrderItem{
										Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "rownum"},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{Column: "row_number() over (partition by column1 order by column2)", Aliases: []string{"rownum"}},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
					value.NewInteger(3),
				}),
			},
			Filter:       NewEmptyFilter(),
			selectFields: []int{0, 1, 2},
		},
	},
	{
		Name: "Select Analytic Function Not Exist Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "notexist",
						Over: "over",
						AnalyticClause: parser.AnalyticClause{
							PartitionClause: parser.PartitionClause{
								PartitionBy: "partition by",
								Values: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								OrderBy: "order by",
								Items: []parser.QueryExpression{
									parser.OrderItem{
										Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "rownum"},
				},
			},
		},
		Error: "[L:- C:-] function notexist does not exist",
	},
	{
		Name: "Select Analytic Function Partition Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "row_number",
						Over: "over",
						AnalyticClause: parser.AnalyticClause{
							PartitionClause: parser.PartitionClause{
								PartitionBy: "partition by",
								Values: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								OrderBy: "order by",
								Items: []parser.QueryExpression{
									parser.OrderItem{
										Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "rownum"},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Select Analytic Function Order Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "row_number",
						Over: "over",
						AnalyticClause: parser.AnalyticClause{
							PartitionClause: parser.PartitionClause{
								PartitionBy: "partition by",
								Values: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								OrderBy: "order by",
								Items: []parser.QueryExpression{
									parser.OrderItem{
										Value: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "rownum"},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Select User Defined Analytic Function",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
			Filter: &Filter{
				Functions: UserDefinedFunctionScopes{
					UserDefinedFunctionMap{
						"USERAGGFUNC": &UserDefinedFunction{
							Name:         parser.Identifier{Literal: "useraggfunc"},
							IsAggregate:  true,
							Cursor:       parser.Identifier{Literal: "list"},
							RequiredArgs: 0,
							Statements: []parser.Statement{
								parser.VariableDeclaration{
									Assignments: []parser.VariableAssignment{
										{
											Variable: parser.Variable{Name: "@value"},
										},
										{
											Variable: parser.Variable{Name: "@fetch"},
										},
									},
								},
								parser.WhileInCursor{
									Variables: []parser.Variable{
										{Name: "@fetch"},
									},
									Cursor: parser.Identifier{Literal: "list"},
									Statements: []parser.Statement{
										parser.If{
											Condition: parser.Is{
												LHS: parser.Variable{Name: "@value"},
												RHS: parser.NewNullValue(),
											},
											Statements: []parser.Statement{
												parser.VariableSubstitution{
													Variable: parser.Variable{Name: "@value"},
													Value:    parser.Variable{Name: "@fetch"},
												},
												parser.FlowControl{Token: parser.CONTINUE},
											},
										},
										parser.VariableSubstitution{
											Variable: parser.Variable{Name: "@value"},
											Value: parser.Arithmetic{
												LHS:      parser.Variable{Name: "@value"},
												RHS:      parser.Variable{Name: "@fetch"},
												Operator: '+',
											},
										},
									},
								},

								parser.Return{
									Value: parser.Variable{Name: "@value"},
								},
							},
						},
					},
				},
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "useraggfunc",
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
						},
						Over:           "over",
						AnalyticClause: parser.AnalyticClause{},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{Column: "useraggfunc(column2) over ()"},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
					value.NewInteger(15),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
					value.NewInteger(15),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
					value.NewInteger(15),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
					value.NewInteger(15),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
					value.NewInteger(15),
				}),
			},
			selectFields: []int{0, 1, 2},
		},
	},
}

func TestView_Select(t *testing.T) {
	for _, v := range viewSelectTests {
		err := v.View.Select(v.Select)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if !reflect.DeepEqual(v.View.Header, v.Result.Header) {
			t.Errorf("%s: header = %v, want %v", v.Name, v.View.Header, v.Result.Header)
		}
		if !reflect.DeepEqual(v.View.RecordSet, v.Result.RecordSet) {
			t.Errorf("%s: records = %s, want %s", v.Name, v.View.RecordSet, v.Result.RecordSet)
		}
		if !reflect.DeepEqual(v.View.selectFields, v.Result.selectFields) {
			t.Errorf("%s: select indices = %v, want %v", v.Name, v.View.selectFields, v.Result.selectFields)
		}
	}
}

var viewOrderByTests = []struct {
	Name    string
	View    *View
	OrderBy parser.OrderByClause
	Result  *View
	Error   string
}{
	{
		Name: "Order By",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{View: "table1", Column: "column3", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewString("2"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewNull(),
				}),
				NewRecordWithId(5, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
					value.NewString("4"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.QueryExpression{
				parser.OrderItem{
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				},
				parser.OrderItem{
					Value:     parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Direction: parser.Token{Token: parser.DESC, Literal: "desc"},
				},
				parser.OrderItem{
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
				},
				parser.OrderItem{
					Value: parser.NewIntegerValueFromString("1"),
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{View: "table1", Column: "column3", IsFromTable: true},
				{Column: "1"},
			},
			RecordSet: []Record{
				NewRecordWithId(5, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
					value.NewString("4"),
					value.NewInteger(1),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
					value.NewInteger(1),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
					value.NewInteger(1),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewNull(),
					value.NewInteger(1),
				}),
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewString("2"),
					value.NewInteger(1),
				}),
			},
			Filter: NewEmptyFilter(),
		},
	},
	{
		Name: "Order By with Cached SortValues",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewString("2"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewNull(),
				}),
				NewRecordWithId(5, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
					value.NewString("4"),
				}),
			},
			Filter: NewEmptyFilter(),
			sortValuesInEachCell: [][]*SortValue{
				{nil, nil, NewSortValue(value.NewString("3")), nil},
				{nil, nil, NewSortValue(value.NewString("4")), nil},
				{nil, nil, NewSortValue(value.NewString("4")), nil},
				{nil, nil, NewSortValue(value.NewString("3")), nil},
				{nil, nil, NewSortValue(value.NewString("2")), nil},
			},
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.QueryExpression{
				parser.OrderItem{
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(5, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
					value.NewString("4"),
				}),
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewString("2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewNull(),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
			},
			Filter: NewEmptyFilter(),
			sortValuesInEachCell: [][]*SortValue{
				{nil, nil, NewSortValue(value.NewString("2")), nil},
				{nil, nil, NewSortValue(value.NewString("3")), nil},
				{nil, nil, NewSortValue(value.NewString("3")), nil},
				{nil, nil, NewSortValue(value.NewString("4")), nil},
				{nil, nil, NewSortValue(value.NewString("4")), nil},
			},
		},
	},
	{
		Name: "Order By With Null Positions",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(5, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.QueryExpression{
				parser.OrderItem{
					Value:    parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					Position: parser.Token{Token: parser.LAST, Literal: "last"},
				},
				parser.OrderItem{
					Value:    parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Position: parser.Token{Token: parser.FIRST, Literal: "first"},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
				NewRecordWithId(5, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
	},
	{
		Name: "Order By Record Extend Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.QueryExpression{
				parser.OrderItem{
					Value: parser.AggregateFunction{
						Name:     "sum",
						Distinct: parser.Token{},
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] function sum cannot aggregate not grouping records",
	},
}

func TestView_OrderBy(t *testing.T) {
	for _, v := range viewOrderByTests {
		err := v.View.OrderBy(v.OrderBy)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if !reflect.DeepEqual(v.View.Header, v.Result.Header) {
			t.Errorf("%s: header = %v, want %v", v.Name, v.View.Header, v.Result.Header)
		}
		if !reflect.DeepEqual(v.View.RecordSet, v.Result.RecordSet) {
			t.Errorf("%s: records = %s, want %s", v.Name, v.View.RecordSet, v.Result.RecordSet)
		}
	}
}

var viewExtendRecordCapacity = []struct {
	Name   string
	View   *View
	Exprs  []parser.QueryExpression
	Result int
	Error  string
}{
	{
		Name: "ExtendRecordCapacity",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(2),
				}),
			},
			Filter: &Filter{
				Functions: UserDefinedFunctionScopes{
					UserDefinedFunctionMap{
						"USERFUNC": &UserDefinedFunction{
							Name: parser.Identifier{Literal: "userfunc"},
							Parameters: []parser.Variable{
								{Name: "@arg1"},
							},
							RequiredArgs: 1,
							Statements: []parser.Statement{
								parser.Return{Value: parser.Variable{Name: "@arg1"}},
							},
							IsAggregate: true,
						},
					},
				},
			},
			isGrouped: true,
		},
		Exprs: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.Function{
				Name: "userfunc",
				Args: []parser.QueryExpression{
					parser.NewIntegerValueFromString("1"),
				},
			},
			parser.AggregateFunction{
				Name:     "avg",
				Distinct: parser.Token{},
				Args: []parser.QueryExpression{
					parser.AggregateFunction{
						Name: "avg",
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
			},
			parser.ListAgg{
				ListAgg:  "listagg",
				Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
				Args: []parser.QueryExpression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					parser.NewStringValue(","),
				},
				OrderBy: parser.OrderByClause{
					Items: []parser.QueryExpression{
						parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
					},
				},
			},
			parser.AnalyticFunction{
				Name: "rank",
				AnalyticClause: parser.AnalyticClause{
					PartitionClause: parser.PartitionClause{
						Values: []parser.QueryExpression{
							parser.Arithmetic{
								LHS:      parser.NewIntegerValueFromString("1"),
								RHS:      parser.NewIntegerValueFromString("2"),
								Operator: '+',
							},
						},
					},
					OrderByClause: parser.OrderByClause{
						Items: []parser.QueryExpression{
							parser.OrderItem{
								Value: parser.Arithmetic{
									LHS:      parser.NewIntegerValueFromString("3"),
									RHS:      parser.NewIntegerValueFromString("4"),
									Operator: '+',
								},
							},
						},
					},
				},
			},
			parser.Arithmetic{
				LHS:      parser.NewIntegerValueFromString("5"),
				RHS:      parser.NewIntegerValueFromString("6"),
				Operator: '+',
			},
		},
		Result: 9,
	},
	{
		Name: "ExtendRecordCapacity UserDefinedFunction Not Grouped Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(2),
				}),
			},
			Filter: &Filter{
				Functions: UserDefinedFunctionScopes{
					UserDefinedFunctionMap{
						"USERFUNC": &UserDefinedFunction{
							Name: parser.Identifier{Literal: "userfunc"},
							Parameters: []parser.Variable{
								{Name: "@arg1"},
							},
							RequiredArgs: 1,
							Statements: []parser.Statement{
								parser.Return{Value: parser.Variable{Name: "@arg1"}},
							},
							IsAggregate: true,
						},
					},
				},
			},
		},
		Exprs: []parser.QueryExpression{
			parser.Function{
				Name: "userfunc",
				Args: []parser.QueryExpression{
					parser.NewIntegerValueFromString("1"),
				},
			},
		},
		Error: "[L:- C:-] function userfunc cannot aggregate not grouping records",
	},
	{
		Name: "ExtendRecordCapacity AggregateFunction Not Grouped Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(2),
				}),
			},
		},
		Exprs: []parser.QueryExpression{
			parser.AggregateFunction{
				Name: "avg",
				Args: []parser.QueryExpression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				},
			},
		},
		Error: "[L:- C:-] function avg cannot aggregate not grouping records",
	},
	{
		Name: "ExtendRecordCapacity ListAgg Not Grouped Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(2),
				}),
			},
		},
		Exprs: []parser.QueryExpression{
			parser.ListAgg{
				ListAgg:  "listagg",
				Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
				Args: []parser.QueryExpression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					parser.NewStringValue(","),
				},
				OrderBy: parser.OrderByClause{
					Items: []parser.QueryExpression{
						parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
					},
				},
			},
		},
		Error: "[L:- C:-] function listagg cannot aggregate not grouping records",
	},
	{
		Name: "ExtendRecordCapacity AnalyticFunction Partition Value Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(2),
				}),
			},
		},
		Exprs: []parser.QueryExpression{
			parser.AnalyticFunction{
				Name: "rank",
				AnalyticClause: parser.AnalyticClause{
					PartitionClause: parser.PartitionClause{
						Values: []parser.QueryExpression{
							parser.AggregateFunction{
								Name: "avg",
								Args: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
								},
							},
						},
					},
					OrderByClause: parser.OrderByClause{
						Items: []parser.QueryExpression{
							parser.OrderItem{
								Value: parser.Arithmetic{
									LHS:      parser.NewIntegerValueFromString("3"),
									RHS:      parser.NewIntegerValueFromString("4"),
									Operator: '+',
								},
							},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] function avg cannot aggregate not grouping records",
	},
	{
		Name: "ExtendRecordCapacity AnalyticFunction OrderBy Value Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(2),
				}),
			},
		},
		Exprs: []parser.QueryExpression{
			parser.AnalyticFunction{
				Name: "rank",
				AnalyticClause: parser.AnalyticClause{
					PartitionClause: parser.PartitionClause{
						Values: []parser.QueryExpression{
							parser.Arithmetic{
								LHS:      parser.NewIntegerValueFromString("1"),
								RHS:      parser.NewIntegerValueFromString("2"),
								Operator: '+',
							},
						},
					},
					OrderByClause: parser.OrderByClause{
						Items: []parser.QueryExpression{
							parser.OrderItem{
								Value: parser.AggregateFunction{
									Name: "avg",
									Args: []parser.QueryExpression{
										parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
									},
								},
							},
						},
					},
				},
			},
		},
		Error: "[L:- C:-] function avg cannot aggregate not grouping records",
	},
}

func TestView_ExtendRecordCapacity(t *testing.T) {
	for _, v := range viewExtendRecordCapacity {
		err := v.View.ExtendRecordCapacity(v.Exprs)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if cap(v.View.RecordSet[0]) != v.Result {
			t.Errorf("%s: record capacity = %d, want %d", v.Name, cap(v.View.RecordSet[0]), v.Result)
		}
	}
}

var viewLimitTests = []struct {
	Name   string
	View   *View
	Limit  parser.LimitClause
	Result *View
	Error  string
}{
	{
		Name: "Limit",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Limit: parser.LimitClause{Value: parser.NewIntegerValueFromString("2")},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
	},
	{
		Name: "Limit With Ties",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			Filter: NewEmptyFilter(),
			sortValuesInEachRecord: []SortValues{
				{
					&SortValue{Type: SORT_VALUE_INTEGER, Integer: 1},
					&SortValue{Type: SORT_VALUE_STRING, String: "str1"},
				},
				{
					&SortValue{Type: SORT_VALUE_INTEGER, Integer: 1},
					&SortValue{Type: SORT_VALUE_STRING, String: "str1"},
				},
				{
					&SortValue{Type: SORT_VALUE_INTEGER, Integer: 1},
					&SortValue{Type: SORT_VALUE_STRING, String: "str1"},
				},
				{
					&SortValue{Type: SORT_VALUE_INTEGER, Integer: 2},
					&SortValue{Type: SORT_VALUE_STRING, String: "str2"},
				},
				{
					&SortValue{Type: SORT_VALUE_INTEGER, Integer: 3},
					&SortValue{Type: SORT_VALUE_STRING, String: "str3"},
				},
			},
		},
		Limit: parser.LimitClause{Value: parser.NewIntegerValueFromString("2"), With: parser.LimitWith{Type: parser.Token{Token: parser.TIES}}},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			Filter: NewEmptyFilter(),
			sortValuesInEachRecord: []SortValues{
				{
					&SortValue{Type: SORT_VALUE_INTEGER, Integer: 1},
					&SortValue{Type: SORT_VALUE_STRING, String: "str1"},
				},
				{
					&SortValue{Type: SORT_VALUE_INTEGER, Integer: 1},
					&SortValue{Type: SORT_VALUE_STRING, String: "str1"},
				},
				{
					&SortValue{Type: SORT_VALUE_INTEGER, Integer: 1},
					&SortValue{Type: SORT_VALUE_STRING, String: "str1"},
				},
				{
					&SortValue{Type: SORT_VALUE_INTEGER, Integer: 2},
					&SortValue{Type: SORT_VALUE_STRING, String: "str2"},
				},
				{
					&SortValue{Type: SORT_VALUE_INTEGER, Integer: 3},
					&SortValue{Type: SORT_VALUE_STRING, String: "str3"},
				},
			},
		},
	},
	{
		Name: "Limit By Percentage",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			Filter: NewEmptyFilter(),
			offset: 1,
		},
		Limit: parser.LimitClause{Value: parser.NewFloatValue(50.5), Percent: "percent"},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			Filter: NewEmptyFilter(),
			offset: 1,
		},
	},
	{
		Name: "Limit By Over 100 Percentage",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Limit: parser.LimitClause{Value: parser.NewFloatValue(150), Percent: "percent"},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
	},
	{
		Name: "Limit By Negative Percentage",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Limit: parser.LimitClause{Value: parser.NewFloatValue(-10), Percent: "percent"},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{},
			Filter:    NewEmptyFilter(),
		},
	},
	{
		Name: "Limit Greater Than RecordSet",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Limit: parser.LimitClause{Value: parser.NewIntegerValueFromString("5")},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
	},
	{
		Name: "Limit Evaluate Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Limit: parser.LimitClause{Value: parser.Variable{Name: "notexist"}},
		Error: "[L:- C:-] variable notexist is undeclared",
	},
	{
		Name: "Limit Value Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Limit: parser.LimitClause{Value: parser.NewStringValue("str")},
		Error: "[L:- C:-] limit number of records 'str' is not an integer value",
	},
	{
		Name: "Limit Negative Value",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Limit: parser.LimitClause{Value: parser.NewIntegerValue(-1)},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{},
			Filter:    NewEmptyFilter(),
		},
	},
	{
		Name: "Limit By Percentage Value Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Limit: parser.LimitClause{Value: parser.NewStringValue("str"), Percent: "percent"},
		Error: "[L:- C:-] limit percentage 'str' is not a float value",
	},
}

func TestView_Limit(t *testing.T) {
	for _, v := range viewLimitTests {
		err := v.View.Limit(v.Limit)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if !reflect.DeepEqual(v.View, v.Result) {
			t.Errorf("%s: view = %v, want %v", v.Name, v.View, v.Result)
		}
	}
}

var viewOffsetTests = []struct {
	Name   string
	View   *View
	Offset parser.OffsetClause
	Result *View
	Error  string
}{
	{
		Name: "Offset",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Offset: parser.OffsetClause{Value: parser.NewIntegerValueFromString("3")},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			Filter: NewEmptyFilter(),
			offset: 3,
		},
	},
	{
		Name: "Offset Equal To Record Length",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Offset: parser.OffsetClause{Value: parser.NewIntegerValueFromString("4")},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{},
			Filter:    NewEmptyFilter(),
			offset:    4,
		},
	},
	{
		Name: "Offset Evaluate Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Offset: parser.OffsetClause{Value: parser.Variable{Name: "notexist"}},
		Error:  "[L:- C:-] variable notexist is undeclared",
	},
	{
		Name: "Offset Value Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Offset: parser.OffsetClause{Value: parser.NewStringValue("str")},
		Error:  "[L:- C:-] offset number 'str' is not an integer value",
	},
	{
		Name: "Offset Negative Number",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			Filter: NewEmptyFilter(),
		},
		Offset: parser.OffsetClause{Value: parser.NewIntegerValue(-3)},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: INTERNAL_ID_COLUMN},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			Filter: NewEmptyFilter(),
			offset: 0,
		},
	},
}

func TestView_Offset(t *testing.T) {
	for _, v := range viewOffsetTests {
		err := v.View.Offset(v.Offset)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if !reflect.DeepEqual(v.View, v.Result) {
			t.Errorf("%s: view = %v, want %v", v.Name, v.View, v.Result)
		}
	}
}

var viewInsertValuesTests = []struct {
	Name       string
	Fields     []parser.QueryExpression
	ValuesList []parser.QueryExpression
	Result     *View
	Error      string
}{
	{
		Name: "InsertValues",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("4"),
					},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewNull(),
					value.NewInteger(3),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewNull(),
					value.NewInteger(4),
					value.NewNull(),
				}),
			},
			Filter:          NewEmptyFilter(),
			OperatedRecords: 2,
		},
	},
	{
		Name: "InsertValues Field Length Does Not Match Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "[L:- C:-] row value should contain exactly 2 values",
	},
	{
		Name: "InsertValues Value Evaluation Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "InsertValues Field Does Not Exist Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestView_InsertValues(t *testing.T) {
	view := &View{
		Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
		RecordSet: []Record{
			NewRecordWithId(1, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecordWithId(2, []value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
		Filter: NewEmptyFilter(),
	}

	for _, v := range viewInsertValuesTests {
		err := view.InsertValues(v.Fields, v.ValuesList)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if !reflect.DeepEqual(view, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, view, v.Result)
		}
	}
}

var viewInsertFromQueryTests = []struct {
	Name   string
	Fields []parser.QueryExpression
	Query  parser.SelectQuery
	Result *View
	Error  string
}{
	{
		Name: "InsertFromQuery",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.NewIntegerValueFromString("3")},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewNull(),
					value.NewInteger(3),
					value.NewNull(),
				}),
			},
			Filter:          NewEmptyFilter(),
			OperatedRecords: 1,
		},
	},
	{
		Name: "InsertFromQuery Field Lenght Does Not Match Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.NewIntegerValueFromString("3")},
					},
				},
			},
		},
		Error: "[L:- C:-] select query should return exactly 2 fields",
	},
	{
		Name: "Insert Values Query Exuecution Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestView_InsertFromQuery(t *testing.T) {
	view := &View{
		Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
		RecordSet: []Record{
			NewRecordWithId(1, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecordWithId(2, []value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
		Filter: NewEmptyFilter(),
	}

	for _, v := range viewInsertFromQueryTests {
		err := view.InsertFromQuery(v.Fields, v.Query)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if !reflect.DeepEqual(view, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, view, v.Result)
		}
	}
}

func TestView_Fix(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: INTERNAL_ID_COLUMN},
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: []Record{
			NewRecordWithId(1, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecordWithId(2, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
		},
		selectFields: []int{2},
	}
	expect := &View{
		Header: NewHeader("table1", []string{"column2"}),
		RecordSet: []Record{
			NewRecord([]value.Primary{
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("str1"),
			}),
		},
		selectFields: []int(nil),
	}

	view.Fix()
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("fix: view = %v, want %v", view, expect)
	}
}

func TestView_Union(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	calcView := &View{
		Header: []HeaderField{
			{View: "table2", Column: "column3", IsFromTable: true},
			{View: "table2", Column: "column4", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("3"),
				value.NewString("str3"),
			}),
		},
	}

	expect := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("3"),
				value.NewString("str3"),
			}),
		},
	}

	view.Union(calcView, false)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("union: view = %v, want %v", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	expect = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("3"),
				value.NewString("str3"),
			}),
		},
	}

	view.Union(calcView, true)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("union all: view = %v, want %v", view, expect)
	}
}

func TestView_Except(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	calcView := &View{
		Header: []HeaderField{
			{View: "table2", Column: "column3", IsFromTable: true},
			{View: "table2", Column: "column4", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("3"),
				value.NewString("str3"),
			}),
		},
	}

	expect := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
		},
	}

	view.Except(calcView, false)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("except: view = %v, want %v", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	expect = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
		},
	}

	view.Except(calcView, true)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("except all: view = %v, want %v", view, expect)
	}
}

func TestView_Intersect(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	calcView := &View{
		Header: []HeaderField{
			{View: "table2", Column: "column3", IsFromTable: true},
			{View: "table2", Column: "column4", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("3"),
				value.NewString("str3"),
			}),
		},
	}

	expect := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	view.Intersect(calcView, false)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("intersect: view = %v, want %v", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	expect = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	view.Intersect(calcView, true)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("intersect all: view = %v, want %v", view, expect)
	}
}

func TestView_FieldIndex(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
			{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
		},
	}
	fieldRef := parser.FieldReference{
		Column: parser.Identifier{Literal: "column1"},
	}
	expect := 0

	idx, _ := view.FieldIndex(fieldRef)
	if idx != expect {
		t.Errorf("field index = %d, want %d", idx, expect)
	}

	columnNum := parser.ColumnNumber{
		View:   parser.Identifier{Literal: "table1"},
		Number: value.NewInteger(2),
	}
	expect = 1

	idx, _ = view.FieldIndex(columnNum)
	if idx != expect {
		t.Errorf("field index = %d, want %d", idx, expect)
	}
}

func TestView_FieldIndices(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
	}
	fields := []parser.QueryExpression{
		parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
	}
	expect := []int{1, 0}

	indices, _ := view.FieldIndices(fields)
	if !reflect.DeepEqual(indices, expect) {
		t.Errorf("field indices = %v, want %v", indices, expect)
	}

	fields = []parser.QueryExpression{
		parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
	}
	expectErr := "[L:- C:-] field notexist does not exist"
	_, err := view.FieldIndices(fields)
	if err.Error() != expectErr {
		t.Errorf("error = %s, want %s", err, expectErr)
	}
}

func TestView_FieldViewName(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table2", Column: "column2", IsFromTable: true},
		},
	}
	fieldRef := parser.FieldReference{
		Column: parser.Identifier{Literal: "column1"},
	}
	expect := "table1"

	ref, _ := view.FieldViewName(fieldRef)
	if ref != expect {
		t.Errorf("field reference = %s, want %s", ref, expect)
	}

	fieldRef = parser.FieldReference{
		Column: parser.Identifier{Literal: "notexist"},
	}
	expectErr := "[L:- C:-] field notexist does not exist"
	_, err := view.FieldViewName(fieldRef)
	if err.Error() != expectErr {
		t.Errorf("error = %s, want %s", err, expectErr)
	}
}

func TestView_InternalRecordId(t *testing.T) {
	view := &View{
		Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
		RecordSet: []Record{
			NewRecordWithId(0, []value.Primary{value.NewInteger(1), value.NewString("str1")}),
			NewRecordWithId(1, []value.Primary{value.NewInteger(2), value.NewString("str2")}),
			NewRecordWithId(2, []value.Primary{value.NewInteger(3), value.NewString("str3")}),
		},
	}
	ref := "table1"
	recordIndex := 1
	expect := 1

	id, _ := view.InternalRecordId(ref, recordIndex)
	if id != expect {
		t.Errorf("field internal id = %d, want %d", id, expect)
	}

	view.RecordSet[1][0] = NewCell(value.NewNull())
	expectError := "[L:- C:-] internal record id is empty"
	_, err := view.InternalRecordId(ref, recordIndex)
	if err.Error() != expectError {
		t.Errorf("error = %q, want error %q", err, expectError)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table2", Column: "column2", IsFromTable: true},
		},
	}
	expectError = "[L:- C:-] internal record id does not exist"
	_, err = view.InternalRecordId(ref, recordIndex)
	if err.Error() != expectError {
		t.Errorf("error = %q, want error %q", err, expectError)
	}
}
