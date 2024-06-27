package cli

import (
	"context"

	"github.com/enrichman/kubectl-rancher-migration/pkg/client"
	v1_10_0 "github.com/enrichman/kubectl-rancher-migration/pkg/migrations/v1_10_0"
	apiv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/spf13/cobra"
)

func NewV1_10_0_Cmd(c *client.RancherClient) (*cobra.Command, error) {
	lConn := &client.LdapClient{}
	adConfig := &apiv3.ActiveDirectoryConfig{}

	cmd := &cobra.Command{
		Use:          "v1.10.0",
		Short:        "v1.10.0",
		Long:         `Handle v1.10.0 migration`,
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			adConfig2 := &apiv3.ActiveDirectoryConfig{}
			err := c.Rancher.Get().Resource("authconfigs").Name("activedirectory").Do(context.Background()).Into(adConfig2)
			if err != nil {
				return err
			}
			*adConfig = *adConfig2

			ldapConfig, err := client.NewLDAPConfigFromActiveDirectory(c.Kube.CoreV1(), adConfig)
			if err != nil {
				return err
			}

			conn, err := client.NewLDAPConn(ldapConfig)
			if err != nil {
				return err
			}
			*lConn = client.LdapClient{Conn: conn}

			return nil
		},
	}

	cmd.AddCommand(
		NewV1_10_0_CheckCmd(c, lConn, adConfig),
		NewV1_10_0_MigrateCmd(c, lConn, adConfig),
		NewV1_10_0_RollbackCmd(c, lConn, adConfig),
	)

	return cmd, nil
}

func NewV1_10_0_CheckCmd(c *client.RancherClient, lConn *client.LdapClient, adConfig *apiv3.ActiveDirectoryConfig) *cobra.Command {
	return &cobra.Command{
		Use:          "check",
		Short:        "check",
		Long:         `v1.10.0 migration check`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return v1_10_0.Check(c, lConn, adConfig)
		},
	}
}

func NewV1_10_0_MigrateCmd(c *client.RancherClient, lConn *client.LdapClient, adConfig *apiv3.ActiveDirectoryConfig) *cobra.Command {
	return &cobra.Command{
		Use:          "migrate",
		Short:        "migrate",
		Long:         `v1.10.0 migration`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return v1_10_0.Migrate(c, lConn, adConfig)
		},
	}
}

func NewV1_10_0_RollbackCmd(c *client.RancherClient, lConn *client.LdapClient, adConfig *apiv3.ActiveDirectoryConfig) *cobra.Command {
	return &cobra.Command{
		Use:          "rollback",
		Short:        "rollback",
		Long:         `v1.10.0 rollback`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return v1_10_0.Rollback(c, lConn, adConfig)
		},
	}
}