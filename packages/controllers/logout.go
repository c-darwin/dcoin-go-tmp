package controllers
import (

)

func (c *Controller) Logout() (string, error) {
	c.sess.Delete("user_id")
	c.sess.Delete("public_key")
	c.sess.Delete("private_key")
	err := c.ExecSql(`DELETE private_key FROM `+c.MyPrefix+`my_keys WHERE block_id = (SELECT max(block_id) FROM `+c.MyPrefix+`my_keys)`)
	if err != nil {
		return "", err
	}
	return "", nil
}


