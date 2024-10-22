package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type ParsedResult struct {
	Conclusion string `json:"conclusion"`
	Cap        string `json:"cap"`
	URL        string `json:"html_url"`
}

type section struct {
	name string
	jobs []ParsedResult
}

type ResultsMap map[string][]ParsedResult

func parseJSON(jsonData []byte, namedKey string) ([]ParsedResult, error) {
	var jobs []ParsedResult
	if namedKey == "" {
		err := json.Unmarshal(jsonData, &jobs)
		return jobs, err
	}

	var results ResultsMap
	err := json.Unmarshal(jsonData, &results)
	if err != nil {
		return nil, err
	}

	if val, ok := results[namedKey]; ok {
		return val, nil
	}

	return jobs, nil
}

func calculateColumnWidths(firstColumnHeader, secondColumnHeader string, sections []section) (int, int) {
	maxFirstColumnLen := len(firstColumnHeader)
	maxSecondColumnLen := len(secondColumnHeader)

	for _, s := range sections {
		if len(s.name) > maxFirstColumnLen {
			maxFirstColumnLen = len(s.name)
		}
		for _, job := range s.jobs {
			if len(job.Cap) > maxFirstColumnLen {
				maxFirstColumnLen = len(job.Cap)
			}
		}
	}

	return maxFirstColumnLen, maxSecondColumnLen
}

func writeResultsToFile(fileName string, firstColumnHeader, secondColumnHeader, currentSection string, jobs []ParsedResult) error {
	orderedSections := make([]section, 0)
	if _, err := os.Stat(fileName); err == nil {
		data, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		lines := strings.Split(string(data), "\n")
		var sectionName string
		for _, line := range lines {
			if strings.HasPrefix(line, "|") {
				parts := strings.Split(line, "|")
				if len(parts) == 3 { // It's a section header
					sectionName = strings.TrimSpace(parts[1])
				} else if len(parts) == 4 { // It might be a job entry, but can also be the header
					if strings.TrimSpace(parts[1]) == firstColumnHeader || strings.TrimSpace(parts[2]) == secondColumnHeader {
						continue
					}
					parsedResults := []ParsedResult{
						{
							Cap:        strings.TrimSpace(parts[1]),
							Conclusion: strings.TrimSpace(parts[2]),
						}}

					sectionFound := false
					for i, orderedSection := range orderedSections {
						if orderedSection.name == sectionName {
							orderedSections[i].jobs = append(orderedSections[i].jobs, parsedResults...)
							sectionFound = true
							break
						}
					}
					if !sectionFound {
						orderedSections = append(orderedSections, section{name: sectionName, jobs: parsedResults})
					}
				}
			}
		}
	}

	orderedSections = append(orderedSections, section{name: currentSection, jobs: jobs})
	maxFirstColumnLen, maxSecondColumnLen := calculateColumnWidths(firstColumnHeader, secondColumnHeader, orderedSections)

	firstColumnFormat := fmt.Sprintf("%%-%ds", maxFirstColumnLen)
	secondColumnFormat := fmt.Sprintf("%%-%ds", maxSecondColumnLen)
	rowFormat := fmt.Sprintf("| %s | %s |\n", firstColumnFormat, secondColumnFormat)
	separator := fmt.Sprintf("+-%s-+-%s-+\n", strings.Repeat("-", maxFirstColumnLen), strings.Repeat("-", maxSecondColumnLen))

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	header := separator
	header += fmt.Sprintf(rowFormat, firstColumnHeader, secondColumnHeader)
	header += separator
	_, err = file.WriteString(header)
	if err != nil {
		return err
	}

	for _, s := range orderedSections {
		if s.name != "" {
			sectionHeader := fmt.Sprintf("| %s |\n", centerText(s.name, maxFirstColumnLen+maxSecondColumnLen+3))
			_, err = file.WriteString(sectionHeader)
			if err != nil {
				return err
			}
			_, err = file.WriteString(separator)
			if err != nil {
				return err
			}
		}

		for _, job := range s.jobs {
			result := "X"
			if job.Conclusion == ":white_check_mark:" || job.Conclusion == "√" {
				result = "√"
			}
			line := fmt.Sprintf(rowFormat, job.Cap, result)
			_, err = file.WriteString(line)
			if err != nil {
				return err
			}
		}

		_, err = file.WriteString(separator)
		if err != nil {
			return err
		}
	}

	return nil
}

func centerText(s string, width int) string {
	spaces := (width - len(s)) / 2
	if (spaces*2+width+len(s))%2 == 0 {
		return strings.Repeat(" ", spaces) + s + strings.Repeat(" ", spaces)
	}
	return strings.Repeat(" ", spaces) + s + strings.Repeat(" ", spaces+1)
}

func main() {
	firstColumnHeader := flag.String("firstColumn", "Value", "Header for the first column")
	secondColumnHeader := flag.String("secondColumn", "Result", "Header for the second column")
	jsonFileFlag := flag.String("jsonfile", "", "Path to JSON input file")
	section := flag.String("section", "", "Optional section name")
	namedKey := flag.String("namedKey", "", "Optional named key to look for in the JSON input")
	outputFile := flag.String("outputFile", "", "Optional output file to save results (default: output.txt)")

	flag.Parse()

	if *jsonFileFlag == "" {
		panic(fmt.Errorf("please provide a path to the JSON input file using --jsonfile flag"))
	}

	jsonFile, err := os.ReadFile(*jsonFileFlag)
	if err != nil {
		panic(fmt.Errorf("error reading JSON file: %v", err))
	}

	jobs, err := parseJSON(jsonFile, *namedKey)
	if err != nil {
		panic(fmt.Errorf("error parsing JSON file: %v", err))
	}

	if len(jobs) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "No results found in the JSON file")
		return
	}

	outputFileName := "output.txt"
	if *outputFile != "" {
		outputFileName = *outputFile
	}

	err = writeResultsToFile(outputFileName, *firstColumnHeader, *secondColumnHeader, *section, jobs)
	if err != nil {
		panic(fmt.Errorf("error writing to file: %v", err))
	}

	msg := fmt.Sprintf("Found results for '%s'. Updating file %s", *namedKey, outputFileName)
	if *namedKey == "" {
		msg = fmt.Sprintf("Results updated successfully in %s", outputFileName)
	}

	fmt.Println(msg)
}
