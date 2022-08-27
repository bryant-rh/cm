package output

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

//Write
func Write(data [][]string, out_type string, outType bool) {
	//var table *tablewriter.Table

	table := table(outType)
	var header []string
	//header = append(header, "NODE")
	switch {
	case out_type == "project":
		header = append(header,
			"ID", "Project_ID", "Project Name", "Description", "CreatedAt", "UpdatedAt",
		)
	case out_type == "cluster":
		header = append(header,
			"ID", "Cluster_ID", "Cluster Name", "Description", "Labels", "CreatedAt", "UpdatedAt",
		)
	case out_type == "sa":
		header = append(header,
			"ID", "Sa_ID", "SaName", "SaToken", "NameSpace", "CreatedAt", "UpdatedAt",
		)
	case out_type == "label":
		header = append(header,
			"Project Name","Cluster Name","Labels",
		)
	case out_type == "token":
		header = append(header,
			"SaToken",
		)
	default:
		header = append(header,
			"ID", "Project_ID", "Project Name", "Description", "CreatedAt", "UpdatedAt",
		)
	}
	table.SetHeader(header)
	for _, i := range data {
		table.Append(i)
	}
	table.Render()

}

//table
func table(outType bool) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	if outType {
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t") // pad with tabs
		table.SetNoWhiteSpace(true)
		return table
	}
	return table
}
