package bubbletea

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/flyteorg/flyte/flytectl/cmd/config/subcommand/project"
	cmdcore "github.com/flyteorg/flyte/flytectl/cmd/core"
	"github.com/flyteorg/flyte/flytectl/pkg/pkce"
	"github.com/flyteorg/flyte/flyteidl/clients/go/admin"

	"github.com/spf13/cobra"
)

type Command struct {
	Cmd   *cobra.Command
	Name  string
	Short string
}

var (
	rootCmd       *cobra.Command
	args          []string
	isCommand     = true
	unhandleFlags []string
	existingFlags []string
	listErrMsg    error = nil
)
var (
	domainName    = [3]string{"development", "staging", "production"}
	nameToCommand = map[string]Command{}
)

// Generate a []list.Item of cmd's subcommands
func generateSubCmdItems(cmd *cobra.Command) []list.Item {
	items := []list.Item{}

	for _, subcmd := range cmd.Commands() {
		subCmdName := strings.Fields(subcmd.Use)[0]
		nameToCommand[subCmdName] = Command{
			Cmd:   subcmd,
			Name:  subCmdName,
			Short: subcmd.Short,
		}
		items = append(items, item(subCmdName))
	}

	return items
}

// Generate list.Model for domain names
func genDomainListModel(m *listModel) {
	items := []list.Item{}
	for _, domain := range domainName {
		items = append(items, item(domain))
	}

	m.list = genList(items, "Please choose one of the domains")
}

// Get the "get" "project" cobra.Command item
func extractGetProjectCmd() *cobra.Command {
	var getCmd *cobra.Command
	var getProjectCmd *cobra.Command

	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "get" {
			getCmd = cmd
			break
		}
	}
	for _, cmd := range getCmd.Commands() {
		if cmd.Use == "project" {
			getProjectCmd = cmd
			break
		}
	}
	return getProjectCmd
}

// Get all the project names from the configured endpoint
func getProjects(getProjectCmd *cobra.Command) ([]string, error) {
	ctx := context.Background()
	err := rootCmd.PersistentPreRunE(rootCmd, []string{})
	if err != nil {
		return nil, err
	}
	adminCfg := admin.GetConfig(ctx)

	clientSet, err := admin.ClientSetBuilder().WithConfig(admin.GetConfig(ctx)).
		WithTokenCache(pkce.TokenCacheKeyringProvider{
			ServiceUser: fmt.Sprintf("%s:%s", adminCfg.Endpoint.String(), pkce.KeyRingServiceUser),
			ServiceName: pkce.KeyRingServiceName,
		}).Build(ctx)
	if err != nil {
		return nil, err
	}
	cmdCtx := cmdcore.NewCommandContext(clientSet, getProjectCmd.OutOrStdout())

	projects, err := cmdCtx.AdminFetcherExt().ListProjects(ctx, project.DefaultConfig.Filter)
	if err != nil {
		return nil, err
	}

	projectNames := []string{}
	for _, p := range projects.Projects {
		projectNames = append(projectNames, p.Id)
	}

	return projectNames, nil
}

// Generate list.Model for project names from the configured endpoint
func genProjectListModel(m *listModel) error {
	getProjectCmd := extractGetProjectCmd()
	projects, err := getProjects(getProjectCmd)
	if err != nil {
		return err
	}

	items := []list.Item{}
	for _, project := range projects {
		items = append(items, item(project))
	}

	m.list = genList(items, "Please choose one of the projects")

	return nil
}

// Generate list.Model of options for different flags
func genFlagListModel(m *listModel) error {
	i := 0
	// If flag already specified by user, skip.
	for i < len(unhandleFlags) {
		if !rootCmd.Flags().ShorthandLookup(strings.TrimPrefix(unhandleFlags[i], "-")).Changed {
			break
		}
		i++
	}
	if i == len(unhandleFlags) {
		m.quitting = true
		return nil
	}

	flag := unhandleFlags[i]
	args = append(args, flag)
	unhandleFlags = unhandleFlags[i+1:]

	switch flag {
	case "-p":
		err := genProjectListModel(m)
		if err != nil {
			return err
		}
	case "-d":
		genDomainListModel(m)
	}

	return nil
}

// Generate list.Model of subcommands from a given command
func genCmdListModel(m *listModel, c *cobra.Command) {
	items := generateSubCmdItems(c)
	l := genList(items, "")
	m.list = l

}

// Generate list.Model after user chose one of the item
func genListModel(m *listModel, item string) error {
	// Still in the stage of handling subcommands
	if isCommand {
		var ok bool
		// check if we reach a runnable command
		if unhandleFlags, ok = commandFlagMap[sliceToString(args)]; !ok {
			genCmdListModel(m, nameToCommand[item].Cmd)
			return nil
		}
		isCommand = false
	}

	// Handled all flags, quit.
	if len(unhandleFlags) == 0 {
		m.quitting = true
		return nil
	}

	// Still have flags to handle
	err := genFlagListModel(m)
	if err != nil {
		return err
	}

	return nil
}

func ifRunBubbleTea() (*cobra.Command, bool, error) {
	cmd, flags, err := rootCmd.Find(os.Args[1:])
	if err != nil {
		return cmd, false, err
	}
	existingFlags = flags

	err = rootCmd.ParseFlags(flags)
	if err != nil {
		return cmd, false, err
	}

	if ok, err := rootCmd.Flags().GetBool("interactive"); !ok || err != nil {
		return cmd, false, nil
	}

	tempCmd := cmd
	for tempCmd.HasParent() {
		args = append([]string{tempCmd.Use}, args...)
		tempCmd = tempCmd.Parent()
	}

	return cmd, true, nil
}