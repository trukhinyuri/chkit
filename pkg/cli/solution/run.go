package clisolution

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/containerum/chkit/pkg/context"
	solutionControls "github.com/containerum/chkit/pkg/controls/solution"
	"github.com/containerum/chkit/pkg/model/solution"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/pairs"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var aliases = []string{"sol", "solutions", "sols", "solu", "so"}

func Run(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:     "solution",
		Aliases: aliases,
		Short:   "Run solution from template",
		Example: "chkit run solution [$PUBLIC_SOLUTION] [--env=KEY1:VALUE1,KEY2:VALUE2] [--file $FILENAME] [--force]",
		Run: func(cmd *cobra.Command, args []string) {
			sol := buildSolution(ctx, cmd, args)
			if force, _ := cmd.Flags().GetBool("force"); force {
				if err := ctx.Client.RunSolution(sol); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Printf("Solution %s is ready to run\n", sol.Name)
				return
			}
			solutions, err := ctx.Client.GetSolutionsTemplatesList()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			config := solutionControls.WizardConfig{
				EditName:  true,
				Templates: solutions.Names(),
				Solution:  &sol,
			}
			sol = solutionControls.WizardF(ctx, config)
			if activekit.YesNo("Are you sure you want to run solution %s?", sol.Name) {
				for k := range sol.Env {
					if sol.Env[k] == "" {
						delete(sol.Env, k)
					}
				}

				if err := ctx.Client.RunSolution(sol); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Printf("Solution %s is ready to run\n", sol.Name)
				return
			}
		},
	}
	command.PersistentFlags().
		BoolP("force", "f", false, "create solution without confirmation, optional")
	command.PersistentFlags().
		String("file", "", "file with solution data, .yaml or .json, stdin if '-', optional")
	command.PersistentFlags().
		String("env", "", "solution environment variables, optional")
	command.PersistentFlags().
		String("name", "", "solution name, optional, autogenerated if void")
	command.PersistentFlags().
		String("branch", "master", "solution git repo branch, optional")
	return command
}

func buildSolution(ctx *context.Context, cmd *cobra.Command, args []string) solution.Solution {
	var sol solution.Solution
	var flags = cmd.Flags()
	if flags.Changed("file") {
		sol = solutionFromFile(cmd)
	} else if len(args) == 1 {
		sol.Template = args[0]
	} else if force, _ := flags.GetBool("force"); force {
		cmd.Help()
		os.Exit(1)
	}
	if flags.Changed("name") {
		sol.Name, _ = flags.GetString("name")
	}
	if flags.Changed("namespace") {
		sol.Namespace, _ = flags.GetString("namespace")
	} else {
		sol.Namespace = ctx.Namespace.ID
	}
	if flags.Changed("branch") {
		sol.Branch, _ = flags.GetString("branch")
	} else {
		sol.Branch = "master"
	}
	if flags.Changed("env") {
		envString, _ := flags.GetString("env")
		env, err := pairs.ParseMap(envString, ":")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		sol.Env = env
	}
	return sol
}

func solutionFromFile(cmd *cobra.Command) solution.Solution {
	flags := cmd.Flags()
	fName, _ := flags.GetString("file")
	var data = func() []byte {
		if fName == "-" {
			buf := &bytes.Buffer{}
			_, err := buf.ReadFrom(os.Stdin)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return buf.Bytes()
		}
		data, err := ioutil.ReadFile(fName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return data
	}()
	var sol solution.Solution
	if path.Ext(fName) == "json" {
		if err := json.Unmarshal(data, &sol); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if path.Ext(fName) == "yaml" {
		if err := yaml.Unmarshal(data, &sol); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Error: invalid file name %q, want extensions 'yaml' or 'json'\n%s", fName, cmd.Flag("file").Usage)
		os.Exit(1)
	}
	return sol
}
