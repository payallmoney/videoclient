package main

import (
	"net/http"
	"io/ioutil"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"bytes"
	"strings"
)

func clearPlayList() interface{} {
	params := bson.M{"jsonrpc": "2.0", "method": "Playlist.Clear", "params": bson.M{"playlistid": 1}, "id": 1}
	return jsonrpc(params, false)
}

func getPlayeditem() string {
	params := bson.M{"jsonrpc": "2.0", "method": "Player.GetItem", "params": bson.M{"playerid": 1}, "id": 1}
	ret := jsonrpc(params, false);
	result, _ := ret.(map[string]interface{})["result"].(map[string]interface{})["item"].(map[string]interface{})["label"].(string)
	return result
}
func activePlayer() bool {
	params := bson.M{"jsonrpc": "2.0", "method": "Player.GetActivePlayers", "id": 1}
	ret := jsonrpc(params, false).(map[string]interface{})["result"].([]interface{});
	if (len(ret) > 0) {
		return true
	}else {
		return false
	}
}
func playPause() interface{} {
	params := bson.M{"jsonrpc": "2.0", "method": "Player.PlayPause", "id": 1, "params":  bson.M{
		"playerid": 1 }}
	return jsonrpc(params, false)
}
func play() interface{} {
	params := bson.M{"jsonrpc": "2.0", "method": "Player.Open", "id": 1, "params":  bson.M{
		"item":bson.M{"playlistid":1},
		"options":bson.M{"repeat":"all"}}}
	return jsonrpc(params, false)
}
func playVideo(video string) interface{} {
	video = strings.Replace(video, "\\", "/", -1)
	params := bson.M{"jsonrpc": "2.0", "method": "Player.Open", "id": 1, "params":  bson.M{
		"item":bson.M{"file":video} }}
	return jsonrpc(params, false)
}
func additem(file string) interface{} {
	params := bson.M{
		"jsonrpc": "2.0",
		"method": "Playlist.Add",
		"id": 1,
		"params": bson.M{
			"playlistid": 1,
			"item":bson.M{"file":file}}}
	return jsonrpc(params, false)
}

func jsonrpc(param interface{}, flag bool) interface{} {
	flag = true
	params_str, err := json.Marshal(param)
	checkerr(err)
	if (flag) {
		log_print(string(params_str))
	}
	resp, err := http.Post(KodiUrl("/jsonrpc"), "application/json", bytes.NewBuffer([]byte(params_str)))
	checkerr(err)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if (flag) {
		log_printf("=====ret====%s", string(body))
	}
	var result  interface{}
	json.Unmarshal(body, &result)
	return result
}

func getlist() interface{} {
	params := bson.M{
		"jsonrpc": "2.0",
		"method": "Playlist.GetItems",
		"params": bson.M{"properties": []string{"runtime", "showtitle", "season", "title", "artist" },
			"playlistid": 1},
		"id": 1}

	return jsonrpc(params, true)
}
func getPlayer(id int) {
	params := bson.M{"jsonrpc": "2.0", "method": "Playlist.Clear", "params": bson.M{"playlistid": id}, "id": 1}
	jsonrpc(params, false)
}