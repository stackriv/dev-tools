package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func Invoice(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/invoice" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "invoice", model.PageData{
			Title:       "Invoice Generator",
			Description: "Generate professional invoices for your clients.",
			Page:        "invoice",
		})
	}
}

func GenerateInvoice(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/invoice" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body model.InvoiceData
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		err1 := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: "invalid request"}})
		fmt.Println(http.StatusBadRequest, "invalid request")
		return
	}

	if body.Date == "" {
		body.Date = time.Now().Format("2006-01-02")
	}
	if body.Number == "" {
		body.Number = fmt.Sprintf("INV-%s", time.Now().Format("20060102"))
	}
	if body.Currency == "" {
		body.Currency = "USD"
	}

	symbols := map[string]string{
		"USD": "$", "EUR": "€", "GBP": "£", "XOF": "FCFA", "CAD": "CA$",
	}
	symbol := symbols[body.Currency]
	if symbol == "" {
		symbol = body.Currency
	}

	// Calculate
	var subtotal float64
	for _, item := range body.Items {
		subtotal += item.Quantity * item.UnitPrice
	}
	tax := subtotal * body.TaxRate / 100
	total := subtotal + tax

	var rows strings.Builder
	for _, item := range body.Items {
		amount := item.Quantity * item.UnitPrice
		rows.WriteString(fmt.Sprintf(`
			<tr>
				<td>%s</td>
				<td class="text-right">%.2f</td>
				<td class="text-right">%s %.2f</td>
				<td class="text-right">%s %.2f</td>
			</tr>`,
			pkg.EscapeHTML(item.Description),
			item.Quantity,
			symbol, item.UnitPrice,
			symbol, amount,
		))
	}

	taxRow := ""
	if body.TaxRate > 0 {
		taxRow = fmt.Sprintf(`<tr><td colspan="3" class="text-right">Tax (%.1f%%)</td><td class="text-right">%s %.2f</td></tr>`, body.TaxRate, symbol, tax)
	}

	notes := ""
	if body.Notes != "" {
		notes = fmt.Sprintf(`<div class="invoice-notes"><h4>Notes</h4><p>%s</p></div>`, pkg.EscapeHTML(body.Notes))
	}

	dueDate := ""
	if body.DueDate != "" {
		dueDate = fmt.Sprintf(`<p><strong>Due Date:</strong> %s</p>`, pkg.EscapeHTML(body.DueDate))
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<title>Invoice %s</title>
		<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: 'Segoe UI', Arial, sans-serif; color: #1e293b; background: #fff; padding: 40px; max-width: 800px; margin: 0 auto; }
		.invoice-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 40px; padding-bottom: 24px; border-bottom: 2px solid #6366f1; }
		.invoice-title { font-size: 2rem; font-weight: 800; color: #6366f1; }
		.invoice-meta { text-align: right; font-size: 0.9rem; color: #64748b; }
		.invoice-meta p { margin-bottom: 4px; }
		.invoice-parties { display: grid; grid-template-columns: 1fr 1fr; gap: 40px; margin-bottom: 32px; }
		.invoice-party h3 { font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.1em; color: #6366f1; margin-bottom: 8px; }
		.invoice-party p { font-size: 0.9rem; color: #475569; line-height: 1.6; }
		.invoice-party .name { font-weight: 700; color: #1e293b; font-size: 1rem; }
		table { width: 100%%; border-collapse: collapse; margin-bottom: 24px; }
		th { background: #6366f1; color: white; padding: 10px 12px; text-align: left; font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.05em; }
		th.text-right, td.text-right { text-align: right; }
		td { padding: 10px 12px; border-bottom: 1px solid #e2e8f0; font-size: 0.9rem; }
		tr:last-child td { border-bottom: none; }
		tr:nth-child(even) { background: #f8fafc; }
		.invoice-totals { margin-left: auto; width: 280px; }
		.invoice-totals table { margin-bottom: 0; }
		.invoice-totals td { border-bottom: 1px solid #e2e8f0; }
		.invoice-totals tr.total { font-weight: 700; font-size: 1.1rem; background: #f1f5f9; }
		.invoice-totals tr.total td { border-bottom: none; color: #6366f1; }
		.invoice-notes { margin-top: 32px; padding: 16px; background: #f8fafc; border-left: 3px solid #6366f1; border-radius: 4px; }
		.invoice-notes h4 { font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.05em; color: #6366f1; margin-bottom: 6px; }
		.invoice-notes p { font-size: 0.9rem; color: #475569; }
		.invoice-footer { margin-top: 40px; padding-top: 16px; border-top: 1px solid #e2e8f0; text-align: center; font-size: 0.8rem; color: #94a3b8; }
		</style>
		</head>
		<body>
		<div class="invoice-header">
			<div>
				<div class="invoice-title">INVOICE</div>
				<div style="font-size:0.9rem;color:#64748b;margin-top:4px">%s</div>
			</div>
			<div class="invoice-meta">
				<p><strong>Invoice #:</strong> %s</p>
				<p><strong>Date:</strong> %s</p>
				%s
			</div>
		</div>
		 
		<div class="invoice-parties">
			<div class="invoice-party">
				<h3>From</h3>
				<p class="name">%s</p>
				<p>%s</p>
				<p>%s</p>
				<p>%s</p>
			</div>
			<div class="invoice-party">
				<h3>Bill To</h3>
				<p class="name">%s</p>
				<p>%s</p>
				<p>%s</p>
			</div>
		</div>
		 
		<table>
			<thead>
				<tr>
					<th>Description</th>
					<th class="text-right">Qty</th>
					<th class="text-right">Unit Price</th>
					<th class="text-right">Amount</th>
				</tr>
			</thead>
			<tbody>%s</tbody>
		</table>
		 
		<div class="invoice-totals">
			<table>
				<tr><td colspan="1">Subtotal</td><td class="text-right">%s %.2f</td></tr>
				%s
				<tr class="total"><td>Total</td><td class="text-right">%s %.2f</td></tr>
			</table>
		</div>
		 
		%s
		 
		<div class="invoice-footer">
			<p>Generated by Stackriv Dev Tools</p>
		</div>
		</body>
		</html>`,
		body.Number,
		pkg.EscapeHTML(body.From.Name),
		pkg.EscapeHTML(body.Number),
		pkg.EscapeHTML(body.Date),
		dueDate,
		pkg.EscapeHTML(body.From.Name),
		pkg.EscapeHTML(body.From.Email),
		pkg.EscapeHTML(body.From.Address),
		pkg.EscapeHTML(body.From.Phone),
		pkg.EscapeHTML(body.To.Name),
		pkg.EscapeHTML(body.To.Email),
		pkg.EscapeHTML(body.To.Address),
		rows.String(),
		symbol, subtotal,
		taxRow,
		symbol, total,
		notes,
	)

	err := json.NewEncoder(w).Encode(map[string]string{"html": html})
	if err != nil {
		err1 := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: err1["msg"]}})
		fmt.Println(err1)
	}
}
