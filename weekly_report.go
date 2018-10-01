package main

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type weekID struct {
	year int
	week int
}

func (id weekID) same(other weekID) bool {
	return id.year == other.year && id.week == other.week
}

func (id weekID) before(other weekID) bool {
	if id.year < other.year {
		return true
	}

	if id.year == other.year && id.week < other.week {
		return true
	}

	return false
}

func (id weekID) after(other weekID) bool {
	return !id.before(other)
}

func newWeekID(t time.Time) weekID {
	// FIXME: t.ISOWeek() function return 52 for Jan 01 to Jan 03 and might return
	// 1 for Dec 29 to Dec 31. This is ok for now, but should be fixed at some point.
	year, week := t.ISOWeek()
	return weekID{year: year, week: week}
}

type weeklyReportItem struct {
	id     weekID
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
			id:     newWeekID(t),
			report: rp.newReport(),
		}

		return rp.head
	}

	var prev, cur *weeklyReportItem

	id := newWeekID(t)
	cur = rp.head
	for {
		if cur == nil {
			cur = rp.insert(id, prev, nil)
			break
		}

		if cur.id.same(id) {
			break
		}

		if cur.id.after(id) {
			cur = rp.insert(id, prev, cur)
			break
		}

		prev = cur
		cur = cur.next
	}

	return cur
}

func (rp *weeklyReport) insert(id weekID, prev, next *weeklyReportItem) *weeklyReportItem {
	ri := &weeklyReportItem{
		id:     id,
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

func weeklyExpensesReport(rd recordReader, assets []string, filter filterFunc) {
	report := &weeklyReport{assets: assets}

	for {
		r := rd.Next()
		if r == nil {
			break
		}

		r = filter(r)
		if r == nil {
			continue
		}

		report.update(r)
	}

	for ri := report.head; ri != nil; ri = ri.next {
		fmt.Println(ri.id, ri.report.total())
	}
}
