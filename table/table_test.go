package table

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testCaption         = "A Song of Ice and Fire"
	testColor           = text.Colors{text.FgGreen}
	testColorHiRedBold  = text.Colors{text.FgHiRed, text.Bold}
	testColorHiBlueBold = text.Colors{text.FgHiBlue, text.Bold}
	testCSSClass        = "test-css-class"
	testFooter          = Row{"", "", "Total", 10000}
	testFooterMultiLine = Row{"", "", "Total\nSalary", 10000}
	testHeader          = Row{"#", "First Name", "Last Name", "Salary"}
	testHeaderMultiLine = Row{"#", "First\nName", "Last\nName", "Salary"}
	testHeader2         = Row{"ID", "Text1", "Date", "Text2"}
	testRows            = []Row{
		{1, "Arya", "Stark", 3000},
		{20, "Jon", "Snow", 2000, "You know nothing, Jon Snow!"},
		{300, "Tyrion", "Lannister", 5000},
	}
	testRowMultiLine    = Row{0, "Winter", "Is", 0, "Coming.\r\nThe North Remembers!\nThis is known."}
	testRowNewLines     = Row{0, "Valar", "Morghulis", 0, "Faceless\nMen"}
	testRowPipes        = Row{0, "Valar", "Morghulis", 0, "Faceless|Men"}
	testRowTabs         = Row{0, "Valar", "Morghulis", 0, "Faceless\tMen"}
	testRowDoubleQuotes = Row{0, "Valar", "Morghulis", 0, "Faceless\"Men"}
	testTitle1          = "Game of Thrones"
	testTitle2          = "When you play the Game of Thrones, you win or you die. There is no middle ground."
)

func init() {
	text.EnableColors()
}

type myMockOutputMirror struct {
	mirroredOutput string
}

func (t *myMockOutputMirror) Write(p []byte) (n int, err error) {
	t.mirroredOutput += string(p)
	return len(p), nil
}

func TestNewWriter(t *testing.T) {
	tw := NewWriter()
	assert.NotNil(t, tw.Style())
	assert.Equal(t, StyleDefault, *tw.Style())

	tw.SetStyle(StyleBold)
	assert.NotNil(t, tw.Style())
	assert.Equal(t, StyleBold, *tw.Style())
}

func TestTable_AppendFooter(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, len(table.rowsFooterRaw))

	table.AppendFooter([]interface{}{})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 1, len(table.rowsFooterRaw))
	assert.Equal(t, 0, len(table.rowsHeaderRaw))

	table.AppendFooter([]interface{}{})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 2, len(table.rowsFooterRaw))
	assert.Equal(t, 0, len(table.rowsHeaderRaw))

	table.AppendFooter([]interface{}{}, RowConfig{AutoMerge: true})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 3, len(table.rowsFooterRaw))
	assert.Equal(t, 0, len(table.rowsHeaderRaw))
	assert.False(t, table.rowsFooterConfigMap[0].AutoMerge)
	assert.False(t, table.rowsFooterConfigMap[1].AutoMerge)
	assert.True(t, table.rowsFooterConfigMap[2].AutoMerge)
}

func TestTable_AppendHeader(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, len(table.rowsHeaderRaw))

	table.AppendHeader([]interface{}{})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 0, len(table.rowsFooterRaw))
	assert.Equal(t, 1, len(table.rowsHeaderRaw))

	table.AppendHeader([]interface{}{})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 0, len(table.rowsFooterRaw))
	assert.Equal(t, 2, len(table.rowsHeaderRaw))

	table.AppendHeader([]interface{}{}, RowConfig{AutoMerge: true})
	assert.Equal(t, 0, table.Length())
	assert.Equal(t, 0, len(table.rowsFooterRaw))
	assert.Equal(t, 3, len(table.rowsHeaderRaw))
	assert.False(t, table.rowsHeaderConfigMap[0].AutoMerge)
	assert.False(t, table.rowsHeaderConfigMap[1].AutoMerge)
	assert.True(t, table.rowsHeaderConfigMap[2].AutoMerge)
}

func TestTable_AppendRow(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, table.Length())

	table.AppendRow([]interface{}{})
	assert.Equal(t, 1, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))

	table.AppendRow([]interface{}{})
	assert.Equal(t, 2, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))

	table.AppendRow([]interface{}{}, RowConfig{AutoMerge: true})
	assert.Equal(t, 3, table.Length())
	assert.Equal(t, 0, len(table.rowsFooterRaw))
	assert.Equal(t, 0, len(table.rowsHeaderRaw))
	assert.False(t, table.rowsConfigMap[0].AutoMerge)
	assert.False(t, table.rowsConfigMap[1].AutoMerge)
	assert.True(t, table.rowsConfigMap[2].AutoMerge)
}

func TestTable_AppendRows(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, table.Length())

	table.AppendRows([]Row{{}})
	assert.Equal(t, 1, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))

	table.AppendRows([]Row{{}})
	assert.Equal(t, 2, table.Length())
	assert.Equal(t, 0, len(table.rowsFooter))
	assert.Equal(t, 0, len(table.rowsHeader))

	table.AppendRows([]Row{{}, {}}, RowConfig{AutoMerge: true})
	assert.Equal(t, 4, table.Length())
	assert.Equal(t, 0, len(table.rowsFooterRaw))
	assert.Equal(t, 0, len(table.rowsHeaderRaw))
	assert.False(t, table.rowsConfigMap[0].AutoMerge)
	assert.False(t, table.rowsConfigMap[1].AutoMerge)
	assert.True(t, table.rowsConfigMap[2].AutoMerge)
	assert.True(t, table.rowsConfigMap[3].AutoMerge)
}

func TestTable_ImportGrid(t *testing.T) {
	t.Run("invalid grid", func(t *testing.T) {
		table := Table{}

		assert.False(t, table.ImportGrid(nil))
		require.Len(t, table.rowsRawFiltered, 0)

		assert.False(t, table.ImportGrid(123))
		require.Len(t, table.rowsRawFiltered, 0)

		assert.False(t, table.ImportGrid("abc"))
		require.Len(t, table.rowsRawFiltered, 0)

		assert.False(t, table.ImportGrid(Table{}))
		require.Len(t, table.rowsRawFiltered, 0)

		assert.False(t, table.ImportGrid(&Table{}))
		require.Len(t, table.rowsRawFiltered, 0)
	})

	a, b, c := 1, 2, 3
	d, e, f := 4, 5, 6
	g, h, i := 7, 8, 9

	t.Run("valid 1d", func(t *testing.T) {
		inputs := []interface{}{
			[3]int{a, b, c},     // array
			[]int{a, b, c},      // slice
			&[]int{a, b, c},     // pointer to slice
			[]*int{&a, &b, &c},  // slice of pointers-to-slices
			&[]*int{&a, &b, &c}, // pointer to slice of pointers
		}

		for _, grid := range inputs {
			message := fmt.Sprintf("grid: %#v", grid)

			table := Table{}
			table.Style().Options.SeparateRows = true
			assert.True(t, table.ImportGrid(grid), message)
			compareOutput(t, table.Render(), `
+---+
| 1 |
+---+
| 2 |
+---+
| 3 |
+---+`, message)
		}
	})

	t.Run("valid 2d", func(t *testing.T) {
		inputs := []interface{}{
			[3][3]int{{a, b, c}, {d, e, f}, {g, h, i}},           // array of arrays
			[3][]int{{a, b, c}, {d, e, f}, {g, h, i}},            // array of slices
			[][]int{{a, b, c}, {d, e, f}, {g, h, i}},             // slice of slices
			&[][]int{{a, b, c}, {d, e, f}, {g, h, i}},            // pointer-to-slice of slices
			[]*[]int{{a, b, c}, {d, e, f}, {g, h, i}},            // slice of pointers-to-slices
			&[]*[]int{{a, b, c}, {d, e, f}, {g, h, i}},           // pointer-to-slice of pointers-to-slices
			&[]*[]*int{{&a, &b, &c}, {&d, &e, &f}, {&g, &h, &i}}, // pointer-to-slice of pointers-to-slices of pointers
		}

		for _, grid := range inputs {
			message := fmt.Sprintf("grid: %#v", grid)

			table := Table{}
			table.Style().Options.SeparateRows = true
			assert.True(t, table.ImportGrid(grid), message)
			compareOutput(t, table.Render(), `
+---+---+---+
| 1 | 2 | 3 |
+---+---+---+
| 4 | 5 | 6 |
+---+---+---+
| 7 | 8 | 9 |
+---+---+---+`, message)
		}
	})

	t.Run("valid 2d with nil rows", func(t *testing.T) {
		inputs := []interface{}{
			[]*[]int{{a, b, c}, {d, e, f}, nil},         // slice of pointers-to-slices
			&[]*[]int{{a, b, c}, {d, e, f}, nil},        // pointer-to-slice of pointers-to-slices
			&[]*[]*int{{&a, &b, &c}, {&d, &e, &f}, nil}, // pointer-to-slice of pointers-to-slices of pointers
		}

		for _, grid := range inputs {
			message := fmt.Sprintf("grid: %#v", grid)

			table := Table{}
			table.Style().Options.SeparateRows = true
			assert.True(t, table.ImportGrid(grid), message)
			compareOutput(t, table.Render(), `
+---+---+---+
| 1 | 2 | 3 |
+---+---+---+
| 4 | 5 | 6 |
+---+---+---+`, message)
		}
	})

	t.Run("valid 2d with nil columns and rows", func(t *testing.T) {
		inputs := []interface{}{
			&[]*[]*int{{&a, &b, &c}, {&d, &e, nil}, nil},
		}

		for _, grid := range inputs {
			message := fmt.Sprintf("grid: %#v", grid)

			table := Table{}
			table.Style().Options.SeparateRows = true
			assert.True(t, table.ImportGrid(grid), message)
			compareOutput(t, table.Render(), `
+---+---+---+
| 1 | 2 | 3 |
+---+---+---+
| 4 | 5 |   |
+---+---+---+`, message)
		}
	})
}

func TestTable_Length(t *testing.T) {
	table := Table{}
	assert.Zero(t, table.Length())

	table.AppendRow(testRows[0])
	assert.Equal(t, 1, table.Length())
	table.AppendRow(testRows[1])
	assert.Equal(t, 2, table.Length())

	table.AppendHeader(testHeader)
	assert.Equal(t, 2, table.Length())
}

func TestTable_ResetFooters(t *testing.T) {
	table := Table{}
	table.AppendFooter(testFooter)
	assert.NotEmpty(t, table.rowsFooterRaw)

	table.ResetFooters()
	assert.Empty(t, table.rowsFooterRaw)
}

func TestTable_ResetHeaders(t *testing.T) {
	table := Table{}
	table.AppendHeader(testHeader)
	assert.NotEmpty(t, table.rowsHeaderRaw)

	table.ResetHeaders()
	assert.Empty(t, table.rowsHeaderRaw)
}

func TestTable_ResetRows(t *testing.T) {
	table := Table{}
	table.AppendRows(testRows)
	assert.NotEmpty(t, table.rowsRawFiltered)

	table.ResetRows()
	assert.Empty(t, table.rowsRawFiltered)
}

func TestTable_SetAllowedRowLength(t *testing.T) {
	table := Table{}
	table.AppendRows(testRows)
	table.SetStyle(styleTest)

	expectedOutWithNoRowLimit := `(-----^--------^-----------^------^-----------------------------)
[<  1>|<Arya  >|<Stark    >|<3000>|<                           >]
[< 20>|<Jon   >|<Snow     >|<2000>|<You know nothing, Jon Snow!>]
[<300>|<Tyrion>|<Lannister>|<5000>|<                           >]
\-----v--------v-----------v------v-----------------------------/`
	assert.Zero(t, table.allowedRowLength)
	assert.Equal(t, expectedOutWithNoRowLimit, table.Render())

	table.SetAllowedRowLength(utf8.RuneCountInString(table.style.Box.UnfinishedRow))
	assert.Equal(t, utf8.RuneCountInString(table.style.Box.UnfinishedRow), table.allowedRowLength)
	assert.Equal(t, "", table.Render())

	table.SetAllowedRowLength(5)
	expectedOutWithRowLimit := `( ~~~
[ ~~~
[ ~~~
[ ~~~
\ ~~~`
	assert.Equal(t, 5, table.allowedRowLength)
	assert.Equal(t, expectedOutWithRowLimit, table.Render())

	table.SetAllowedRowLength(30)
	expectedOutWithRowLimit = `(-----^--------^---------- ~~~
[<  1>|<Arya  >|<Stark     ~~~
[< 20>|<Jon   >|<Snow      ~~~
[<300>|<Tyrion>|<Lannister ~~~
\-----v--------v---------- ~~~`
	assert.Equal(t, 30, table.allowedRowLength)
	assert.Equal(t, expectedOutWithRowLimit, table.Render())

	table.SetAllowedRowLength(300)
	assert.Equal(t, 300, table.allowedRowLength)
	assert.Equal(t, expectedOutWithNoRowLimit, table.Render())
}

func TestTable_SetAutoIndex(t *testing.T) {
	table := Table{}
	table.AppendRows(testRows)
	table.SetStyle(styleTest)

	expectedOut := `(-----^--------^-----------^------^-----------------------------)
[<  1>|<Arya  >|<Stark    >|<3000>|<                           >]
[< 20>|<Jon   >|<Snow     >|<2000>|<You know nothing, Jon Snow!>]
[<300>|<Tyrion>|<Lannister>|<5000>|<                           >]
\-----v--------v-----------v------v-----------------------------/`
	assert.False(t, table.autoIndex)
	assert.Equal(t, expectedOut, table.Render())

	table.SetAutoIndex(true)
	expectedOut = `(---^-----^--------^-----------^------^-----------------------------)
[< >|<  A>|<   B  >|<    C    >|<   D>|<             E             >]
{---+-----+--------+-----------+------+-----------------------------}
[<1>|<  1>|<Arya  >|<Stark    >|<3000>|<                           >]
[<2>|< 20>|<Jon   >|<Snow     >|<2000>|<You know nothing, Jon Snow!>]
[<3>|<300>|<Tyrion>|<Lannister>|<5000>|<                           >]
\---v-----v--------v-----------v------v-----------------------------/`
	assert.True(t, table.autoIndex)
	assert.Equal(t, expectedOut, table.Render())

	table.AppendHeader(testHeader)
	expectedOut = `(---^-----^------------^-----------^--------^-----------------------------)
[< >|<  #>|<FIRST NAME>|<LAST NAME>|<SALARY>|<                           >]
{---+-----+------------+-----------+--------+-----------------------------}
[<1>|<  1>|<Arya      >|<Stark    >|<  3000>|<                           >]
[<2>|< 20>|<Jon       >|<Snow     >|<  2000>|<You know nothing, Jon Snow!>]
[<3>|<300>|<Tyrion    >|<Lannister>|<  5000>|<                           >]
\---v-----v------------v-----------v--------v-----------------------------/`
	assert.True(t, table.autoIndex)
	assert.Equal(t, expectedOut, table.Render())

	table.AppendRow(testRowMultiLine)
	expectedOut = `(---^-----^------------^-----------^--------^-----------------------------)
[< >|<  #>|<FIRST NAME>|<LAST NAME>|<SALARY>|<                           >]
{---+-----+------------+-----------+--------+-----------------------------}
[<1>|<  1>|<Arya      >|<Stark    >|<  3000>|<                           >]
[<2>|< 20>|<Jon       >|<Snow     >|<  2000>|<You know nothing, Jon Snow!>]
[<3>|<300>|<Tyrion    >|<Lannister>|<  5000>|<                           >]
[<4>|<  0>|<Winter    >|<Is       >|<     0>|<Coming.                    >]
[< >|<   >|<          >|<         >|<      >|<The North Remembers!       >]
[< >|<   >|<          >|<         >|<      >|<This is known.             >]
\---v-----v------------v-----------v--------v-----------------------------/`
	assert.Equal(t, expectedOut, table.Render())

	table.SetStyle(StyleLight)
	expectedOut = `┌───┬─────┬────────────┬───────────┬────────┬─────────────────────────────┐
│   │   # │ FIRST NAME │ LAST NAME │ SALARY │                             │
├───┼─────┼────────────┼───────────┼────────┼─────────────────────────────┤
│ 1 │   1 │ Arya       │ Stark     │   3000 │                             │
│ 2 │  20 │ Jon        │ Snow      │   2000 │ You know nothing, Jon Snow! │
│ 3 │ 300 │ Tyrion     │ Lannister │   5000 │                             │
│ 4 │   0 │ Winter     │ Is        │      0 │ Coming.                     │
│   │     │            │           │        │ The North Remembers!        │
│   │     │            │           │        │ This is known.              │
└───┴─────┴────────────┴───────────┴────────┴─────────────────────────────┘`
	assert.Equal(t, expectedOut, table.Render())
}

func TestTable_SetCaption(t *testing.T) {
	table := Table{}
	assert.Empty(t, table.caption)

	table.SetCaption(testCaption)
	assert.NotEmpty(t, table.caption)
	assert.Equal(t, testCaption, table.caption)
}

func TestTable_SetColumnConfigs(t *testing.T) {
	table := Table{}
	assert.Empty(t, table.columnConfigs)

	table.SetColumnConfigs([]ColumnConfig{{}, {}, {}})
	assert.NotEmpty(t, table.columnConfigs)
	assert.Equal(t, 3, len(table.columnConfigs))
}

func TestTable_SetHTMLCSSClass(t *testing.T) {
	table := Table{}
	table.AppendRow(testRows[0])
	expectedHTML := `<table class="` + DefaultHTMLCSSClass + `">
  <tbody>
  <tr>
    <td align="right">1</td>
    <td>Arya</td>
    <td>Stark</td>
    <td align="right">3000</td>
  </tr>
  </tbody>
</table>`
	assert.Equal(t, "", table.htmlCSSClass)
	assert.Equal(t, expectedHTML, table.RenderHTML())

	table.SetHTMLCSSClass(testCSSClass)
	assert.Equal(t, testCSSClass, table.htmlCSSClass)
	assert.Equal(t, strings.Replace(expectedHTML, DefaultHTMLCSSClass, testCSSClass, -1), table.RenderHTML())
}

func TestTable_SetOutputMirror(t *testing.T) {
	table := Table{}
	table.AppendRow(testRows[0])
	expectedOut := `+---+------+-------+------+
| 1 | Arya | Stark | 3000 |
+---+------+-------+------+`
	assert.Equal(t, nil, table.outputMirror)
	assert.Equal(t, expectedOut, table.Render())

	mockOutputMirror := &myMockOutputMirror{}
	table.SetOutputMirror(mockOutputMirror)
	assert.Equal(t, mockOutputMirror, table.outputMirror)
	assert.Equal(t, expectedOut, table.Render())
	assert.Equal(t, expectedOut+"\n", mockOutputMirror.mirroredOutput)
}

func TestTable_SePageSize(t *testing.T) {
	table := Table{}
	assert.Equal(t, 0, table.pager.size)

	table.SetPageSize(13)
	assert.Equal(t, 13, table.pager.size)
}

func TestTable_SortByColumn(t *testing.T) {
	table := Table{}
	assert.Empty(t, table.sortBy)

	table.SortBy([]SortBy{{Name: "#", Mode: Asc}})
	assert.Equal(t, 1, len(table.sortBy))

	table.SortBy([]SortBy{{Name: "First Name", Mode: Dsc}, {Name: "Last Name", Mode: Asc}})
	assert.Equal(t, 2, len(table.sortBy))
}

func TestTable_SetStyle(t *testing.T) {
	table := Table{}
	assert.NotNil(t, table.Style())
	assert.Equal(t, StyleDefault, *table.Style())

	table.SetStyle(StyleDefault)
	assert.NotNil(t, table.Style())
	assert.Equal(t, StyleDefault, *table.Style())
}

func TestTable_ColumsHorizontalMerge(t *testing.T) {
	expectedOutput := `
╭───┬──────────────────────────────────────┬────────────────────┬─────────────┬─────────────────────────┬──────────┬────────┬───────────┬─────────╮
│ # │ Run ID                               │ Plugin Name        │ Started By  │ Start Time              │ Duration │ Done   │ Cancelled │ Success │
├───┼──────────────────────────────────────┼────────────────────┼─────────────┼─────────────────────────┼──────────┼────────┼───────────┼─────────┤
│ 1 │ 446056cf-8737-459b-8295-2f0fc34b1f30 │ clients_benchmarks │ john_dow_12 │ 2025-12-05 08:51:52 +08 │ 2s       │ true   │ false     │ true    │
│ 2 │ 12140ce3-74ce-474e-b618-572649b6be47 │ clients_benchmarks │ john_dow_12 │ 2025-12-05 04:53:50 +08 │ 42s      │ true   │ false     │ false   │
├───┼──────────────────────────────────────┴────────────────────┴─────────────┴─────────────────────────┴──────────┴────────┴───────────┴─────────┤
│ > │ Failed with error: Run(): File upload failed. The file size exceeds the maximum limit of 5MB for this operation. Please compress your data, │
│   │ remove unnecessary content, or split it into multiple smaller files, then try uploading again with a file that is within the allowed size   │
│   │ limit. If the problem persists, check your connection or contact support.                                                                   │
╰───┴─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯`

	// build vertical table using go-pretty
	tw := NewWriter()
	tw.SetStyle(StyleRounded)
	tw.SetOutputMirror(os.Stdout) // ensure the table is printed to stdout

	minColWidth := 1
	maxColWidth := 200
	maxWrap := 139

	tw.SetColumnConfigs([]ColumnConfig{
		{
			Number:   2,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   3,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   4,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   5,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   6,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   7,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   8,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   9,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
	})

	rowConfig := RowConfig{AutoMerge: false, AutoMergeAlign: text.AlignLeft}

	pluginNameStyle := text.Colors{}
	normalStyle := text.Colors{}
	successStyle := text.Colors{}
	failureStyle := text.Colors{}
	inProgressStyle := text.Colors{}

	isDone := true
	isDoneStr := ""
	if isDone {
		isDoneStr = successStyle.Sprint(isDone)
	} else {
		isDoneStr = inProgressStyle.Sprint(isDone)
	}

	isCancelled := false
	isCancelledStr := ""
	if isDone {
		if isCancelled {
			isCancelledStr = inProgressStyle.Sprint(isCancelled)
		} else {
			isCancelledStr = successStyle.Sprint(isCancelled)
		}
	} else {
		isCancelledStr = inProgressStyle.Sprint(isCancelled)
	}

	isSuccess := true
	isSuccessStr := ""
	if isDone {
		if isSuccess {
			isSuccessStr = successStyle.Sprint(isSuccess)
		} else {
			isSuccessStr = failureStyle.Sprint(isSuccess)
		}
	} else {
		isSuccessStr = inProgressStyle.Sprint(isSuccess)
	}

	tw.AppendRow(Row{"#", "Run ID", "Plugin Name", "Started By", "Start Time", "Duration", "Done", "Cancelled", "Success"})
	tw.AppendSeparator()

	// row 1 - success
	row := Row{
		1,
		"446056cf-8737-459b-8295-2f0fc34b1f30",
		pluginNameStyle.Sprint("clients_benchmarks"),
		"john_dow_12",
		"2025-12-05 08:51:52 +08",
		pluginNameStyle.Sprint("2s"),
		isDoneStr,
		isCancelledStr,
		isSuccessStr,
	}

	tw.AppendRow(row, rowConfig)

	isSuccess = false
	isSuccessStr = ""
	if isDone {
		if isSuccess {
			isSuccessStr = successStyle.Sprint(isSuccess)
		} else {
			isSuccessStr = failureStyle.Sprint(isSuccess)
		}
	} else {
		isSuccessStr = inProgressStyle.Sprint(isSuccess)
	}

	// row 2 - failure
	row = Row{
		2,
		"12140ce3-74ce-474e-b618-572649b6be47",
		pluginNameStyle.Sprint("clients_benchmarks"),
		"john_dow_12",
		"2025-12-05 04:53:50 +08",
		pluginNameStyle.Sprint("42s"),
		isDoneStr,
		isCancelledStr,
		isSuccessStr,
	}

	tw.AppendRow(row, rowConfig)

	// row 3 - failure message

	errorMsg := "Run(): File upload failed. The file size exceeds the maximum limit of 5MB for this operation. Please compress your data, remove unnecessary content, or split it into multiple smaller files, then try uploading again with a file that is within the allowed size limit. If the problem persists, check your connection or contact support. "
	ff := failureStyle.Sprint("Failed with error: ") + normalStyle.Sprint(errorMsg)
	arr := failureStyle.Sprint(">")
	tw.AppendSeparator()
	tw.AppendRow(Row{arr, ff, ff, ff, ff, ff, ff, ff, ff}, RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignLeft})
	tw.AppendSeparator()

	compareOutput(t, tw.Render(), expectedOutput)
}

func TestTable_ColumsHorizontalMerge2(t *testing.T) {
	expectedOutput := `
╭───┬──────────────────────────────────────┬────────────────────┬─────────────┬─────────────────────────┬──────────┬────────┬───────────┬─────────╮
│ # │ Run ID                               │ Plugin Name        │ Started By  │ Start Time              │ Duration │ Done   │ Cancelled │ Success │
├───┼──────────────────────────────────────┼────────────────────┼─────────────┼─────────────────────────┼──────────┼────────┼───────────┼─────────┤
│ 1 │ 446056cf-8737-459b-8295-2f0fc34b1f30 │ clients_benchmarks │ john_dow_12 │ 2025-12-05 08:51:52 +08 │ 2s       │ true   │ false     │ true    │
│ 2 │ 12140ce3-74ce-474e-b618-572649b6be47 │ clients_benchmarks │ john_dow_12 │ 2025-12-05 04:53:50 +08 │ 42s      │ true   │ false     │ false   │
├───┼──────────────────────────────────────┴────────────────────┴─────────────┴─────────────────────────┴──────────┴────────┴───────────┴─────────┤
│ > │ Failed with error: Run(): File upload failed. The file size exceeds the maximum limit of 5MB for this operation. Please compress your data, │
│   │ remove unnecessary content, or split it into multiple smaller files, then try uploading again with a file that is within the allowed size   │
│   │ limit. If the problem persists, check your connection or contact support.                                                                   │
╰───┴─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯`

	// build vertical table using go-pretty
	tw := NewWriter()
	tw.SetStyle(StyleRounded)
	tw.SetOutputMirror(os.Stdout) // ensure the table is printed to stdout

	minColWidth := 1
	maxColWidth := 200
	maxWrap := 139

	tw.SetColumnConfigs([]ColumnConfig{
		{
			Number:   2,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   3,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   4,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   5,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   6,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   7,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   8,
			WidthMin: minColWidth,
			WidthMax: maxColWidth,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
		{
			Number:   9,
			WidthMin: minColWidth,
			WidthMax: 9,
			Transformer: text.Transformer(func(val interface{}) string {
				return text.WrapSoft(fmt.Sprintf("%v", val), maxWrap)
			}),
		},
	})

	rowConfig := RowConfig{AutoMerge: false, AutoMergeAlign: text.AlignLeft}

	pluginNameStyle := text.Colors{}
	normalStyle := text.Colors{}
	successStyle := text.Colors{}
	failureStyle := text.Colors{}
	inProgressStyle := text.Colors{}

	isDone := true
	isDoneStr := ""
	if isDone {
		isDoneStr = successStyle.Sprint(isDone)
	} else {
		isDoneStr = inProgressStyle.Sprint(isDone)
	}

	isCancelled := false
	isCancelledStr := ""
	if isDone {
		if isCancelled {
			isCancelledStr = inProgressStyle.Sprint(isCancelled)
		} else {
			isCancelledStr = successStyle.Sprint(isCancelled)
		}
	} else {
		isCancelledStr = inProgressStyle.Sprint(isCancelled)
	}

	isSuccess := true
	isSuccessStr := ""
	if isDone {
		if isSuccess {
			isSuccessStr = successStyle.Sprint(isSuccess)
		} else {
			isSuccessStr = failureStyle.Sprint(isSuccess)
		}
	} else {
		isSuccessStr = inProgressStyle.Sprint(isSuccess)
	}

	tw.AppendRow(Row{"#", "Run ID", "Plugin Name", "Started By", "Start Time", "Duration", "Done", "Cancelled", "Success"})
	tw.AppendSeparator()

	// row 1 - success
	row := Row{
		1,
		"446056cf-8737-459b-8295-2f0fc34b1f30",
		pluginNameStyle.Sprint("clients_benchmarks"),
		"john_dow_12",
		"2025-12-05 08:51:52 +08",
		pluginNameStyle.Sprint("2s"),
		isDoneStr,
		isCancelledStr,
		isSuccessStr,
	}

	tw.AppendRow(row, rowConfig)

	isSuccess = false
	isSuccessStr = ""
	if isDone {
		if isSuccess {
			isSuccessStr = successStyle.Sprint(isSuccess)
		} else {
			isSuccessStr = failureStyle.Sprint(isSuccess)
		}
	} else {
		isSuccessStr = inProgressStyle.Sprint(isSuccess)
	}

	// row 2 - failure
	row = Row{
		2,
		"12140ce3-74ce-474e-b618-572649b6be47",
		pluginNameStyle.Sprint("clients_benchmarks"),
		"john_dow_12",
		"2025-12-05 04:53:50 +08",
		pluginNameStyle.Sprint("42s"),
		isDoneStr,
		isCancelledStr,
		isSuccessStr,
	}

	tw.AppendRow(row, rowConfig)

	// row 3 - failure message

	errorMsg := "Run(): File upload failed. The file size exceeds the maximum limit of 5MB for this operation. Please compress your data, remove unnecessary content, or split it into multiple smaller files, then try uploading again with a file that is within the allowed size limit. If the problem persists, check your connection or contact support. "
	ff := failureStyle.Sprint("Failed with error: ") + normalStyle.Sprint(errorMsg)
	arr := failureStyle.Sprint(">")
	tw.AppendSeparator()
	tw.AppendRow(Row{arr, ff, ff, ff, ff, ff, ff, ff, ff}, RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignLeft})
	tw.AppendSeparator()

	compareOutput(t, tw.Render(), expectedOutput)
}
