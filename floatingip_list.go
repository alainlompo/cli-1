package cli

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newFloatingIPListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list [FLAGS]",
		Short:                 "List Floating IPs",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runFloatingIPList),
	}
	return cmd
}

func runFloatingIPList(cli *CLI, cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetStringArray("output")
	outOpts, err := parseOutputOpts(out)
	if err != nil {
		return err
	}

	floatingIPs, err := cli.Client().FloatingIP.All(cli.Context)
	if err != nil {
		return err
	}

	cols := []string{"id", "type", "description", "ip", "home", "server", "dns"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := newTableOutput().
		AddAllowedFields(hcloud.FloatingIP{}).
		AddFieldOutputFn("dns", fieldOutputFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			var dns string
			if len(floatingIP.DNSPtr) == 1 {
				for _, v := range floatingIP.DNSPtr {
					dns = v
				}
			}
			if len(floatingIP.DNSPtr) > 1 {
				dns = fmt.Sprintf("%d entries", len(floatingIP.DNSPtr))
			}
			return na(dns)
		})).
		AddFieldOutputFn("server", fieldOutputFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			var server string
			if floatingIP.Server != nil {
				server = strconv.Itoa(floatingIP.Server.ID)
			}
			return na(server)
		})).
		AddFieldOutputFn("home", fieldOutputFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			return floatingIP.HomeLocation.Name
		}))

	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, floatingIP := range floatingIPs {
		tw.Write(cols, floatingIP)
	}
	tw.Flush()
	return nil
}
