package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"encoding/json"
)

func CfProjects() {

	const GoroutineName = "CfProjects"
	d := new(daemon)
	d.DCDB = DbConnect()
	if d.DCDB == nil {
		return
	}
	d.goRoutineName = GoroutineName
	if !d.CheckInstall(DaemonCh, AnswerDaemonCh) {
		return
	}
	d.DCDB = DbConnect()
	if d.DCDB == nil {
		return
	}

	BEGIN:
	for {
		log.Info(GoroutineName)
		MonitorDaemonCh <- []string{GoroutineName, utils.Int64ToStr(utils.Time())}

		// проверим, не нужно ли нам выйти из цикла
		if CheckDaemonsRestart() {
			break BEGIN
		}

		err, restart := d.dbLock()
		if restart {
			break BEGIN
		}
		if err != nil {
			d.PrintSleep(err, 1)
			continue BEGIN
		}

		// гео-декодирование
		all, err := d.GetAll(`
				SELECT id,
							latitude,
							longitude
				FROM cf_projects
				WHERE geo_checked= 0
				`, -1)
		for _, cf_projects := range all {
			gmapData, err := utils.GetHttpTextAnswer("http://maps.googleapis.com/maps/api/geocode/json?latlng="+cf_projects["latitude"]+","+cf_projects["longitude"]+"&sensor=true_or_false")
			if err != nil {
				d.PrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			var gmap map[string][]map[string][]map[string]string
			json.Unmarshal([]byte(gmapData), &gmap)
			if len(gmap["results"]) > 1 && len(gmap["results"][len(gmap["results"])-2]["address_components"]) > 0 {
				country := gmap["results"][len(gmap["results"])-2]["address_components"][0]["long_name"]
				city := gmap["results"][len(gmap["results"])-2]["address_components"][1]["short_name"]
				err = d.ExecSql("UPDATE cf_projects SET country = ?, city = ?, geo_checked= 1 WHERE id = ?", country, city, cf_projects["id"])
				if err != nil {
					d.unlockPrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
			}
		}

		// финансирование проектов
		cf_funding, err := d.GetAll(`
				SELECT  id,
							 project_id,
							 amount
				FROM cf_funding
				WHERE checked= 0
				`, -1)
		for _, data := range cf_funding {
			// отмечаем, чтобы больше не брать
			err = d.ExecSql("UPDATE cf_funding SET checked = 1 WHERE id = ?", data["id"])
			if err != nil {
				d.unlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			// сколько собрано средств
			funding, err := d.Single("SELECT sum(amount) FROM cf_funding WHERE project_id  =  ? AND del_block_id  =  0", data["project_id"]).Float64()
			if err != nil {
				d.unlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}

			// сколько всего фундеров
			countFunders, err := d.Single("SELECT count(id) FROM cf_funding WHERE project_id  = ? AND del_block_id  =  0 GROUP BY user_id", data["project_id"]).Int64()
			if err != nil {
				d.unlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}

			// обновляем кол-во фундеров и собранные средства
			err = d.ExecSql("UPDATE cf_projects SET funding = ?, funders = ? WHERE id = ?", funding, countFunders, data["project_id"])
			if err != nil {
				d.unlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
		}

		d.dbUnlock()

		for i:=0; i < 60; i++ {
			utils.Sleep(1)
			// проверим, не нужно ли нам выйти из цикла
			if CheckDaemonsRestart() {
				break BEGIN
			}
		}
	}
}
