# Wrapped Gorilla Session
将 gorilla session https://github.com/gorilla/sessions 转化为 externalauth库可以使用的Session 对象

## 范例代码

### 与External auth 一起使用
    gstore := sessions.NewCookieStore([]byte("something-very-secret"))
    Session := New(gstore, "test")
	mux := http.NewServeMux()
	mux.Handle(Auth.Path+"/", http.StripPrefix(Auth.Path, WrapRecover(
		Session.Wrap(
			http.HandlerFunc(
				Auth.Serve(ActionSuccess),
			),
		),
	)),
	)

### 作为session正常使用

	gstore := sessions.NewCookieStore([]byte("something-very-secret"))
	session := New(gstore, "test")

    value := &data{}
	err := session.Set(r, "test", value)
    
    Value2:=&data{}

    err := session.Get(r, "test", value)