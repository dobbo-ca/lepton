package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dobbo-ca/lepton/internal/api"
	"github.com/dobbo-ca/lepton/internal/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "leptonctl",
	Short: "lepton — AI agent orchestration platform",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		dsn := viper.GetString("db.dsn")
		if dsn == "" {
			dsn = "lepton.db"
		}
		database, err := db.Open(dsn)
		if err != nil {
			return fmt.Errorf("open database: %w", err)
		}
		addr := viper.GetString("server.addr")
		if addr == "" {
			addr = ":8080"
		}
		srv := api.New(database)
		fmt.Fprintf(os.Stdout, "leptonctl: listening on %s\n", addr)
		return http.ListenAndServe(addr, srv)
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		dsn := viper.GetString("db.dsn")
		if dsn == "" {
			dsn = "lepton.db"
		}
		database, err := db.Open(dsn)
		if err != nil {
			return fmt.Errorf("open database: %w", err)
		}
		if err := database.Migrate(); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
		fmt.Fprintln(os.Stdout, "migrations complete")
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stdout, "leptonctl v0.1.0")
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("config", "", "config file (default: lepton.yaml)")
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.AddCommand(serveCmd, migrateCmd, versionCmd)
}

func initConfig() {
	cfgFile := viper.GetString("config")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("lepton")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
	}
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
