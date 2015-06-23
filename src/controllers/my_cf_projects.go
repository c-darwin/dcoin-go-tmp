package controllers
import (
	"utils"
	"log"
	"sort"
)

type MyCfProjectsPage struct {
	Lang map[string]string
	CfLng string
	CurrencyList map[int64]string
	Projects map[string]map[string]string
	UserId int64
	ProjectsLang map[string]map[string]string
}

type SortMyCfProjects []map[string]string

func (s SortMyCfProjects) Len() int {
	return len(s)
}
func (s SortMyCfProjects) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortMyCfProjects) Less(i, j int) bool {
	return s[i]["name"] < s[j]["name"]
}

func (c *Controller) MyCfProjects() (string, error) {

	var err error
	log.Println("MyCfProjects")

	projectsLang := make(map[string]map[string]string)
	projects := make(map[string]map[string]string)
	cfProjects, err := c.GetAll(`
			SELECT cf_projects.id, lang_id, blurb_img, country, city, currency_id, end_time
			FROM cf_projects
			WHERE user_id = ? AND del_block_id = ?
			`, -1, c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	for _, data := range cfProjects {
		CfProjectData, err := c.GetCfProjectData(utils.StrToInt64(data["id"]), utils.StrToInt64(data["end_time"]), c.LangInt, utils.StrToFloat64(data["end_time"]), "")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		for k, v := range CfProjectData {
			data[k] = v
		}
		projects[data["id"]] = data
		lang, err:=c.GetMap(`SELECT id, lang_id FROM cf_projects_data WHERE project_id = ?`, "id", "lang_id", data["id"])
		projectsLang[data["id"]] = map[string]string{lang}
	}

	cfLng, err := c.GetAllCfLng()

	TemplateStr, err := makeTemplate("cf_catalog", "MyCfProjects", &MyCfProjectsPage{
		Lang: c.Lang,
		CfLng: cfLng,
		CurrencyList: c.CurrencyList,
		Projects: projects,
		UserId: c.SessUserId,
		ProjectsLang: projectsLang})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}


