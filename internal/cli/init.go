package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/lamtuanvu/gh-runner-ctl/internal/config"
)

func newInitCmd() *cobra.Command {
	var importEnv bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create config file in ~/.ghr/",
		RunE: func(cmd *cobra.Command, args []string) error {
			outPath := config.DefaultConfigPath()
			if cfgFile != "" {
				outPath = cfgFile
			}

			if _, err := os.Stat(outPath); err == nil {
				fmt.Printf("Config file %s already exists. Overwrite? [y/N] ", outPath)
				reader := bufio.NewReader(os.Stdin)
				answer, _ := reader.ReadString('\n')
				if strings.TrimSpace(strings.ToLower(answer)) != "y" {
					fmt.Println("Aborted.")
					return nil
				}
			}

			newCfg := config.Default()

			// Try to import from .env if requested or if .env exists
			// Check both ~/.ghr/.env and ./.env
			envPath := ""
			if importEnv {
				envPath = findEnvFile()
				if envPath != "" {
					if err := importFromEnvFile(newCfg, envPath); err != nil {
						fmt.Printf("Warning: could not import %s: %v\n", envPath, err)
					}
				} else {
					fmt.Println("Warning: no .env file found")
				}
			} else {
				envPath = findEnvFile()
				if envPath != "" {
					fmt.Printf("Found %s. Import settings? [Y/n] ", envPath)
					reader := bufio.NewReader(os.Stdin)
					answer, _ := reader.ReadString('\n')
					answer = strings.TrimSpace(strings.ToLower(answer))
					if answer == "" || answer == "y" {
						if err := importFromEnvFile(newCfg, envPath); err != nil {
							fmt.Printf("Warning: could not import %s: %v\n", envPath, err)
						} else {
							fmt.Printf("Imported settings from %s\n", envPath)
						}
					}
				}
			}

			// Interactive prompts
			reader := bufio.NewReader(os.Stdin)

			fmt.Printf("Scope [%s]: ", newCfg.Scope)
			if s := readLine(reader); s != "" {
				newCfg.Scope = s
			}

			if newCfg.Scope == "org" {
				fmt.Printf("Organization [%s]: ", newCfg.Org)
				if s := readLine(reader); s != "" {
					newCfg.Org = s
				}
			} else {
				fmt.Printf("Repo owner [%s]: ", newCfg.Repo.Owner)
				if s := readLine(reader); s != "" {
					newCfg.Repo.Owner = s
				}
				fmt.Printf("Repo name [%s]: ", newCfg.Repo.Name)
				if s := readLine(reader); s != "" {
					newCfg.Repo.Name = s
				}
			}

			fmt.Printf("Token reference [%s]: ", newCfg.Token)
			if s := readLine(reader); s != "" {
				newCfg.Token = s
			}

			fmt.Printf("Runner image [%s]: ", newCfg.Runners.Image)
			if s := readLine(reader); s != "" {
				newCfg.Runners.Image = s
			}

			fmt.Printf("Labels [%s]: ", strings.Join(newCfg.Runners.Labels, ","))
			if s := readLine(reader); s != "" {
				newCfg.Runners.Labels = strings.Split(s, ",")
				for i := range newCfg.Runners.Labels {
					newCfg.Runners.Labels[i] = strings.TrimSpace(newCfg.Runners.Labels[i])
				}
			}

			fmt.Printf("Runner group [%s]: ", newCfg.Runners.Group)
			if s := readLine(reader); s != "" {
				newCfg.Runners.Group = s
			}

			fmt.Printf("Name prefix [%s]: ", newCfg.Runners.NamePrefix)
			if s := readLine(reader); s != "" {
				newCfg.Runners.NamePrefix = s
			}

			if err := config.Save(newCfg, outPath); err != nil {
				return err
			}
			fmt.Printf("Config written to %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().BoolVar(&importEnv, "import-env", false, "import settings from .env file")
	return cmd
}

// findEnvFile returns the first .env file found, checking ~/.ghr/.env then ./.env.
func findEnvFile() string {
	for _, p := range []string{config.DotenvPath(), ".env"} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func importFromEnvFile(cfg *config.Config, path string) error {
	envMap, err := config.ParseDotenv(path)
	if err != nil {
		return err
	}

	if org, ok := envMap["GH_ORG"]; ok {
		cfg.Scope = "org"
		cfg.Org = org
	}
	if labels, ok := envMap["RUNNER_LABELS"]; ok {
		cfg.Runners.Labels = strings.Split(labels, ",")
	}
	if group, ok := envMap["RUNNER_GROUP"]; ok {
		cfg.Runners.Group = group
	}
	// Keep token as env reference
	cfg.Token = "env:GH_TOKEN"

	// Copy all env vars to ~/.ghr/.env (skip if source is already that file)
	dest := config.DotenvPath()
	srcAbs, _ := filepath.Abs(path)
	destAbs, _ := filepath.Abs(dest)
	if srcAbs != destAbs {
		if err := config.SaveDotenv(dest, envMap); err != nil {
			return fmt.Errorf("copying env vars to %s: %w", dest, err)
		}
		fmt.Printf("Environment variables copied to %s\n", dest)
	}

	return nil
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}
