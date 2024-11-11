package web

import (
	"context"
)

type session interface {
	Get(string) any
	Set(key string, val interface{})
	Save() error
}

type Store struct {
}

type store interface {
	Get(*context.Context) (session, error)
}

type Presenter struct {
	// store?
	store store
}

func (p *Presenter) GetSession(c *context.Context) session {
	// Get session from storage
	sess, err := p.store.Get(c)
	if err != nil {
		panic(err)
	}
	return sess
}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func (p *Presenter) Render(c *context.Context, statusCode int, data map[string]interface{}, templateName string, redirectHTML bool) error {
	// session := p.GetSession(c)
	// userId, ok := session.Get("user_id").(string)
	// // admin, okAdmin := session.Get("admin").(bool)
	// username := session.Get("username")

	// if data["nav"] == nil {
	// 	data["nav"] = ""
	// }
	// data["isLogin"] = false
	// if ok == true {
	// 	data["isLogin"] = true
	// 	data["isUser"] = userId
	// 	data["username"] = username
	// 	user, _ := authHandler.CurrentUser(c)
	// 	// dont forget to use in template with {{with}}
	// 	data["user"] = user
	// }
	// // queryArgs := map[string]string{}
	// // c.QueryArgs().VisitAll(func(key []byte, value []byte) {
	// // 	queryArgs[string(key)] = string(value)
	// // })
	// // data["Query"] = queryArgs
	// // data["Version"] = handler.Version
	// // if okAdmin == true {
	// // 	data["isAdmin"] = admin
	// // }

	// data["notice"] = session.Get("notice")
	// data["alert"] = session.Get("alert")
	// if redirectHTML == false {
	// 	// after showing alert just reset it
	// 	session.Set("notice", nil)
	// 	session.Set("alert", nil)
	// 	session.Save()
	// }

	// // switch c.GetReqHeaders()["Accept"] {
	// // case "application/json":
	// // 	// Respond with internal.JSON
	// // 	message, _ := data["message"].(string)
	// // 	return p.JSON(c, statusCode, message, data["payload"])
	// // case "application/xml":
	// // 	// Respond with XML
	// // 	// Fiber doesnt have XML response, make a generic one if needed
	// // 	return errors.New("XML Not supported")
	// // default:
	// // 	if redirectHTML == false {
	// // 		return p.HTML(c, statusCode, templateName, data)
	// // 	} else {
	// // 		return c.Redirect(templateName)
	// // 	}
	// // }
	return nil
}

func (p *Presenter) SetNotice(c *context.Context, notice string) error {
	session := p.GetSession(c)
	// Save the username in the session
	session.Set("notice", notice)
	if err := session.Save(); err != nil {
		return err
	}
	return nil
}

func (p *Presenter) SetAlert(c *context.Context, notice string) error {
	session := p.GetSession(c)
	// Save the username in the session
	session.Set("alert", notice)
	if err := session.Save(); err != nil {
		return err
	}
	return nil
}
