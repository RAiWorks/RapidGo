package cli

import (
	"fmt"

	"github.com/RAiWorks/RapidGo/v2/core/container"
	"github.com/RAiWorks/RapidGo/v2/core/service"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var dbSeedCmd = &cobra.Command{
	Use:   "db:seed",
	Short: "Seed the database with records",
	RunE: func(cmd *cobra.Command, args []string) error {
		application := NewApp(service.ModeAll)
		db := container.MustMake[*gorm.DB](application.Container, "db")

		if seederFn == nil {
			return fmt.Errorf("no seeder registered — call cli.SetSeeder() in main.go")
		}

		name, _ := cmd.Flags().GetString("seeder")
		if err := seederFn(db, name); err != nil {
			return err
		}
		if name != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Seeder %s complete.\n", name)
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "Database seeding complete.")
		}
		return nil
	},
}

func init() {
	dbSeedCmd.Flags().String("seeder", "", "Run a specific seeder by name")
}
