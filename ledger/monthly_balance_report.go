package ledger

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/kolo/ledger-go/datetime"
	"github.com/shopspring/decimal"
)

type monthlyReportItem struct {
	month  datetime.Month
	report *balanceReport

	next *monthlyReportItem
}

func (ri *monthlyReportItem) update(r *Record) {
	ri.report.update(r)
}

type monthlyReport struct {
	assets []*Account
	head   *monthlyReportItem
}

func (rp *monthlyReport) update(r *Record) {
	rp.findOrCreateReportItem(r.Date).update(r)
}

func (rp *monthlyReport) findOrCreateReportItem(t time.Time) *monthlyReportItem {
	if rp.head == nil {
		rp.head = &monthlyReportItem{
			month:  datetime.NewMonth(t),
			report: newBalanceReport(rp.assets),
		}
		return rp.head
	}

	var prev, ri *monthlyReportItem
	m := datetime.NewMonth(t)
	ri = rp.head
	for {
		if ri == nil {
			ri = rp.insert(m, prev, nil)
			break
		}

		if ri.month.Same(m) {
			break
		}

		if ri.month.After(m) {
			ri = rp.insert(m, prev, ri)
			break
		}

		prev = ri
		ri = ri.next
	}

	return ri
}

func (rp *monthlyReport) insert(m datetime.Month, prev, next *monthlyReportItem) *monthlyReportItem {
	ri := &monthlyReportItem{
		month:  m,
		report: newBalanceReport(rp.assets),
	}

	prev.next = ri
	ri.next = next

	return ri
}

func (rp *monthlyReport) print(output io.Writer) {
	w := tabwriter.NewWriter(output, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Month\tBalance\n")

	total := decimal.Zero
	for ri := rp.head; ri != nil; ri = ri.next {
		monthTotal := ri.report.total()
		total = total.Add(monthTotal)
		fmt.Fprintf(w, "%v\t%s\n", ri.month, monthTotal.StringFixed(2))
	}

	fmt.Fprintf(w, "-----\t\n")
	fmt.Fprintf(w, "Total\t%s\n", total.StringFixed(2))

	w.Flush()
}

func MonthlyBalanceReport(iter RecordIterator, assets []*Account, output io.Writer) {
	report := monthlyReport{assets: assets}
	for {
		r := iter.Next()
		if r == nil {
			break
		}

		report.update(r)
	}
	report.print(output)
}
