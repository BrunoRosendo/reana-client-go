/*
This file is part of REANA.
Copyright (C) 2022 CERN.

REANA is free software; you can redistribute it and/or modify it
under the terms of the MIT License; see LICENSE file for more details.
*/

package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/spf13/cobra"
)

var FilesBlacklist = []string{".git/", "/.git/"}
var InteractiveSessionTypes = []string{"jupyter"}

func GetRunStatuses(includeDeleted bool) []string {
	runStatuses := []string{
		"created",
		"running",
		"finished",
		"failed",
		"stopped",
		"queued",
		"pending",
	}
	if includeDeleted {
		runStatuses = append(runStatuses, "deleted")
	}
	return runStatuses
}

func ExecuteCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()

	return buf.String(), err
}

func NewRequest(token string, serverURL string, endpoint string) []byte {
	// disable certificate security checks
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	url := serverURL + endpoint + "?access_token=" + token
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	return respBytes
}

func ParseFilterParameters(filter []string, filterNames []string) ([]string, string) {
	searchFilters := make(map[string][]string)
	var statusFilters []string

	for _, value := range filter {
		if !strings.Contains(value, "=") {
			fmt.Println("Error: Wrong input format. Please use --filter filter_name=filter_value")
			os.Exit(1)
		}

		filterNameAndValue := strings.SplitN(value, "=", 2)
		filterName := strings.ToLower(filterNameAndValue[0])
		filterValue := filterNameAndValue[1]

		if !slices.Contains(filterNames, filterName) {
			fmt.Printf("Error: Filter %s is not valid", filterName)
			os.Exit(1)
		}

		if filterName == "status" && !slices.Contains(GetRunStatuses(true), filterValue) {
			fmt.Printf("Error: Input status value %s is not valid. ", filterValue)
			os.Exit(1)
		}

		if filterName == "status" {
			statusFilters = append(statusFilters, filterValue)
		} else {
			searchFilters[filterName] = append(searchFilters[filterName], filterValue)
		}
	}

	searchFiltersString := ""
	if len(searchFilters) > 0 {
		searchFiltersByteArray, err := json.Marshal(searchFilters)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		searchFiltersString = string(searchFiltersByteArray)
	}

	return statusFilters, searchFiltersString
}

func ParseFormatParameters(filters []string) map[string]string {
	parsedFilters := make(map[string]string)
	for _, filter := range filters {
		filterNameAndValue := strings.SplitN(filter, "=", 2)
		if len(filterNameAndValue) < 2 {
			parsedFilters[filterNameAndValue[0]] = ""
		} else {
			parsedFilters[filterNameAndValue[0]] = filterNameAndValue[1]
		}
	}
	return parsedFilters
}

func FormatData(data *[][]any, header *[]string, formatFilters map[string]string) {
	if len(formatFilters) == 0 {
		return
	}

	for col := 0; col < len(*header); col++ {
		filter, include := formatFilters[(*header)[col]]

		if !include {
			// Remove the column from header and data
			*header = append((*header)[:col], (*header)[col+1:]...)
			for row := range *data {
				(*data)[row] = append((*data)[row][:col], (*data)[row][col+1:]...)
			}
			col--
		} else if filter != "" {
			// Remove rows not containing filter
			for row := 0; row < len(*data); row++ {
				val := fmt.Sprint((*data)[row][col])
				if val != filter {
					*data = append((*data)[:row], (*data)[row+1:]...)
					row--
				}
			}
		}
	}
}

func HasAnyPrefix(s string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func FromIsoToTimestamp(date string) time.Time {
	timestamp, err := time.Parse("2006-01-02T15:04:05", date)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	return timestamp
}

func GetWorkflowNameAndRunNumber(workflowName string) (string, string) {
	workflowNameAndRunNumber := strings.SplitN(workflowName, ".", 2)
	if len(workflowNameAndRunNumber) < 2 {
		return workflowName, ""
	}
	return workflowNameAndRunNumber[0], workflowNameAndRunNumber[1]
}

func FormatSessionURI(serverURL string, path string, token string) string {
	return serverURL + path + "?token=" + token
}
