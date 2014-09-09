package action

import (
	"net/http"

	"github.com/vti/twigo/app/model"
)

type ListArticlesArchive struct {
	BaseAction
}

type MonthInfo struct {
	Name      string
	Documents []*model.Document
}
type YearInfo struct {
	Name   string
	Months []MonthInfo
}

var monthsNames = [...]string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

func (action *ListArticlesArchive) Execute(w http.ResponseWriter, r *http.Request) {
	home := action.Context.App.Home

	dm := &model.DocumentManager{Root: home + "/articles/"}
	documents, err := dm.LoadDocuments(0, "")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	years := []YearInfo{}

	for _, document := range documents {
		yearName := document.Created["Year"].String()
		monthName := monthsNames[document.Created["Month"]-1]

		yearFound := false
		for i, yearInfo := range years {
			if yearInfo.Name == yearName {

				monthFound := false
				for j, monthInfo := range yearInfo.Months {
					if monthInfo.Name == monthName {
						yearInfo.Months[j].Documents = append(monthInfo.Documents, document)

						monthFound = true
						break
					}
				}

				if !monthFound {
					years[i].Months = append(yearInfo.Months, MonthInfo{
						Name:      monthName,
						Documents: []*model.Document{document}})
				}

				yearFound = true
				break
			}
		}

		if !yearFound {
			month := MonthInfo{Name: monthName,
				Documents: []*model.Document{document}}
			years = append(years, YearInfo{
				Name:   yearName,
				Months: []MonthInfo{month}})
		}
	}

	action.Context.SetTemplateName("layouts/html")
	action.Context.SetTemplateFiles([]string{
		"layouts/html.tpl",
		"archive.tpl",
		"article-meta.tpl",
	})

	vars := map[string]interface{}{
		"Conf":  action.Context.App.Conf,
		"Years": years,
	}
	action.Context.SetTemplateVars(vars)
}
