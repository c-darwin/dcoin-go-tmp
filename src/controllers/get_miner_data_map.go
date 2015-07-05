package controllers
import (
	"utils"
)

func (c *Controller) GetMinerDataMap() (string, error) {

	rows, err := c.Query(c.FormatQuery("SELECT user_id, latitude, longitude FROM miners_data WHERE status  =  'miner' AND user_id>7 AND user_id != 106"))
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	defer rows.Close()
	result:=""
	for rows.Next() {
		var user_id, latitude, longitude string
		err = rows.Scan(&user_id, &latitude, &longitude)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		result+="{\"user_id\": "+user_id+",\"longitude\": "+longitude+", \"latitude\": "+latitude+"},";
	}
	c.w.Header().Set("Access-Control-Allow-Origin", "*")
	return string(`{ "info": [`+result[:len(result)-1]+`}`), nil
}
