package ledger

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/kolo/ledger-go/datetime"
	"github.com/shopspring/decimal"
)

type weeklyReportItem struct {
	week   datetime.Week
	report report

	next *weeklyReportItem
}

func (ri *weeklyReportItem) update(r *Record) {
	if updated, found := ri.report[r.Credit.Name]; found {
		updated.increase(r.Amount)
	}
}

// weeklyReport is a sorted double linked list.
type weeklyReport struct {
	assets []*Account
	head   *weeklyReportItem
}

func (r *weeklyReport) print(output io.Writer) {
	// Sort assets
	assets := make([]string, len(r.assets))
	for i, asset := range r.assets {
		assets[i] = asset.Name
	}
	sort.Strings(assets)

	w := tabwriter.NewWriter(output, 0, 0, 2, ' ', 0)

	// Print header
	fmt.Fprintf(w, "\t%s\tTotal\n", strings.Join(assets, "\t"))

	// Print report body
	totals := make([]decimal.Decimal, len(assets))
	for ri := r.head; ri != nil; ri = ri.next {
		subtotals := []string{}
		total := decimal.Zero
		for i, asset := range assets {
			total = total.Add(ri.report[asset].total)
			totals[i] = totals[i].Add(ri.report[asset].total)
			subtotals = append(subtotals, ri.report[asset].total.StringFixed(2))
		}
		fmt.Fprintf(w, "%v\t%s\t%s\n", ri.week, strings.Join(subtotals, "\t"), total.StringFixed(2))
	}

	// Print footer
	sum := decimal.Zero
	stringTotals := make([]string, len(totals))
	for i, total := range totals {
		sum = sum.Add(total)
		stringTotals[i] = total.StringFixed(2)
	}
	fmt.Fprintf(w, "-----\t%s\n", strings.Join(make([]string, len(totals)), "\t"))
	fmt.Fprintf(w, "Total\t%s\t%s\n", strings.Join(stringTotals, "\t"), sum.StringFixed(2))

	w.Flush()
}
func (rp *weeklyReport) update(r *Record) {
	rp.findOrCreateReportItem(r.Date).update(r)
}

func (rp *weeklyReport) findOrCreateReportItem(t time.Time) *weeklyReportItem {
	if rp.head == nil {
		rp.head = &weeklyReportItem{
			week:   datetime.NewWeek(t),
			report: rp.newReport(),
		}

		return rp.head
	}

	var prev, cur *weeklyReportItem

	id := datetime.NewWeek(t)
	cur = rp.head
	for {
		if cur == nil {
			cur = rp.insert(id, prev, nil)
			break
		}

		if cur.week.Same(id) {
			break
		}

		if cur.week.After(id) {
			cur = rp.insert(id, prev, cur)
			break
		}

		prev = cur
		cur = cur.next
	}

	return cur
}

func (rp *weeklyReport) insert(w datetime.Week, prev, next *weeklyReportItem) *weeklyReportItem {
	ri := &weeklyReportItem{
		week:   w,
		report: rp.newReport(),
	}

	prev.next = ri
	ri.next = next

	return ri
}

func (rp *weeklyReport) newReport() report {
	r := report{}

	for _, asset := range rp.assets {
		r[asset.Name] = &reportItem{
			account: asset,
			total:   decimal.Zero,
		}
	}

	return r
}

type reportItem struct {
	account *Account
	total   decimal.Decimal
}

func (ri *reportItem) increase(amount decimal.Decimal) {
	ri.total = ri.total.Add(amount)
}

func (ri *reportItem) decrease(amount decimal.Decimal) {
	ri.total = ri.total.Sub(amount)
}

type report map[string]*reportItem

func (r report) update(rec *Record) {
	r.item(rec.Credit).decrease(rec.Amount)
	r.item(rec.Debit).increase(rec.Amount)
}

func (r report) item(ac *Account) *reportItem {
	item, found := r[ac.Name]
	if !found {
		item = &reportItem{
			account: ac,
			total:   decimal.Zero,
		}
		r[ac.Name] = item
	}

	return item
}

func (r report) total() decimal.Decimal {
	total := decimal.Zero
	for _, ri := range r {
		total = total.Add(ri.total)
	}

	return total
}

func WeeklyExpensesReport(iter RecordIterator, assets []*Account, output io.Writer) {
	iter = NewFilteredIterator(iter, filterExpenses)
	report := &weeklyReport{assets: assets}

	for {
		r := iter.Next()
		if r == nil {
			break
		}

		report.update(r)
	}

	report.print(output)
}
