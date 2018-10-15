package main

import (
	"time"

	"github.com/kolo/ledger-go/datetime"
	"github.com/shopspring/decimal"
)

type weeklyReportItem struct {
	week   datetime.Week
	report report

	next *weeklyReportItem
}

func (ri *weeklyReportItem) update(r *record) {
	if updated, found := ri.report[r.credit.name]; found {
		updated.increase(r.amount)
	}
}

// weeklyReport is a sorted doubly linked list.
type weeklyReport struct {
	assets []string
	head   *weeklyReportItem
}

func (rp *weeklyReport) update(r *record) {
	rp.findOrCreateReportItem(r.recordedAt).update(r)
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
		r[asset] = &reportItem{
			account: &account{
				name:  asset,
				asset: true,
			},
			total: decimal.Zero,
		}
	}

	return r
}

func weeklyExpensesReport(rd recordReader, assets []string) {
	report := &weeklyReport{assets: assets}

	for {
		r := rd.Next()
		if r == nil {
			break
		}

		report.update(r)
	}

	printWeeklyReport(report)
}
