package handlers

import (
	"fmt"
	"math"
	"time"
	"work-at-olist/storage"
)

type Call struct {
	*storage.Call
}

func (h *Handler) CalculateCall(dateStart, dateEnd time.Time) (float64, error) {
	var hour int
	var start, end time.Time

	if !inTimeRange(dateStart) {
		hour = dateStart.Hour()
		if hour >= 22 {
			dateStart = dateStart.AddDate(0, 0, 1)
		}
		start = time.Date(dateStart.Year(), dateStart.Month(), dateStart.Day(), 6, 0, 0, 0, time.UTC)
	} else {
		start = dateStart
	}

	if !inTimeRange(dateEnd) {
		hour = dateEnd.Hour()
		if hour < 6 {
			dateEnd = dateEnd.AddDate(0, 0, -1)
		}
		end = time.Date(dateEnd.Year(), dateEnd.Month(), dateEnd.Day(), 22, 0, 0, 0, time.UTC)
	} else {
		end = dateEnd
	}

	du := end.Sub(start)
	x := math.Floor(du.Hours()/24) * (8 * 60)

	t := float64(int(du.Minutes()) - int(x))
	rst := math.Round((t*0.09+0.36)*100) / 100

	return rst, nil
}

func (h *Handler) SaveCall(rs, re storage.Record) error {

	dateStart, err := time.Parse(time.RFC3339, rs.Timestamp)
	if err != nil {
		return err
	}
	dateEnd, err := time.Parse(time.RFC3339, re.Timestamp)
	if err != nil {
		return err
	}

	cc, err := h.CalculateCall(dateStart, dateEnd)
	if err != nil {
		return err
	}

	bill, err := h.DB.GetBillByPeriod(rs.Source, int(dateEnd.Month()), dateEnd.Year())
	if err != nil {
		return err
	}

	duration := dateEnd.Sub(dateStart)
	hr, mn, sc := dateStart.Clock()
	startTime := fmt.Sprintf("%d:%d:%d", hr, mn, sc)

	if bill.Id != 0 {
		c := storage.NewCall(bill.Id, rs.Destination, duration.String(), dateStart.Format("2006-01-02"), startTime, cc)
		if err := h.DB.CreateCall(c); err != nil {
			return err
		}
	} else {
		b := storage.NewBill(rs.Source, int(dateEnd.Month()), dateEnd.Year())
		bid, err := h.DB.CreateBill(b)
		if err != nil {
			return err
		}
		c := storage.NewCall(int(bid), rs.Destination, duration.String(), dateStart.Format("2006-01-02"), startTime, cc)
		if err := h.DB.CreateCall(c); err != nil {
			return err
		}
	}
	return nil
}
