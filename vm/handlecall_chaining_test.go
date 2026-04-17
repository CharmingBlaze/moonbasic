package vm

import (
	"testing"

	"moonbasic/vm/heap"
)

// TestHandleCallBuiltinModelSettersDispatch documents registry keys for TagModel handle methods.
func TestHandleCallBuiltinModelSettersDispatch(t *testing.T) {
	core := []struct {
		method string
		want   string
	}{
		{"SETPOS", "MODEL.SETPOS"},
		{"SETROT", "MODEL.SETROT"},
		{"SETSCALE", "MODEL.SETSCALE"},
		{"SETCOLOR", "MODEL.SETCOLOR"},
		{"SETALPHA", "MODEL.SETALPHA"},
		{"DRAW", "MODEL.DRAW"},
		{"FREE", "MODEL.FREE"},
		{"MOVE", "MODEL.MOVE"},
		{"SETCULL", "MODEL.SETCULL"},
		{"SETWIREFRAME", "MODEL.SETWIREFRAME"},
		{"SETTEXTURESTAGE", "MODEL.SETTEXTURESTAGE"},
		{"SETMATERIAL", "MODEL.SETMATERIAL"},
		{"SETLODDISTANCES", "MODEL.SETLODDISTANCES"},
		{"ATTACHTO", "MODEL.ATTACHTO"},
	}
	for _, tc := range core {
		k, prep, ok := handleCallBuiltin(heap.TagModel, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagModel %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

// Short names normalize to SET* (e.g. .cull() -> SETCULL) via normalizeHandleMethod.
func TestHandleCallBuiltinModelNormalizedAliases(t *testing.T) {
	tests := []struct {
		method string
		want   string
	}{
		{"cull", "MODEL.SETCULL"},
		{"wireframe", "MODEL.SETWIREFRAME"},
		{"fog", "MODEL.SETFOG"},
		{"metal", "MODEL.SETMETAL"},
		{"stagescale", "MODEL.SETSTAGESCALE"},
		{"attach", "MODEL.ATTACHTO"},
	}
	for _, tc := range tests {
		k, prep, ok := handleCallBuiltin(heap.TagModel, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagModel %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallBuiltinTagLODModelMatchesModel(t *testing.T) {
	k1, p1, ok1 := handleCallBuiltin(heap.TagModel, "SETDIFFUSE")
	k2, p2, ok2 := handleCallBuiltin(heap.TagLODModel, "SETDIFFUSE")
	if !ok1 || !ok2 || k1 != k2 || p1 != p2 || k1 != "MODEL.SETDIFFUSE" {
		t.Fatalf("TagModel vs TagLODModel: (%q,%v,%v) vs (%q,%v,%v)", k1, p1, ok1, k2, p2, ok2)
	}
}

func TestHandleCallBuiltinMaterialDispatch(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"SETSHADER", "MATERIAL.SETSHADER"},
		{"SETTEXTURE", "MATERIAL.SETTEXTURE"},
		{"SETCOLOR", "MATERIAL.SETCOLOR"},
		{"SETFLOAT", "MATERIAL.SETFLOAT"},
		{"SETEFFECT", "MATERIAL.SETEFFECT"},
		{"shader", "MATERIAL.SETSHADER"},
	}
	for _, tc := range cases {
		k, prep, ok := handleCallBuiltin(heap.TagMaterial, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagMaterial %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallRegistryPrefixMaterial(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagMaterial); p != "MATERIAL." {
		t.Fatalf("TagMaterial prefix: got %q", p)
	}
}

func TestHandleCallBuiltinTerrainOps(t *testing.T) {
	tests := []struct {
		method string
		want   string
	}{
		{"APPLYMAP", "TERRAIN.APPLYMAP"},
		{"DRAW", "TERRAIN.DRAW"},
		{"FILLPERLIN", "TERRAIN.FILLPERLIN"},
		{"SETDETAIL", "TERRAIN.SETDETAIL"},
		{"applymap", "TERRAIN.APPLYMAP"},
	}
	for _, tc := range tests {
		k, prep, ok := handleCallBuiltin(heap.TagTerrain, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagTerrain %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallDispatchTerrainDetailGetter(t *testing.T) {
	k, prep, ok := handleCallDispatch(heap.TagTerrain, "detail", 0)
	if !ok || !prep || k != "TERRAIN.GETDETAIL" {
		t.Fatalf("terrain detail() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
}

func TestHandleCallRegistryPrefixWater(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagWater); p != "WATER." {
		t.Fatalf("TagWater prefix: got %q", p)
	}
}

func TestHandleCallBuiltinWater(t *testing.T) {
	tests := []struct {
		method string
		want   string
	}{
		{"DRAW", "WATER.DRAW"},
		{"SETPOS", "WATER.SETPOS"},
		{"wave", "WATER.SETWAVE"},
		{"GETWAVESPEED", "WATER.GETWAVESPEED"},
		{"SETSHALLOWCOLOR", "WATER.SETSHALLOWCOLOR"},
	}
	for _, tc := range tests {
		k, prep, ok := handleCallBuiltin(heap.TagWater, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagWater %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallRegistryPrefixShader(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagShader); p != "SHADER." {
		t.Fatalf("TagShader prefix: got %q", p)
	}
}

func TestHandleCallBuiltinShader(t *testing.T) {
	tests := []struct {
		method string
		want   string
	}{
		{"SETFLOAT", "SHADER.SETFLOAT"},
		{"setVec3", "SHADER.SETVEC3"},
		{"SETVECTOR", "SHADER.SETVEC3"},
		{"GETLOC", "SHADER.GETLOC"},
		{"SETTEXTURE", "SHADER.SETTEXTURE"},
		{"FREE", "SHADER.FREE"},
	}
	for _, tc := range tests {
		k, prep, ok := handleCallBuiltin(heap.TagShader, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagShader %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallRegistryPrefixSky(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagSky); p != "SKY." {
		t.Fatalf("TagSky prefix: got %q", p)
	}
}

func TestHandleCallBuiltinSky(t *testing.T) {
	tests := []struct {
		method string
		want   string
	}{
		{"UPDATE", "SKY.UPDATE"},
		{"DRAW", "SKY.DRAW"},
		{"SETTIME", "SKY.SETTIME"},
		{"SETDAYLENGTH", "SKY.SETDAYLENGTH"},
		{"FREE", "SKY.FREE"},
	}
	for _, tc := range tests {
		k, prep, ok := handleCallBuiltin(heap.TagSky, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagSky %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallDispatchSkyTimeGetter(t *testing.T) {
	k, prep, ok := handleCallDispatch(heap.TagSky, "time", 0)
	if !ok || !prep || k != "SKY.GETTIMEHOURS" {
		t.Fatalf("sky time() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
}

func TestHandleCallRegistryPrefixDecal(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagDecal); p != "DECAL." {
		t.Fatalf("TagDecal prefix: got %q", p)
	}
}

func TestHandleCallBuiltinDecal(t *testing.T) {
	tests := []struct {
		method string
		want   string
	}{
		{"SETPOS", "DECAL.SETPOS"},
		{"DRAW", "DECAL.DRAW"},
		{"GETPOS", "DECAL.GETPOS"},
		{"SETLIFETIME", "DECAL.SETLIFETIME"},
	}
	for _, tc := range tests {
		k, prep, ok := handleCallBuiltin(heap.TagDecal, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagDecal %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallDispatchDecalPosGetter(t *testing.T) {
	k, prep, ok := handleCallDispatch(heap.TagDecal, "pos", 0)
	if !ok || !prep || k != "DECAL.GETPOS" {
		t.Fatalf("decal pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
}

func TestHandleCallRegistryPrefixCloud(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagCloud); p != "CLOUD." {
		t.Fatalf("TagCloud prefix: got %q", p)
	}
}

func TestHandleCallBuiltinCloud(t *testing.T) {
	tests := []struct {
		method string
		want   string
	}{
		{"UPDATE", "CLOUD.UPDATE"},
		{"DRAW", "CLOUD.DRAW"},
		{"SETCOVERAGE", "CLOUD.SETCOVERAGE"},
		{"coverage", "CLOUD.SETCOVERAGE"},
		{"FREE", "CLOUD.FREE"},
	}
	for _, tc := range tests {
		k, prep, ok := handleCallBuiltin(heap.TagCloud, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagCloud %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallRegistryPrefixWeather(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagWeather); p != "WEATHER." {
		t.Fatalf("TagWeather prefix: got %q", p)
	}
}

func TestHandleCallBuiltinWeather(t *testing.T) {
	tests := []struct {
		method string
		want   string
	}{
		{"UPDATE", "WEATHER.UPDATE"},
		{"DRAW", "WEATHER.DRAW"},
		{"SETTYPE", "WEATHER.SETTYPE"},
		{"type", "WEATHER.SETTYPE"},
		{"GETCOVERAGE", "WEATHER.GETCOVERAGE"},
		{"GETTYPE", "WEATHER.GETTYPE"},
		{"FREE", "WEATHER.FREE"},
	}
	for _, tc := range tests {
		k, prep, ok := handleCallBuiltin(heap.TagWeather, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagWeather %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallDispatchWeatherGetters(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"type", "WEATHER.GETTYPE"},
		{"coverage", "WEATHER.GETCOVERAGE"},
	}
	for _, tc := range cases {
		k, prep, ok := handleCallDispatch(heap.TagWeather, tc.method, 0)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagWeather %q (0 args): got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallRegistryPrefixScatterProp(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagScatterSet); p != "SCATTER." {
		t.Fatalf("TagScatterSet prefix: got %q", p)
	}
	if p := handleCallRegistryPrefix(heap.TagProp); p != "PROP." {
		t.Fatalf("TagProp prefix: got %q", p)
	}
}

func TestHandleCallBuiltinScatterProp(t *testing.T) {
	scatter := []struct {
		method string
		want   string
	}{
		{"APPLY", "SCATTER.APPLY"},
		{"DRAWALL", "SCATTER.DRAWALL"},
		{"FREE", "SCATTER.FREE"},
	}
	for _, tc := range scatter {
		k, prep, ok := handleCallBuiltin(heap.TagScatterSet, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagScatterSet %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
	k, prep, ok := handleCallBuiltin(heap.TagProp, "FREE")
	if !ok || !prep || k != "PROP.FREE" {
		t.Fatalf("TagProp FREE: got %q prepend=%v ok=%v", k, prep, ok)
	}
	kp, prepP, okP := handleCallBuiltin(heap.TagProp, "PLACE")
	if !okP || !prepP || kp != "PROP.PLACE" {
		t.Fatalf("TagProp PLACE: got %q prepend=%v ok=%v", kp, prepP, okP)
	}
}

func TestHandleCallBuiltinLobbyFindNotReceiverFirst(t *testing.T) {
	_, _, ok := handleCallBuiltin(heap.TagLobby, "FIND")
	if ok {
		t.Fatalf("LOBBY.FIND is not a handle method; expected no mapping")
	}
}

func TestHandleCallBuiltinTweenUpdate(t *testing.T) {
	k, prep, ok := handleCallBuiltin(heap.TagTween, "UPDATE")
	if !ok || !prep || k != "TWEEN.UPDATE" {
		t.Fatalf("TagTween UPDATE: got %q prepend=%v ok=%v", k, prep, ok)
	}
}

func TestHandleCallBuiltinMusicUpdate(t *testing.T) {
	km, pm, okm := handleCallBuiltin(heap.TagMusic, "UPDATE")
	if !okm || !pm || km != "AUDIO.UPDATEMUSIC" {
		t.Fatalf("TagMusic UPDATE: got %q prepend=%v ok=%v", km, pm, okm)
	}
	_, _, oks := handleCallBuiltin(heap.TagSound, "UPDATE")
	if oks {
		t.Fatalf("TagSound UPDATE: expected no mapping (use stream/sound APIs as documented)")
	}
}

func TestHandleCallBuiltinInstancedModelExplicitRegistryNames(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"SETINSTANCEPOS", "INSTANCE.SETINSTANCEPOS"},
		{"SETINSTANCESCALE", "INSTANCE.SETINSTANCESCALE"},
		{"SETPOS", "INSTANCE.SETPOS"},
		{"GETALPHA", "INSTANCE.GETALPHA"},
	}
	for _, tc := range cases {
		k, prep, ok := handleCallBuiltin(heap.TagInstancedModel, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagInstancedModel %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallDispatchInstancedModelColorAlphaGetters(t *testing.T) {
	tests := []struct {
		method string
		want   string
	}{
		{"col", "INSTANCE.GETCOLOR"},
		{"alpha", "INSTANCE.GETALPHA"},
		{"getColor", "INSTANCE.GETCOLOR"},
		{"getAlpha", "INSTANCE.GETALPHA"},
	}
	for _, tc := range tests {
		k, prep, ok := handleCallDispatch(heap.TagInstancedModel, tc.method, 0)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagInstancedModel %q 0-arg: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallRegistryPrefixTweenBiomeNoise(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagTween); p != "TWEEN." {
		t.Fatalf("TagTween prefix: got %q", p)
	}
	if p := handleCallRegistryPrefix(heap.TagBiome); p != "BIOME." {
		t.Fatalf("TagBiome prefix: got %q", p)
	}
	if p := handleCallRegistryPrefix(heap.TagNoise); p != "NOISE." {
		t.Fatalf("TagNoise prefix: got %q", p)
	}
}

func TestHandleCallBuiltinBiomeNoise(t *testing.T) {
	biomeCases := []struct {
		method string
		want   string
	}{
		{"SETTEMP", "BIOME.SETTEMP"},
		{"temp", "BIOME.SETTEMP"},
		{"SETHUMIDITY", "BIOME.SETHUMIDITY"},
		{"humidity", "BIOME.SETHUMIDITY"},
		{"FREE", "BIOME.FREE"},
	}
	for _, tc := range biomeCases {
		k, prep, ok := handleCallBuiltin(heap.TagBiome, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagBiome %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
	noiseCases := []struct {
		method string
		want   string
	}{
		{"SETFREQUENCY", "NOISE.SETFREQUENCY"},
		{"SETDOMAINWARPAMPLITUDE", "NOISE.SETDOMAINWARPAMPLITUDE"},
		{"GET", "NOISE.GET"},
		{"FILLARRAY", "NOISE.FILLARRAY"},
		{"FREE", "NOISE.FREE"},
	}
	for _, tc := range noiseCases {
		k, prep, ok := handleCallBuiltin(heap.TagNoise, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagNoise %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallRegistryPrefixTablePool(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagTable); p != "TABLE." {
		t.Fatalf("TagTable prefix: got %q", p)
	}
	if p := handleCallRegistryPrefix(heap.TagPool); p != "POOL." {
		t.Fatalf("TagPool prefix: got %q", p)
	}
}

func TestHandleCallBuiltinTablePool(t *testing.T) {
	tableCases := []struct {
		method string
		want   string
	}{
		{"ADDROW", "TABLE.ADDROW"},
		{"GET", "TABLE.GET"},
		{"SET", "TABLE.SET"},
		{"ROWCOUNT", "TABLE.ROWCOUNT"},
		{"rows", "TABLE.ROWCOUNT"},
		{"FREE", "TABLE.FREE"},
	}
	for _, tc := range tableCases {
		k, prep, ok := handleCallBuiltin(heap.TagTable, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagTable %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
	poolCases := []struct {
		method string
		want   string
	}{
		{"SETFACTORY", "POOL.SETFACTORY"},
		{"GET", "POOL.GET"},
		{"RETURN", "POOL.RETURN"},
		{"FREE", "POOL.FREE"},
	}
	for _, tc := range poolCases {
		k, prep, ok := handleCallBuiltin(heap.TagPool, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagPool %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallDispatchTableRowColCount(t *testing.T) {
	k, prep, ok := handleCallDispatch(heap.TagTable, "rowCount", 0)
	if !ok || !prep || k != "TABLE.ROWCOUNT" {
		t.Fatalf("table rowCount(): got %q prepend=%v ok=%v", k, prep, ok)
	}
	k2, prep2, ok2 := handleCallDispatch(heap.TagTable, "cols", 0)
	if !ok2 || !prep2 || k2 != "TABLE.COLCOUNT" {
		t.Fatalf("table cols(): got %q prepend=%v ok=%v", k2, prep2, ok2)
	}
}

func TestHandleCallRegistryPrefixJSONCSV(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagJSON); p != "JSON." {
		t.Fatalf("TagJSON prefix: got %q", p)
	}
	if p := handleCallRegistryPrefix(heap.TagCSV); p != "CSV." {
		t.Fatalf("TagCSV prefix: got %q", p)
	}
}

func TestHandleCallBuiltinJSONCSV(t *testing.T) {
	jk, jprep, jok := handleCallBuiltin(heap.TagJSON, "SETSTRING")
	if !jok || !jprep || jk != "JSON.SETSTRING" {
		t.Fatalf("TagJSON SETSTRING: got %q prepend=%v ok=%v", jk, jprep, jok)
	}
	ck, cprep, cok := handleCallBuiltin(heap.TagCSV, "SET")
	if !cok || !cprep || ck != "CSV.SET" {
		t.Fatalf("TagCSV SET: got %q prepend=%v ok=%v", ck, cprep, cok)
	}
}

func TestHandleCallDispatchJSONCSVLenRow(t *testing.T) {
	k, prep, ok := handleCallDispatch(heap.TagJSON, "len", 0)
	if !ok || !prep || k != "JSON.LEN" {
		t.Fatalf("json len(): got %q prepend=%v ok=%v", k, prep, ok)
	}
	k2, prep2, ok2 := handleCallDispatch(heap.TagCSV, "rowCount", 0)
	if !ok2 || !prep2 || k2 != "CSV.ROWCOUNT" {
		t.Fatalf("csv rowCount(): got %q prepend=%v ok=%v", k2, prep2, ok2)
	}
}

func TestHandleCallRegistryPrefixDBRng(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagDB); p != "DB." {
		t.Fatalf("TagDB prefix: got %q", p)
	}
	if p := handleCallRegistryPrefix(heap.TagDBRows); p != "ROWS." {
		t.Fatalf("TagDBRows prefix: got %q", p)
	}
	if p := handleCallRegistryPrefix(heap.TagDBStmt); p != "DB." {
		t.Fatalf("TagDBStmt prefix: got %q", p)
	}
	if p := handleCallRegistryPrefix(heap.TagDBTx); p != "DB." {
		t.Fatalf("TagDBTx prefix: got %q", p)
	}
	if p := handleCallRegistryPrefix(heap.TagRng); p != "RAND." {
		t.Fatalf("TagRng prefix: got %q", p)
	}
}

func TestHandleCallBuiltinDBRng(t *testing.T) {
	k, prep, ok := handleCallBuiltin(heap.TagDB, "QUERY")
	if !ok || !prep || k != "DB.QUERY" {
		t.Fatalf("TagDB QUERY: got %q prepend=%v ok=%v", k, prep, ok)
	}
	kr, prepr, okr := handleCallBuiltin(heap.TagDBRows, "NEXT")
	if !okr || !prepr || kr != "ROWS.NEXT" {
		t.Fatalf("TagDBRows NEXT: got %q prepend=%v ok=%v", kr, prepr, okr)
	}
	ks, preps, oks := handleCallBuiltin(heap.TagDBStmt, "STMTEXEC")
	if !oks || !preps || ks != "DB.STMTEXEC" {
		t.Fatalf("TagDBStmt STMTEXEC: got %q prepend=%v ok=%v", ks, preps, oks)
	}
	kt, prept, okt := handleCallBuiltin(heap.TagDBTx, "COMMIT")
	if !okt || !prept || kt != "DB.COMMIT" {
		t.Fatalf("TagDBTx COMMIT: got %q prepend=%v ok=%v", kt, prept, okt)
	}
	kn, prepn, okn := handleCallBuiltin(heap.TagRng, "NEXT")
	if !okn || !prepn || kn != "RAND.NEXT" {
		t.Fatalf("TagRng NEXT: got %q prepend=%v ok=%v", kn, prepn, okn)
	}
}

func TestHandleCallDispatchDBIsOpen(t *testing.T) {
	k, prep, ok := handleCallDispatch(heap.TagDB, "isOpen", 0)
	if !ok || !prep || k != "DB.ISOPEN" {
		t.Fatalf("db isOpen(): got %q prepend=%v ok=%v", k, prep, ok)
	}
}

func TestHandleCallRegistryPrefixMemLobbyPacket(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagMem); p != "MEM." {
		t.Fatalf("TagMem prefix: got %q", p)
	}
	if p := handleCallRegistryPrefix(heap.TagLobby); p != "LOBBY." {
		t.Fatalf("TagLobby prefix: got %q", p)
	}
	if p := handleCallRegistryPrefix(heap.TagNetPacket); p != "PACKET." {
		t.Fatalf("TagNetPacket prefix: got %q", p)
	}
}

func TestHandleCallBuiltinMemLobbyPacket(t *testing.T) {
	k, prep, ok := handleCallBuiltin(heap.TagMem, "GETBYTE")
	if !ok || !prep || k != "MEM.GETBYTE" {
		t.Fatalf("TagMem GETBYTE: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k2, prep2, ok2 := handleCallBuiltin(heap.TagLobby, "SETPROPERTY")
	if !ok2 || !prep2 || k2 != "LOBBY.SETPROPERTY" {
		t.Fatalf("TagLobby SETPROPERTY: got %q prepend=%v ok=%v", k2, prep2, ok2)
	}
	k3, prep3, ok3 := handleCallBuiltin(heap.TagNetPacket, "DATA")
	if !ok3 || !prep3 || k3 != "PACKET.DATA" {
		t.Fatalf("TagNetPacket DATA: got %q prepend=%v ok=%v", k3, prep3, ok3)
	}
}

func TestHandleCallDispatchMemLobbyPacket(t *testing.T) {
	k, prep, ok := handleCallDispatch(heap.TagMem, "size", 0)
	if !ok || !prep || k != "MEM.SIZE" {
		t.Fatalf("mem size(): got %q prepend=%v ok=%v", k, prep, ok)
	}
	k2, prep2, ok2 := handleCallDispatch(heap.TagLobby, "name", 0)
	if !ok2 || !prep2 || k2 != "LOBBY.GETNAME" {
		t.Fatalf("lobby name(): got %q prepend=%v ok=%v", k2, prep2, ok2)
	}
	k3, prep3, ok3 := handleCallDispatch(heap.TagNetPacket, "data", 0)
	if !ok3 || !prep3 || k3 != "PACKET.DATA" {
		t.Fatalf("packet data(): got %q prepend=%v ok=%v", k3, prep3, ok3)
	}
}

func TestHandleCallRegistryPrefixNetGameHandles(t *testing.T) {
	tests := []struct {
		tag  uint16
		want string
	}{
		{heap.TagHost, "NET."},
		{heap.TagEvent, "EVENT."},
		{heap.TagPlayer2D, "PLAYER2D."},
		{heap.TagGameTimer, "TIMER."},
		{heap.TagGameTimerSim, "TIMER."},
		{heap.TagGameStopwatch, "STOPWATCH."},
	}
	for _, tc := range tests {
		if p := handleCallRegistryPrefix(tc.tag); p != tc.want {
			t.Fatalf("tag %d: got %q want %q", tc.tag, p, tc.want)
		}
	}
}

func TestHandleCallBuiltinNetEventPlayer2DTimers(t *testing.T) {
	k, prep, ok := handleCallBuiltin(heap.TagHost, "UPDATE")
	if !ok || !prep || k != "NET.UPDATE" {
		t.Fatalf("TagHost UPDATE: got %q prepend=%v ok=%v", k, prep, ok)
	}
	kf, prepf, okf := handleCallBuiltin(heap.TagHost, "FLUSH")
	if !okf || !prepf || kf != "NET.FLUSH" {
		t.Fatalf("TagHost FLUSH: got %q prepend=%v ok=%v", kf, prepf, okf)
	}
	ke, prepe, oke := handleCallBuiltin(heap.TagEvent, "TYPE")
	if !oke || !prepe || ke != "EVENT.TYPE" {
		t.Fatalf("TagEvent TYPE: got %q prepend=%v ok=%v", ke, prepe, oke)
	}
	kp, prepk, okp := handleCallBuiltin(heap.TagPlayer2D, "MOVE")
	if !okp || !prepk || kp != "PLAYER2D.MOVE" {
		t.Fatalf("TagPlayer2D MOVE: got %q prepend=%v ok=%v", kp, prepk, okp)
	}
	kt, prept, okt := handleCallBuiltin(heap.TagGameTimer, "FINISHED")
	if !okt || !prept || kt != "TIMER.FINISHED" {
		t.Fatalf("TagGameTimer FINISHED: got %q prepend=%v ok=%v", kt, prept, okt)
	}
	ks, preps, oks := handleCallBuiltin(heap.TagGameTimerSim, "DONE")
	if !oks || !preps || ks != "TIMER.DONE" {
		t.Fatalf("TagGameTimerSim DONE: got %q prepend=%v ok=%v", ks, preps, oks)
	}
	kw, prepw, okw := handleCallBuiltin(heap.TagGameStopwatch, "ELAPSED")
	if !okw || !prepw || kw != "STOPWATCH.ELAPSED" {
		t.Fatalf("TagGameStopwatch ELAPSED: got %q prepend=%v ok=%v", kw, prepw, okw)
	}
}

func TestHandleCallDispatchNetEventPlayer2DTimers(t *testing.T) {
	k, prep, ok := handleCallDispatch(heap.TagHost, "peerCount", 0)
	if !ok || !prep || k != "NET.PEERCOUNT" {
		t.Fatalf("host peerCount(): got %q prepend=%v ok=%v", k, prep, ok)
	}
	ke, prepe, oke := handleCallDispatch(heap.TagEvent, "data", 0)
	if !oke || !prepe || ke != "EVENT.DATA" {
		t.Fatalf("event data(): got %q prepend=%v ok=%v", ke, prepe, oke)
	}
	kp, prepk, okp := handleCallDispatch(heap.TagPlayer2D, "x", 0)
	if !okp || !prepk || kp != "PLAYER2D.GETX" {
		t.Fatalf("player2d x(): got %q prepend=%v ok=%v", kp, prepk, okp)
	}
	kt, prept, okt := handleCallDispatch(heap.TagGameTimer, "finished", 0)
	if !okt || !prept || kt != "TIMER.FINISHED" {
		t.Fatalf("game timer finished(): got %q prepend=%v ok=%v", kt, prept, okt)
	}
	ks, preps, oks := handleCallDispatch(heap.TagGameTimerSim, "done", 0)
	if !oks || !preps || ks != "TIMER.DONE" {
		t.Fatalf("sim timer done(): got %q prepend=%v ok=%v", ks, preps, oks)
	}
	kw, prepw, okw := handleCallDispatch(heap.TagGameStopwatch, "elapsed", 0)
	if !okw || !prepw || kw != "STOPWATCH.ELAPSED" {
		t.Fatalf("stopwatch elapsed(): got %q prepend=%v ok=%v", kw, prepw, okw)
	}
}

func TestHandleCallRegistryPrefixSteerComputeJointBrushBTree(t *testing.T) {
	tests := []struct {
		tag  uint16
		want string
	}{
		{heap.TagBTree, "BTREE."},
		{heap.TagSteerGroup, "STEER."},
		{heap.TagComputeShader, "COMPUTESHADER."},
		{heap.TagShaderBuffer, "COMPUTESHADER."},
		{heap.TagJoint2D, "JOINT2D."},
	}
	for _, tc := range tests {
		if p := handleCallRegistryPrefix(tc.tag); p != tc.want {
			t.Fatalf("tag %d: got %q want %q", tc.tag, p, tc.want)
		}
	}
	if p := handleCallRegistryPrefix(heap.TagBrush); p != "" {
		t.Fatalf("TagBrush prefix: got %q want empty (PascalCase API keys)", p)
	}
}

func TestHandleCallBuiltinSteerComputeJointBrush(t *testing.T) {
	k, prep, ok := handleCallBuiltin(heap.TagSteerGroup, "clear")
	if !ok || !prep || k != "STEER.GROUPCLEAR" {
		t.Fatalf("SteerGroup clear: got %q prepend=%v ok=%v", k, prep, ok)
	}
	kc, prepc, okc := handleCallBuiltin(heap.TagComputeShader, "DISPATCH")
	if !okc || !prepc || kc != "COMPUTESHADER.DISPATCH" {
		t.Fatalf("ComputeShader DISPATCH: got %q prepend=%v ok=%v", kc, prepc, okc)
	}
	kb, prepb, okb := handleCallBuiltin(heap.TagShaderBuffer, "free")
	if !okb || !prepb || kb != "COMPUTESHADER.BUFFERFREE" {
		t.Fatalf("ShaderBuffer free: got %q prepend=%v ok=%v", kb, prepb, okb)
	}
	kj, prepj, okj := handleCallBuiltin(heap.TagJoint2D, "FREE")
	if !okj || !prepj || kj != "JOINT2D.FREE" {
		t.Fatalf("Joint2D FREE: got %q prepend=%v ok=%v", kj, prepj, okj)
	}
	kbr, prepbr, okbr := handleCallBuiltin(heap.TagBrush, "color")
	if !okbr || !prepbr || kbr != "BrushColor" {
		t.Fatalf("Brush color: got %q prepend=%v ok=%v", kbr, prepbr, okbr)
	}
}

func TestHandleCallPathPrefixBuiltinDispatch(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagPath); p != "PATH." {
		t.Fatalf("TagPath prefix: got %q want PATH.", p)
	}
	kb, preb, okb := handleCallBuiltin(heap.TagPath, "NODEX")
	if !okb || !preb || kb != "PATH.NODEX" {
		t.Fatalf("path nodeX: got %q prepend=%v ok=%v", kb, preb, okb)
	}
	kv, prev, okv := handleCallDispatch(heap.TagPath, "isValid", 0)
	if !okv || !prev || kv != "PATH.ISVALID" {
		t.Fatalf("path isValid(): got %q prepend=%v ok=%v", kv, prev, okv)
	}
	kc, prec, okc := handleCallDispatch(heap.TagPath, "nodeCount", 0)
	if !okc || !prec || kc != "PATH.NODECOUNT" {
		t.Fatalf("path nodeCount(): got %q prepend=%v ok=%v", kc, prec, okc)
	}
}

func TestHandleCallRegistryPrefixTerrainDrawMoverAudioSplit(t *testing.T) {
	tests := []struct {
		tag  uint16
		want string
	}{
		{heap.TagTerrain, "TERRAIN."},
		{heap.TagDrawPrim3D, "DRAWPRIM3D."},
		{heap.TagDrawPrim2D, "DRAWPRIM2D."},
		{heap.TagTextDraw, "TEXTDRAW."},
		{heap.TagTextDrawEx, "TEXTEXOBJ."},
		{heap.TagTextureDraw, "DRAWTEX"},
		{heap.TagMoverFacade, "MOVER."},
		{heap.TagTacticalGrid, "GRID."},
		{heap.TagAudioStream, "AUDIOSTREAM."},
		{heap.TagWave, "WAVE."},
	}
	for _, tc := range tests {
		if p := handleCallRegistryPrefix(tc.tag); p != tc.want {
			t.Fatalf("tag %d: got %q want %q", tc.tag, p, tc.want)
		}
	}
	if p := handleCallRegistryPrefix(heap.TagSound); p != "AUDIO." {
		t.Fatalf("TagSound: got %q", p)
	}
	if p := handleCallRegistryPrefix(heap.TagImageSequence); p != "IMAGE." {
		t.Fatalf("TagImageSequence: got %q want IMAGE.", p)
	}
}

func TestHandleCallBuiltinGridMeshBuilderImageSeq(t *testing.T) {
	kg, prepg, okg := handleCallBuiltin(heap.TagTacticalGrid, "SETCELL")
	if !okg || !prepg || kg != "GRID.SETCELL" {
		t.Fatalf("grid setCell: got %q prepend=%v ok=%v", kg, prepg, okg)
	}
	km, prepm, okm := handleCallBuiltin(heap.TagMeshBuilder, "ADDVERTEX")
	if !okm || !prepm || km != "ENTITY.ADDVERTEX" {
		t.Fatalf("mesh builder addVertex: got %q prepend=%v ok=%v", km, prepm, okm)
	}
	kx, prepx, okx := handleCallBuiltin(heap.TagMeshBuilder, "x")
	if !okx || !prepx || kx != "ENTITY.VERTEXX" {
		t.Fatalf("mesh builder x: got %q prepend=%v ok=%v", kx, prepx, okx)
	}
	ks, preps, oks := handleCallBuiltin(heap.TagImageSequence, "FREE")
	if !oks || !preps || ks != "IMAGE.FREE" {
		t.Fatalf("image sequence FREE: got %q prepend=%v ok=%v", ks, preps, oks)
	}
}

func TestHandleCallRegistryPrefixSpriteVecMathAutomation(t *testing.T) {
	tests := []struct {
		tag  uint16
		want string
	}{
		{heap.TagSpriteGroup, "SPRITEGROUP."},
		{heap.TagSpriteLayer, "SPRITELAYER."},
		{heap.TagSpriteBatch, "SPRITEBATCH."},
		{heap.TagSpriteUI, "SPRITEUI."},
		{heap.TagParticle2D, "PARTICLE2D."},
		{heap.TagQuaternion, "QUAT."},
		{heap.TagColor, "COLOR."},
		{heap.TagVec2, "VEC2."},
		{heap.TagVec3, "VEC3."},
		{heap.TagAutomationList, "EVENT."},
	}
	for _, tc := range tests {
		if p := handleCallRegistryPrefix(tc.tag); p != tc.want {
			t.Fatalf("tag %d: got %q want %q", tc.tag, p, tc.want)
		}
	}
}

func TestHandleCallBuiltinSpriteParticleQuatVecColorAutomation(t *testing.T) {
	cases := []struct {
		tag    uint16
		method string
		want   string
	}{
		{heap.TagSpriteGroup, "DRAW", "SPRITEGROUP.DRAW"},
		{heap.TagSpriteLayer, "SETZ", "SPRITELAYER.SETZ"},
		{heap.TagSpriteBatch, "CLEAR", "SPRITEBATCH.CLEAR"},
		{heap.TagSpriteUI, "FREE", "SPRITEUI.FREE"},
		{heap.TagParticle2D, "EMIT", "PARTICLE2D.EMIT"},
		{heap.TagQuaternion, "SLERP", "QUAT.SLERP"},
		{heap.TagColor, "LERP", "COLOR.LERP"},
		{heap.TagVec2, "TRANSFORMMAT4", "VEC2.TRANSFORMMAT4"},
		{heap.TagVec3, "ROTATEBYQUAT", "VEC3.ROTATEBYQUAT"},
		{heap.TagVec3, "ORTHONORMALIZE", "VEC3.ORTHONORMALIZE"},
		{heap.TagAutomationList, "EXPORT", "EVENT.LISTEXPORT"},
		{heap.TagAutomationList, "SETACTIVE", "EVENT.SETACTIVELIST"},
		{heap.TagAutomationList, "LISTCOUNT", "EVENT.LISTCOUNT"},
	}
	for _, tc := range cases {
		k, prep, ok := handleCallBuiltin(tc.tag, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("tag %d %q: got %q prepend=%v ok=%v want %q", tc.tag, tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallDispatchVecColorAutomationListZeroArg(t *testing.T) {
	// 0-arg getters mapped in handleCallDispatch
	kx, px, ox := handleCallDispatch(heap.TagVec2, "x", 0)
	if !ox || !px || kx != "VEC2.X" {
		t.Fatalf("vec2 x(): got %q prepend=%v ok=%v", kx, px, ox)
	}
	kz, pz, oz := handleCallDispatch(heap.TagVec3, "z", 0)
	if !oz || !pz || kz != "VEC3.Z" {
		t.Fatalf("vec3 z(): got %q prepend=%v ok=%v", kz, pz, oz)
	}
	kr, pr, or := handleCallDispatch(heap.TagColor, "r", 0)
	if !or || !pr || kr != "COLOR.R" {
		t.Fatalf("color r(): got %q prepend=%v ok=%v", kr, pr, or)
	}
	ka, pa, oa := handleCallDispatch(heap.TagColor, "alpha", 0)
	if !oa || !pa || ka != "COLOR.A" {
		t.Fatalf("color alpha(): got %q prepend=%v ok=%v", ka, pa, oa)
	}
	// listCount / size (SETSIZE) map to EVENT.LISTCOUNT with prepend
	kc, pc, oc := handleCallDispatch(heap.TagAutomationList, "listCount", 0)
	if !oc || !pc || kc != "EVENT.LISTCOUNT" {
		t.Fatalf("automation listCount(): got %q prepend=%v ok=%v", kc, pc, oc)
	}
}

func TestHandleCallSuggestionsSpriteVecMathIncludesCoreMethods(t *testing.T) {
	has := func(tag uint16, needle string) bool {
		for _, s := range HandleCallSuggestions(tag) {
			if s == needle {
				return true
			}
		}
		return false
	}
	if !has(heap.TagSpriteGroup, "Draw") || !has(heap.TagParticle2D, "Emit") {
		t.Fatalf("sprite/particle suggestions missing expected methods")
	}
	if !has(heap.TagQuaternion, "Slerp") || !has(heap.TagColor, "Lerp") {
		t.Fatalf("quat/color suggestions missing expected methods")
	}
	if !has(heap.TagVec2, "X") || !has(heap.TagVec3, "RotateByQuat") || !has(heap.TagVec3, "OrthoNormalize") {
		t.Fatalf("vec suggestions missing expected methods")
	}
	if !has(heap.TagAutomationList, "ListCount") {
		t.Fatalf("automation list suggestions missing ListCount")
	}
}

func TestHandleCallRegistryPrefixMatrix(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagMatrix); p != "TRANSFORM." {
		t.Fatalf("TagMatrix prefix: got %q want TRANSFORM.", p)
	}
}

func TestHandleCallBuiltinMatrixTransformOps(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"FREE", "TRANSFORM.FREE"},
		{"ROT", "TRANSFORM.SETROTATION"},
		{"INVERSE", "TRANSFORM.INVERSE"},
		{"TRANSPOSE", "TRANSFORM.TRANSPOSE"},
		{"MULTIPLY", "TRANSFORM.MULTIPLY"},
		{"GETELEMENT", "TRANSFORM.GETELEMENT"},
		{"APPLYX", "TRANSFORM.APPLYX"},
		{"TRANSFORMX", "TRANSFORM.APPLYX"},
		{"toQuat", "QUAT.FROMMAT4"},
		{"FROMMAT4", "QUAT.FROMMAT4"},
	}
	for _, tc := range cases {
		k, prep, ok := handleCallBuiltin(heap.TagMatrix, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagMatrix %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}

func TestHandleCallDispatchMatrixInverseZeroArg(t *testing.T) {
	k, prep, ok := handleCallDispatch(heap.TagMatrix, "inverse", 0)
	if !ok || !prep || k != "TRANSFORM.INVERSE" {
		t.Fatalf("matrix inverse(): got %q prepend=%v ok=%v", k, prep, ok)
	}
}

func TestHandleCallBuiltinPeerSendVsSendPacket(t *testing.T) {
	ks, ps, oks := handleCallBuiltin(heap.TagPeer, "SEND")
	if !oks || !ps || ks != "PEER.SEND" {
		t.Fatalf("peer SEND: got %q prepend=%v ok=%v", ks, ps, oks)
	}
	kp, pp, okp := handleCallBuiltin(heap.TagPeer, "SENDPACKET")
	if !okp || !pp || kp != "PEER.SENDPACKET" {
		t.Fatalf("peer SENDPACKET: got %q prepend=%v ok=%v", kp, pp, okp)
	}
}

func TestHandleCallRegistryPrefixPhysicsBuilder(t *testing.T) {
	if p := handleCallRegistryPrefix(heap.TagPhysicsBuilder); p != "BODY3D." {
		t.Fatalf("TagPhysicsBuilder prefix: got %q want BODY3D.", p)
	}
}

func TestHandleCallBuiltinPhysicsBuilderShapeOps(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"ADDBOX", "BODY3D.ADDBOX"},
		{"COMMIT", "BODY3D.COMMIT"},
		{"FREE", "BODY3D.FREE"},
	}
	for _, tc := range cases {
		k, prep, ok := handleCallBuiltin(heap.TagPhysicsBuilder, tc.method)
		if !ok || !prep || k != tc.want {
			t.Fatalf("TagPhysicsBuilder %q: got %q prepend=%v ok=%v want %q", tc.method, k, prep, ok, tc.want)
		}
	}
}
