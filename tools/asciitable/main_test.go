package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseJSON(t *testing.T) {
	jsonData := []byte(`[{"conclusion": ":white_check_mark:", "cap": "Cap1", "html_url": "http://example.com"}]`)
	jobs, err := parseJSON(jsonData, "")
	require.NoError(t, err)
	require.Len(t, jobs, 1)
	require.Equal(t, ":white_check_mark:", jobs[0].Conclusion)
	require.Equal(t, "Cap1", jobs[0].Cap)
	require.Equal(t, "http://example.com", jobs[0].URL)

	jsonDataWithNamedKey := []byte(`{"key1":[{"conclusion": ":white_check_mark:", "cap": "Cap1", "html_url": "http://example.com"}]}`)
	jobs, err = parseJSON(jsonDataWithNamedKey, "key1")
	require.NoError(t, err)
	require.Len(t, jobs, 1)
	require.Equal(t, ":white_check_mark:", jobs[0].Conclusion)
	require.Equal(t, "Cap1", jobs[0].Cap)
	require.Equal(t, "http://example.com", jobs[0].URL)

	_, err = parseJSON(jsonDataWithNamedKey, "key2")
	require.NoError(t, err)
}

func TestWriteResultsToFileOneSection(t *testing.T) {
	jobs := []ParsedResult{
		{Conclusion: ":white_check_mark:", Cap: "Cap1", URL: "http://example.com"},
	}

	fileName := "output_test.txt"
	err := writeResultsToFile(fileName, "Header1", "Header2", "Section1", jobs)
	require.NoError(t, err)
	defer func() { _ = os.Remove(fileName) }()

	data, err := os.ReadFile(fileName)
	require.NoError(t, err)

	expectedContent := `+----------+---------+
| Header1  | Header2 |
+----------+---------+
|      Section1      |
+----------+---------+
| Cap1     | √       |
+----------+---------+
`
	require.Equal(t, expectedContent, string(data))
}

func TestWriteResultsToFile_Two_Sections(t *testing.T) {
	jobs := []ParsedResult{
		{Conclusion: ":white_check_mark:", Cap: "Cap1", URL: "http://example.com"},
	}

	fileName := "output_two_sections_1.txt"
	err := writeResultsToFile(fileName, "Header1", "Header2", "Section1", jobs)
	require.NoError(t, err)
	defer func() { _ = os.Remove(fileName) }()

	err = writeResultsToFile(fileName, "Header1", "Header2", "Section2", jobs)
	require.NoError(t, err)

	data, err := os.ReadFile(fileName)
	require.NoError(t, err)

	expectedContent := `+----------+---------+
| Header1  | Header2 |
+----------+---------+
|      Section1      |
+----------+---------+
| Cap1     | √       |
+----------+---------+
|      Section2      |
+----------+---------+
| Cap1     | √       |
+----------+---------+
`
	require.Equal(t, expectedContent, string(data))
}

func TestMainFunctionSingleWrite(t *testing.T) {
	content := []byte(`[{"conclusion": ":white_check_mark:", "cap": "Cap1", "html_url": "http://example.com"}]`)
	jsonFileName := "test.json"
	err := os.WriteFile(jsonFileName, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(jsonFileName) }()

	outputFileName := "output_simple.txt"
	defer func() { _ = os.Remove(outputFileName) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", jsonFileName, "--outputFile", outputFileName}
	main()

	data, err := os.ReadFile(outputFileName)
	require.NoError(t, err)

	expectedContent := `+-------+--------+
| Value | Result |
+-------+--------+
| Cap1  | √      |
+-------+--------+
`
	require.Equal(t, expectedContent, string(data))
}

func TestMainFunction_Double_Write_No_Sections_Longer_Rewrite(t *testing.T) {
	content := []byte(`[{"conclusion": ":x:", "cap": "Short", "html_url": "http://example.com"}]`)
	test1Json := "test1.json"
	err := os.WriteFile(test1Json, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(test1Json) }()

	outputFileName := "output_no_sections.txt"
	defer func() { _ = os.Remove(outputFileName) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", test1Json, "--outputFile", outputFileName}
	main()

	content = []byte(`[{"conclusion": ":white_check_mark:", "cap": "I am much longer", "html_url": "http://example.com"}]`)
	test2Json := "test2.json"
	err = os.WriteFile(test2Json, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(test2Json) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", test2Json, "--outputFile", outputFileName}
	main()

	data, err := os.ReadFile(outputFileName)
	require.NoError(t, err)

	expectedContent := `+------------------+--------+
| Value            | Result |
+------------------+--------+
| Short            | X      |
+------------------+--------+
| I am much longer | √      |
+------------------+--------+
`
	fmt.Println(string(data))
	require.Equal(t, expectedContent, string(data))
}

func TestMainFunction_Double_Write_No_Sections_No_Rewrite(t *testing.T) {
	content := []byte(`[{"conclusion": ":white_check_mark:", "cap": "I am much longer", "html_url": "http://example.com"}]`)
	test1Json := "test3.json"
	err := os.WriteFile(test1Json, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(test1Json) }()

	outputFileName := "output_no_sections.txt"
	defer func() { _ = os.Remove(outputFileName) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", test1Json, "--outputFile", outputFileName}
	main()

	content = []byte(`[{"conclusion": ":x:", "cap": "Me short", "html_url": "http://example.com"}]`)
	test2Json := "test4.json"
	err = os.WriteFile(test2Json, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(test2Json) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", test2Json, "--outputFile", outputFileName}
	main()

	data, err := os.ReadFile(outputFileName)
	require.NoError(t, err)

	expectedContent := `+------------------+--------+
| Value            | Result |
+------------------+--------+
| I am much longer | √      |
+------------------+--------+
| Me short         | X      |
+------------------+--------+
`
	fmt.Println(string(data))
	require.Equal(t, expectedContent, string(data))
}

func TestMainFunction_Double_Write_Two_Sections_Rewrite(t *testing.T) {
	content := []byte(`[{"conclusion": ":white_check_mark:", "cap": "Short one", "html_url": "http://example.com"}]`)
	test1Json := "test5.json"
	err := os.WriteFile(test1Json, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(test1Json) }()

	outputFileName := "output_two_sections.txt"
	defer func() { _ = os.Remove(outputFileName) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", test1Json, "--outputFile", outputFileName, "--section", "Section1"}
	main()

	content = []byte(`[{"conclusion": ":x:", "cap": "I am much longer", "html_url": "http://example.com"}]`)
	test2Json := "test6.json"
	err = os.WriteFile(test2Json, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(test2Json) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", test2Json, "--outputFile", outputFileName, "--section", "Section2"}
	main()

	data, err := os.ReadFile(outputFileName)
	require.NoError(t, err)

	expectedContent := `+------------------+--------+
| Value            | Result |
+------------------+--------+
|         Section1          |
+------------------+--------+
| Short one        | √      |
+------------------+--------+
|         Section2          |
+------------------+--------+
| I am much longer | X      |
+------------------+--------+
`
	fmt.Println(string(data))
	require.Equal(t, expectedContent, string(data))
}

func TestMainFunction_Double_Write_Two_Long_Sections_Rewrite(t *testing.T) {
	content := []byte(`[{"conclusion": ":white_check_mark:", "cap": "First short one", "html_url": "http://example.com"}, {"conclusion": ":white_check_mark:", "cap": "Second short one", "html_url": "http://example.com"}]`)
	test1Json := "test5a.json"
	err := os.WriteFile(test1Json, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(test1Json) }()

	outputFileName := "output_two_sections.txt"
	defer func() { _ = os.Remove(outputFileName) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", test1Json, "--outputFile", outputFileName, "--section", "Section1"}
	main()

	content = []byte(`[{"conclusion": ":x:", "cap": "I am much longer", "html_url": "http://example.com"}]`)
	test2Json := "test6a.json"
	err = os.WriteFile(test2Json, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(test2Json) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", test2Json, "--outputFile", outputFileName, "--section", "Section2"}
	main()

	data, err := os.ReadFile(outputFileName)
	require.NoError(t, err)

	expectedContent := `+------------------+--------+
| Value            | Result |
+------------------+--------+
|         Section1          |
+------------------+--------+
| First short one  | √      |
| Second short one | √      |
+------------------+--------+
|         Section2          |
+------------------+--------+
| I am much longer | X      |
+------------------+--------+
`
	fmt.Println(string(data))
	require.Equal(t, expectedContent, string(data))
}

func TestMainFunction_Double_Write_Three_Long_Sections_Rewrite(t *testing.T) {
	content := []byte(`[{"conclusion": ":white_check_mark:", "cap": "First short one", "html_url": "http://example.com"},{"conclusion": ":white_check_mark:", "cap": "Second short one", "html_url": "http://example.com"}, {"conclusion": ":x:", "cap": "Third short one", "html_url": "http://example.com"}]`)
	test1Json := "test5b.json"
	err := os.WriteFile(test1Json, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(test1Json) }()

	outputFileName := "output_two_sections.txt"
	defer func() { _ = os.Remove(outputFileName) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", test1Json, "--outputFile", outputFileName, "--section", "Section1"}
	main()

	content = []byte(`[{"conclusion": ":x:", "cap": "I am much longer", "html_url": "http://example.com"},{"conclusion": ":x:", "cap": "I am much longer", "html_url": "http://example.com"}]`)
	test2Json := "test6b.json"
	err = os.WriteFile(test2Json, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(test2Json) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", test2Json, "--outputFile", outputFileName, "--section", "Section2"}
	main()

	content = []byte(`[{"conclusion": ":white_check_mark:", "cap": "I the last of us", "html_url": "http://example.com"}]`)
	test3Json := "test6c.json"
	err = os.WriteFile(test3Json, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(test3Json) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", test3Json, "--outputFile", outputFileName, "--section", "Section3"}
	main()

	data, err := os.ReadFile(outputFileName)
	require.NoError(t, err)

	expectedContent := `+------------------+--------+
| Value            | Result |
+------------------+--------+
|         Section1          |
+------------------+--------+
| First short one  | √      |
| Second short one | √      |
| Third short one  | X      |
+------------------+--------+
|         Section2          |
+------------------+--------+
| I am much longer | X      |
| I am much longer | X      |
+------------------+--------+
|         Section3          |
+------------------+--------+
| I the last of us | √      |
+------------------+--------+
`
	fmt.Println(string(data))
	require.Equal(t, expectedContent, string(data))
}

func TestMainFunction_Double_Write_Mix_Section_No_Section(t *testing.T) {
	content := []byte(`[{"conclusion": ":white_check_mark:", "cap": "No section here", "html_url": "http://example.com"}]`)
	test1Json := "test7.json"
	err := os.WriteFile(test1Json, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(test1Json) }()

	outputFileName := "output_mix.txt"
	defer func() { _ = os.Remove(outputFileName) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", test1Json, "--outputFile", outputFileName}
	main()

	content = []byte(`[{"conclusion": ":x:", "cap": "With a section", "html_url": "http://example.com"},{"conclusion": ":x:", "cap": "With a section too", "html_url": "http://example.com"}]`)
	test2Json := "test8.json"
	err = os.WriteFile(test2Json, content, 0600)
	require.NoError(t, err)
	defer func() { _ = os.Remove(test2Json) }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--jsonfile", test2Json, "--outputFile", outputFileName, "--section", "Some section"}
	main()

	data, err := os.ReadFile(outputFileName)
	require.NoError(t, err)

	expectedContent := `+--------------------+--------+
| Value              | Result |
+--------------------+--------+
| No section here    | √      |
+--------------------+--------+
|        Some section         |
+--------------------+--------+
| With a section     | X      |
| With a section too | X      |
+--------------------+--------+
`
	fmt.Println(string(data))
	require.Equal(t, expectedContent, string(data))
}
