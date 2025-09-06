package templates

import "embed"

//go:embed invoice/invoice.html
var InvoiceFS embed.FS
var InvoicePath = "invoice/invoice.html"
