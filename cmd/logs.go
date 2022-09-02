/*
This file is part of REANA.
Copyright (C) 2022 CERN.

REANA is free software; you can redistribute it and/or modify it
under the terms of the MIT License; see LICENSE file for more details.
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"reanahub/reana-client-go/client"
	"reanahub/reana-client-go/client/operations"
	"reanahub/reana-client-go/pkg/config"
	"reanahub/reana-client-go/pkg/displayer"
	"reanahub/reana-client-go/pkg/filterer"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

const logsDesc = `
Get workflow logs.

The ` + "``logs``" + ` command allows to retrieve logs of running workflow. Note that
only finished steps of the workflow are returned, the logs of the currently
processed step is not returned until it is finished.

Examples:

$ reana-client logs -w myanalysis.42

$ reana-client logs -w myanalysis.42 -s 1st_ste
`

const logsFilterFlagDesc = `Filter job logs to include only those steps that
match certain filtering criteria. Use --filter
name=value pairs. Available filters are
compute_backend, docker_img, status and step.`

// logs struct that contains the logs of a workflow.
// Pointers used for nullable values
type logs struct {
	WorkflowLogs   *string               `json:"workflow_logs"`
	JobLogs        map[string]jobLogItem `json:"job_logs"`
	EngineSpecific *string               `json:"engine_specific"`
}

// jobLogItem struct that contains the log information of a job.
type jobLogItem struct {
	WorkflowUuid   string  `json:"workflow_uuid"`
	JobName        string  `json:"job_name"`
	ComputeBackend string  `json:"compute_backend"`
	BackendJobId   string  `json:"backend_job_id"`
	DockerImg      string  `json:"docker_img"`
	Cmd            string  `json:"cmd"`
	Status         string  `json:"status"`
	Logs           string  `json:"logs"`
	StartedAt      *string `json:"started_at"`
	FinishedAt     *string `json:"finished_at"`
}

type logsOptions struct {
	token      string
	workflow   string
	jsonOutput bool
	filters    []string
	page       int64
	size       int64
}

// newLogsCmd creates a command to get workflow logs.
func newLogsCmd(api *client.API) *cobra.Command {
	o := &logsOptions{}

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Get workflow logs.",
		Long:  logsDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(cmd, api)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&o.token, "access-token", "t", "", "Access token of the current user.")
	f.StringVarP(
		&o.workflow,
		"workflow",
		"w",
		"",
		"Name or UUID of the workflow. Overrides value of REANA_WORKON environment variable.",
	)
	f.BoolVar(&o.jsonOutput, "json", false, "Get output in JSON format.")
	f.StringSliceVar(&o.filters, "filter", []string{}, logsFilterFlagDesc)
	f.Int64Var(&o.page, "page", 1, "Results page number (to be used with --size).")
	f.Int64Var(&o.size, "size", 0, "Size of results per page (to be used with --page).")

	return cmd
}

func (o *logsOptions) run(cmd *cobra.Command, api *client.API) error {
	filters, err := parseLogsFilters(o.filters)
	if err != nil {
		return err
	}
	steps, err := filters.GetMulti("step")
	if err != nil {
		return err
	}

	logsParams := operations.NewGetWorkflowLogsParams()
	logsParams.SetAccessToken(&o.token)
	logsParams.SetWorkflowIDOrName(o.workflow)
	logsParams.SetPage(&o.page)
	logsParams.SetSteps(steps)
	if cmd.Flags().Changed("size") {
		logsParams.SetSize(&o.size)
	}

	logsResp, err := api.Operations.GetWorkflowLogs(logsParams)
	if err != nil {
		return err
	}

	var workflowLogs logs
	err = json.Unmarshal([]byte(logsResp.GetPayload().Logs), &workflowLogs)
	if err != nil {
		return err
	}

	err = filterJobLogs(&workflowLogs.JobLogs, filters)
	if err != nil {
		return err
	}

	if o.jsonOutput {
		err := displayer.DisplayJsonOutput(workflowLogs, cmd.OutOrStdout())
		if err != nil {
			return err
		}
	} else {
		displayHumanFriendlyLogs(cmd, workflowLogs, steps)
	}

	return nil
}

// parseLogsFilters parses a list of filters in the format 'filter=value', for the 'logs' command.
// Returns an error if any of the given filters are not valid.
func parseLogsFilters(filterInput []string) (filterer.Filters, error) {
	filters, err := filterer.NewFilters(
		config.LogsSingleFilters,
		config.LogsMultiFilters,
		filterInput,
	)
	if err != nil {
		return filters, err
	}

	err = filters.ValidateValues("status", config.GetRunStatuses(true))
	if err != nil {
		return filters, err
	}

	err = filters.ValidateValues("compute_backend", config.ReanaComputeBackendKeys)
	if err != nil {
		return filters, err
	}

	return filters, nil
}

// filterJobLogs returns a subset of jobLogs based on the given filters.
func filterJobLogs(jobLogs *map[string]jobLogItem, filters filterer.Filters) error {
	// Convert to a map based on json properties
	var jobLogsMap map[string]map[string]string
	jsonLogs, err := json.Marshal(jobLogs)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonLogs, &jobLogsMap)
	if err != nil {
		return err
	}

	var unwantedLogs []string
	for jobLogKey, jobLogValue := range jobLogsMap {
		for _, filterKey := range filters.SingleFilterKeys {
			filterValue, _ := filters.GetSingle(filterKey)
			if filterKey == "compute_backend" {
				filterValue = config.ReanaComputeBackends[strings.ToLower(filterValue)]
			}
			if filterValue != "" && jobLogValue[filterKey] != filterValue {
				unwantedLogs = append(unwantedLogs, jobLogKey)
				break
			}
		}
	}

	for _, log := range unwantedLogs {
		delete(*jobLogs, log)
	}
	return nil
}

// displayHumanFriendlyLogs displays the logs in a human friendly way.
func displayHumanFriendlyLogs(cmd *cobra.Command, logs logs, steps []string) {
	if logs.WorkflowLogs != nil && *logs.WorkflowLogs != "" {
		displayLogHeader(cmd, "Workflow engine logs")
		cmd.Println(*logs.WorkflowLogs)
	}

	if logs.EngineSpecific != nil && *logs.EngineSpecific != "" {
		displayLogHeader(cmd, "Engine internal logs")
		cmd.Println(*logs.EngineSpecific)
	}

	if len(steps) > 0 {
		var returnedStepNames, missingStepNames []string
		for _, jobItem := range logs.JobLogs {
			returnedStepNames = append(returnedStepNames, jobItem.JobName)
		}

		for _, step := range steps {
			if !slices.Contains(returnedStepNames, step) {
				missingStepNames = append(missingStepNames, step)
			}
		}

		if len(missingStepNames) > 0 {
			errMsg := fmt.Sprintf(
				"The logs of step(s) %s were not found, check for spelling mistakes in the step names",
				strings.Join(missingStepNames, ","),
			)
			displayer.DisplayMessage(errMsg, displayer.Error, false, cmd.ErrOrStderr())
		}
	}

	if len(logs.JobLogs) > 0 {
		displayLogHeader(cmd, "Job logs")
		for jobId, jobItem := range logs.JobLogs {
			jobNameOrId := jobId
			if jobItem.JobName != "" {
				jobNameOrId = jobItem.JobName
			}
			displayer.PrintColorable(
				fmt.Sprintf("%s Step: %s\n", config.LeadingMark, jobNameOrId),
				cmd.OutOrStdout(),
				text.Bold,
				displayer.JobStatusToColor[jobItem.Status],
			)

			displayLogItem(cmd, &jobItem.WorkflowUuid, "Workflow ID", jobItem.Status)
			displayLogItem(cmd, &jobItem.ComputeBackend,
				"Compute backend", jobItem.Status)
			displayLogItem(cmd, &jobItem.BackendJobId, "Job ID", jobItem.Status)
			displayLogItem(cmd, &jobItem.DockerImg, "Docker image", jobItem.Status)
			displayLogItem(cmd, &jobItem.Cmd, "Command", jobItem.Status)
			displayLogItem(cmd, &jobItem.Status, "Status", jobItem.Status)
			displayLogItem(cmd, jobItem.StartedAt, "Started", jobItem.Status)
			displayLogItem(cmd, jobItem.FinishedAt, "Finished", jobItem.Status)

			if jobItem.Logs != "" {
				logsItem := "\n" + jobItem.Logs // break line after title
				displayLogItem(cmd, &logsItem, "Logs", jobItem.Status)
			} else {
				msg := fmt.Sprintf("Step %s emitted no logs.", jobNameOrId)
				displayer.DisplayMessage(msg, displayer.Info, false, cmd.OutOrStdout())
			}
		}
	}
}

// displayLogItem displays an optional log item if it is not nil or an empty string.
// The title is displayed according to the color associated with the job's status.
func displayLogItem(cmd *cobra.Command, item *string, title, status string) {
	if item == nil || *item == "" {
		return
	}
	displayer.PrintColorable(
		fmt.Sprintf("%s %s: ", config.LeadingMark, title),
		cmd.OutOrStdout(),
		displayer.JobStatusToColor[status],
	)
	cmd.Println(*item)
}

// displayLogHeader displays a header for a group of logs represented by title.
func displayLogHeader(cmd *cobra.Command, title string) {
	displayer.PrintColorable(
		fmt.Sprintf("\n%s %s\n", config.LeadingMark, title),
		cmd.OutOrStdout(),
		text.Bold,
		text.FgYellow,
	)
}
