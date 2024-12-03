package forecast

import (
	"fmt"
	"strings"
)

type ForecastKind string

const (
	PeriodToday     ForecastKind = "Today"
	PeriodAfternoon ForecastKind = "Afternoon"
	PeriodTonight   ForecastKind = "Tonight"
)

// Simplified forecast report builder
type Simplified struct {
	periods []SimplifiedPeriod
}

func (sf *Simplified) AddPeriod(sp SimplifiedPeriod) {
	sf.periods = append(sf.periods, sp)
}

func (sf *Simplified) AsString() string {
	out := strings.Builder{}
	itemCount := 0
	joinStatement := "and"
	for _, p := range sf.periods {
		// Add separator between forecast state variations and change join statement
		if itemCount > 0 {
			out.WriteString(", ")
			joinStatement = "with"
		}
		out.WriteString(
			fmt.Sprintf(
				"%s %s %s %s", periodKindToPhrase(p), periodTempToPhrase(p, itemCount > 0), joinStatement, p.shortDesc))
		itemCount++
	}
	return out.String()
}

// SimplifiedPeriod provides necessary information for forecast builder to create a report
type SimplifiedPeriod struct {
	kind      ForecastKind
	temp      int
	shortDesc string
}

func (sp *SimplifiedPeriod) SetKind(fk ForecastKind) {
	sp.kind = fk
}

func (sp *SimplifiedPeriod) SetTemp(temp int) {
	sp.temp = temp
}

func (sp *SimplifiedPeriod) SetShortDesc(desc string) {
	sp.shortDesc = desc
}

// Helpers

// Transform variety of kinds into the better phrasing
func periodKindToPhrase(p SimplifiedPeriod) string {
	switch p.kind {
	case PeriodToday, PeriodTonight:
		return "For " + string(p.kind)
	case PeriodAfternoon:
		return "In the " + string(p.kind)
	default:
		return string(p.kind)
	}
}

// Transform variety of temperatures into the alpha representation with context of following part of the forecast
func periodTempToPhrase(p SimplifiedPeriod, useSuffix bool) string {
	var tempType string
	var tempSuffix string
	switch {
	case p.temp < 50:
		tempType = "cold"
		tempSuffix = "er"
	case p.temp >= 50 && p.temp < 75:
		tempType = "moderate"
	default:
		tempType = "hot"
		tempSuffix = "ter"
	}

	if useSuffix {
		return fmt.Sprintf("expecting %s%s temperatures", tempType, tempSuffix)
	}
	return fmt.Sprintf("expecting %s temperature", tempType)
}
