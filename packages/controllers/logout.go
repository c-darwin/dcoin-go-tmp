package controllers
import (

)

func (c *Controller) Logout() (string, error) {
	c.sess.Delete("user_id")
	c.sess.Delete("public_key")
	c.sess.Delete("private_key")
	return "", nil
}


