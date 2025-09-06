package pdfutil

import (
	"bytes"
	"text/template"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/templates"
)

type InvoiceViewData struct {
	Number        string
	Date          string
	CustomerName  string
	Items         []model.InvoiceItem
	Total         float64
	ProForma      bool
	Subtotal      float64
	TaxTotal      float64
	DiscountTotal float64
}

func GenerateInvoicePDF(invoice *model.Invoice, proForma bool) ([]byte, error) {
	// Parse template
	tmpl, err := template.New("invoice.html").Funcs(funcMap).ParseFS(templates.InvoiceFS, templates.InvoicePath)
	if err != nil {
		return nil, err
	}

	items := make([]model.InvoiceItem, len(invoice.Items))
	for i, item := range invoice.Items {
		items[i] = *item
	}

	customerName := ""
	if invoice.Order != nil && invoice.Order.Customer != nil {
		customerName = invoice.Order.Customer.FirstName + " " + invoice.Order.Customer.LastName
	}

	data := InvoiceViewData{
		Number:        invoice.InvoiceNumber,
		Date:          time.Now().Format("02 Jan 2006"),
		CustomerName:  customerName,
		Items:         items,
		Total:         invoice.Total,
		ProForma:      proForma,
		Subtotal:      invoice.Subtotal,
		TaxTotal:      invoice.TaxTotal,
		DiscountTotal: invoice.DiscountTotal,
	}

	// Render HTML
	var htmlBuf bytes.Buffer
	if err := tmpl.Execute(&htmlBuf, data); err != nil {
		return nil, err
	}

	// Convert to PDF with wkhtmltopdf
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(htmlBuf.Bytes())))
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}
