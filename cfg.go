package main

func Cfg() map[string]interface{} {
	ret := make(map[string]interface{})
	ret["server"] = "127.0.0.1:3000"
	return ret;
}
