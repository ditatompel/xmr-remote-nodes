package server

import "github.com/ditatompel/xmr-remote-nodes/cmd"

func init() {
	cmd.Root.AddCommand(serveCmd)
	cmd.Root.AddCommand(cronCmd)
	cmd.Root.AddCommand(probersCmd)
	probersCmd.AddCommand(listProbersCmd)
	probersCmd.AddCommand(addProbersCmd)
	probersCmd.AddCommand(editProbersCmd)
	probersCmd.AddCommand(deleteProbersCmd)
	listProbersCmd.Flags().StringP("sort-by", "s", "last_submit_ts", "Sort by column name, can be id or last_submit_ts")
	listProbersCmd.Flags().StringP("sort-dir", "d", "desc", "Sort direction, can be asc or desc")
}
