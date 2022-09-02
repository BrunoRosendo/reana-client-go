/*
This file is part of REANA.
Copyright (C) 2022 CERN.

REANA is free software; you can redistribute it and/or modify it
under the terms of the MIT License; see LICENSE file for more details.
*/

package cmd

import (
	"errors"
	"reanahub/reana-client-go/client"
	"reanahub/reana-client-go/client/operations"
	"reanahub/reana-client-go/pkg/config"
	"reanahub/reana-client-go/pkg/datautils"
	"reanahub/reana-client-go/pkg/displayer"
	"reanahub/reana-client-go/pkg/filterer"

	"github.com/spf13/cobra"
)

const duDesc = `
Get workspace disk usage.

The ` + "``du``" + ` command allows to chech the disk usage of given workspace.

Examples:

  $ reana-client du -w myanalysis.42 -s

  $ reana-client du -w myanalysis.42 -s --human-readable

  $ reana-client du -w myanalysis.42 --filter name=data/
`

const duFilterFlagDesc = `Filter results to show only files that match certain filtering
criteria such as file name or size.
Use --filter <columm_name>=<column_value> pairs.
Available filters are 'name' and 'size'.`

type duOptions struct {
	token         string
	workflow      string
	summarize     bool
	humanReadable bool
	filter        []string
}

// newDuCmd creates a command to get workspace disk usage.
func newDuCmd(api *client.API) *cobra.Command {
	o := &duOptions{}

	cmd := &cobra.Command{
		Use:   "du",
		Short: "Get workspace disk usage.",
		Long:  duDesc,
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
		"w", "",
		"Name or UUID of the workflow. Overrides value of REANA_WORKON environment variable.",
	)
	f.BoolVarP(&o.summarize, "summarize", "s", false, "Display total.")
	f.BoolVarP(
		&o.humanReadable,
		"human-readable",
		"h",
		false,
		"Show disk size in human readable format.",
	)
	f.StringSliceVar(&o.filter, "filter", []string{}, duFilterFlagDesc)
	// Remove -h shorthand
	cmd.PersistentFlags().BoolP("help", "", false, "Help for du")

	return cmd
}

func (o *duOptions) run(cmd *cobra.Command, api *client.API) error {
	filters, err := filterer.NewFilters(nil, config.DuMultiFilters, o.filter)
	if err != nil {
		return err
	}
	searchFilter, err := filters.GetJson(config.DuMultiFilters)
	if err != nil {
		return err
	}

	duParams := operations.NewGetWorkflowDiskUsageParams()
	duParams.SetAccessToken(&o.token)
	duParams.SetWorkflowIDOrName(o.workflow)
	additionalParams := operations.GetWorkflowDiskUsageBody{
		Summarize: o.summarize,
		Search:    searchFilter,
	}
	duParams.SetParameters(additionalParams)

	duResp, err := api.Operations.GetWorkflowDiskUsage(duParams)
	if err != nil {
		return err
	}

	err = displayDuPayload(cmd, duResp.Payload, o.humanReadable)
	if err != nil {
		return err
	}
	return nil
}

// displayDuPayload displays the disk usage payload, according to the humanReadable flag.
func displayDuPayload(
	cmd *cobra.Command,
	p *operations.GetWorkflowDiskUsageOKBody,
	humanReadable bool,
) error {
	if len(p.DiskUsageInfo) == 0 {
		return errors.New("no files matching filter criteria")
	}

	header := []string{"SIZE", "NAME"}
	var rows [][]any

	for _, diskUsageInfo := range p.DiskUsageInfo {
		if datautils.HasAnyPrefix(diskUsageInfo.Name, config.FilesBlacklist) {
			continue
		}

		var row []any
		if humanReadable {
			row = append(row, diskUsageInfo.Size.HumanReadable)
		} else {
			row = append(row, diskUsageInfo.Size.Raw)
		}
		row = append(row, "."+diskUsageInfo.Name)
		rows = append(rows, row)
	}

	displayer.DisplayTable(header, rows, cmd.OutOrStdout())
	return nil
}
