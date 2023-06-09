package cmd

import (
	"github.com/ai-lenok/gitlab-adapter/maintainer"
	"github.com/ai-lenok/gitlab-adapter/properties"
	"github.com/ai-lenok/gitlab-adapter/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	cfgFile string

	reqListPipelines properties.ReqListPipelines

	reqCreateRepo properties.ReqCreateRepo
	reqDeleteRepo properties.ReqDeleteRepo

	rootCmd = &cobra.Command{
		Use:   "gitlab-adapter",
		Short: "Adapter between Learning Management System and GitLab",
	}

	createRepoCmd = &cobra.Command{
		Use:   "create-repo",
		Short: "Create new repository",
		Run:   createRepo,
	}

	deleteRepoCmd = &cobra.Command{
		Use:   "delete-repo",
		Short: "Delete repository",
		Run:   deleteRepo,
	}

	verifyPipelineStatusCmd = &cobra.Command{
		Use:   "verify-pipeline",
		Short: "Verify that the latest build in the repository was successful",
		Run:   verifyPipelineStatus,
	}

	startServerCmd = &cobra.Command{
		Use:   "start-server",
		Short: "Start web server app",
		Run:   startServer,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootFlags := rootCmd.PersistentFlags()
	rootFlags.StringVarP(&cfgFile, "config", "c", "", "Path to configuration file (default is './config/application.yaml')")
	rootFlags.String("gitlab.host", "", "GitLab host")
	rootFlags.String("gitlab.token", "", "Authorization token")
	viper.BindPFlags(rootFlags)

	verifyPipelineFlags := verifyPipelineStatusCmd.PersistentFlags()
	verifyPipelineFlags.StringVar(&reqListPipelines.ProjectId, "project-id", "", "Project ID for checking repo")
	verifyPipelineStatusCmd.MarkPersistentFlagRequired("project-id")

	createRepoFlags := createRepoCmd.PersistentFlags()
	createRepoFlags.StringVar(&reqCreateRepo.Name, "name", "", "Display name of new repository")
	createRepoFlags.StringVar(&reqCreateRepo.Path, "path", "", "Real name of new repository")
	createRepoFlags.StringVar(&reqCreateRepo.Description, "description", "", "Description of new repository")
	createRepoFlags.StringVar(&reqCreateRepo.NamespaceId, "namespace", "", "Parent repo's group ID")
	createRepoFlags.StringVar(&reqCreateRepo.ImportUrl, "import-url", "", "Url to open repository for cloning")
	createRepoCmd.MarkPersistentFlagRequired("path")
	createRepoCmd.MarkPersistentFlagRequired("namespace")

	deleteRepoFlags := deleteRepoCmd.PersistentFlags()
	deleteRepoFlags.StringVar(&reqDeleteRepo.ProjectId, "project-id", "", "Repository ID")
	deleteRepoCmd.MarkPersistentFlagRequired("project-id")

	serverFlags := startServerCmd.PersistentFlags()
	serverFlags.Int("server.port", 8080, "Server port, default 8080")
	viper.BindPFlags(serverFlags)

	rootCmd.AddCommand(verifyPipelineStatusCmd)
	rootCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(createRepoCmd)
	rootCmd.AddCommand(deleteRepoCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./config")
		viper.SetConfigType("yaml")
		viper.SetConfigName("application")
	}

	viper.SetEnvPrefix("ga")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func createMaintainer() maintainer.Maintainer {
	return maintainer.Maintainer{
		Client: &http.Client{},
		Config: &properties.GitLabConfig{
			Host:      viper.GetString("gitlab.host"),
			AuthToken: viper.GetString("gitlab.token"),
		},
	}
}

func createRepo(cmd *cobra.Command, args []string) {
	m := createMaintainer()

	resp, err := m.CreateRepo(&reqCreateRepo)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
}

func deleteRepo(cmd *cobra.Command, args []string) {
	m := createMaintainer()

	resp, err := m.DeleteRepo(&reqDeleteRepo)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}

func verifyPipelineStatus(cmd *cobra.Command, args []string) {
	m := createMaintainer()

	log.Printf("Gitlab: %s. Config: %s. Host: %s. Token: %s", viper.Get("gitlab"), m.Config,
		viper.GetString("gitlab.host"), viper.GetString("gitlab.token"))

	isSuccess, err := m.LastBuildIsSuccess(&reqListPipelines)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("is Success: " + strconv.FormatBool(isSuccess))
	if isSuccess {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func startServer(cmd *cobra.Command, args []string) {
	m := createMaintainer()

	err := server.StartServer(
		&m,
		viper.GetInt("server.port"))

	if err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
