/*
This file is part of REANA.
Copyright (C) 2022 CERN.

REANA is free software; you can redistribute it and/or modify it
under the terms of the MIT License; see LICENSE file for more details.
*/

package cmd

import (
	"fmt"
	"os"
	"reanahub/reana-client-go/client"
	"reanahub/reana-client-go/client/operations"
	"reanahub/reana-client-go/utils"
	"reanahub/reana-client-go/validation"
	"strings"

	"github.com/spf13/cobra"
)

const listFormatFlagDesc = `Format output according to column titles or column
values. Use <columm_name>=<column_value> format.
E.g. display workflow with failed status and named test_workflow
--format status=failed,name=test_workflow.
`

const listFilterFlagDesc = `Filter workflow that contains certain filtering
criteria. Use --filter
<columm_name>=<column_value> pairs. Available
filters are 'name' and 'status'.
`

const listDesc = `
List all workflows and sessions.

The ` + "``list``" + ` command lists workflows and sessions. By default, the list of
workflows is returned. If you would like to see the list of your open
interactive sessions, you need to pass the ` + "``--sessions``" + ` command-line
option.

Example:

  $ reana-client list --all

  $ reana-client list --sessions

  $ reana-client list --verbose --bytes
`

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all workflows and sessions.",
		Long:  listDesc,
		Run: func(cmd *cobra.Command, args []string) {
			token, _ := cmd.Flags().GetString("access-token")
			if token == "" {
				token = os.Getenv("REANA_ACCESS_TOKEN")
			}
			serverURL := os.Getenv("REANA_SERVER_URL")
			validation.ValidateAccessToken(token)
			validation.ValidateServerURL(serverURL)
			list(cmd)
		},
	}

	cmd.Flags().StringP("access-token", "t", "", "Access token of the current user.")
	cmd.Flags().StringP("workflow", "w", "", "List all runs of the given workflow.")
	cmd.Flags().StringP("sessions", "s", "", "List all open interactive sessions.")
	cmd.Flags().String("format", "", listFormatFlagDesc)
	cmd.Flags().BoolP("json", "", false, "Get output in JSON format.")
	cmd.Flags().StringArray("filter", []string{}, listFilterFlagDesc)

	return cmd
}

func list(cmd *cobra.Command) {
	token, _ := cmd.Flags().GetString("access-token")
	if token == "" {
		token = os.Getenv("REANA_ACCESS_TOKEN")
	}
	workflow, _ := cmd.Flags().GetString("workflow")
	jsonOutput, _ := cmd.Flags().GetBool("json")
	filter, _ := cmd.Flags().GetStringArray("filter")

	filterNames := []string{"name", "status"}
	statusFilters, searchFilter := utils.ParseListFilters(filter, filterNames)

	listParams := operations.NewGetWorkflowsParams()
	listParams.SetAccessToken(&token)
	listParams.SetWorkflowIDOrName(&workflow)
	listParams.SetStatus(statusFilters)
	listParams.SetSearch(&searchFilter)

	listResp, err := client.ApiClient().Operations.GetWorkflows(listParams)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	if jsonOutput {
		utils.DisplayJsonOutput(listResp.Payload)
	} else {
		displayListPayload(listResp.Payload)
	}
}

func displayListPayload(p *operations.GetWorkflowsOKBody) {
	header := []interface{}{
		"NAME",
		"RUN_NUMBER",
		"CREATED",
		"STARTED",
		"ENDED",
		"STATUS",
	}
	var rows [][]interface{}

	for _, workflow := range p.Items {
		var row []interface{}
		workflowNameAndRunNumber := strings.SplitN(workflow.Name, ".", 2)
		row = append(
			row,
			workflowNameAndRunNumber[0],
			workflowNameAndRunNumber[1],
			workflow.Created,
			displayOptionalField(workflow.Progress.RunStartedAt),
			displayOptionalField(workflow.Progress.RunFinishedAt),
			workflow.Status,
		)
		rows = append(rows, row)
	}

	utils.DisplayTable(header, rows)
}

func displayOptionalField(value *string) string {
	if value == nil {
		return "-"
	}
	return *value
}
