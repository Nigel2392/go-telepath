package telepath_test

import (
	"encoding/json"
	"strings"
	"testing"

	_ "embed"

	"github.com/Nigel2392/go-telepath/telepath"
	"github.com/dop251/goja"
)

var AlbumAdapter = &telepath.ObjectAdapter{
	JSConstructor: "js.funcs.Album",
	GetJSArgs: func(obj interface{}) []interface{} {
		album := obj.(*Album)
		return []interface{}{album.Name, album.Artists}
	},
}

var ArtistAdapter = &telepath.ObjectAdapter{
	JSConstructor: "js.funcs.Artist",
	GetJSArgs: func(obj interface{}) []interface{} {
		artist := obj.(*Artist)
		return []interface{}{artist.Name}
	},
}

type Album struct {
	Name    string
	Artists []*Artist
}

type Artist struct {
	Name string
}

type TelepathEncoder struct {
	json.Encoder
}

func TestPacking(t *testing.T) {
	telepath.Register(AlbumAdapter, &Album{})
	telepath.Register(ArtistAdapter, &Artist{})

	t.Run("TestPackObject", func(t *testing.T) {
		var object = &Album{Name: "Hello"}
		var ctx = telepath.NewContext()
		var result, err = ctx.Pack(object)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		var chk = result.(telepath.TelepathValue)
		if chk.Type != "js.funcs.Album" {
			t.Errorf("Expected js.funcs.Album, got %v", chk.Type)
		}

		if chk.Args[0] != "Hello" {
			t.Errorf("Expected Hello, got %v", chk.Args[0])
		}

	})

	t.Run("TestPackList", func(t *testing.T) {
		var object = []*Album{
			{Name: "Hello 1"},
			{Name: "Hello 2"},
		}

		var ctx = telepath.NewContext()
		var result, err = ctx.Pack(object)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
			return
		}

		var chk = result.(telepath.TelepathValue)
		if chk.List == nil {
			t.Errorf("Expected list, got nil")
			return
		}

		if len(chk.List) != 2 {
			t.Errorf("Expected 2, got %v", len(chk.List))
			return
		}

		if chk.List[0].(telepath.TelepathValue).Type != "js.funcs.Album" {
			t.Errorf("Expected js.funcs.Album, got %v", chk.List[0].(telepath.TelepathValue).Type)
		}

		if chk.List[0].(telepath.TelepathValue).Args[0] != "Hello 1" {
			t.Errorf("Expected Hello 1, got %v", chk.List[0].(telepath.TelepathValue).Args[0])
		}

		if chk.List[1].(telepath.TelepathValue).Type != "js.funcs.Album" {
			t.Errorf("Expected js.funcs.Album, got %v", chk.List[1].(telepath.TelepathValue).Type)
		}

		if chk.List[1].(telepath.TelepathValue).Args[0] != "Hello 2" {
			t.Errorf("Expected Hello 2, got %v", chk.List[1].(telepath.TelepathValue).Args[0])
		}
	})

	t.Run("TestPackMap", func(t *testing.T) {

		var object = map[string]*Album{
			"1": {Name: "Hello 1"},
			"2": {Name: "Hello 2"},
		}

		var ctx = telepath.NewContext()
		var result, err = ctx.Pack(object)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
			return
		}

		var chk = result.(map[string]interface{})
		if len(chk) != 2 {
			t.Errorf("Expected 2, got %v", len(chk))
			return
		}

		if chk["1"].(telepath.TelepathValue).Type != "js.funcs.Album" {
			t.Errorf("Expected js.funcs.Album, got %v", chk["1"].(telepath.TelepathValue).Type)
		}

		if chk["1"].(telepath.TelepathValue).Args[0] != "Hello 1" {
			t.Errorf("Expected Hello 1, got %v", chk["1"].(telepath.TelepathValue).Args[0])
		}

		if chk["2"].(telepath.TelepathValue).Type != "js.funcs.Album" {
			t.Errorf("Expected js.funcs.Album, got %v", chk["2"].(telepath.TelepathValue).Type)
		}

		if chk["2"].(telepath.TelepathValue).Args[0] != "Hello 2" {
			t.Errorf("Expected Hello 2, got %v", chk["2"].(telepath.TelepathValue).Args[0])
		}
	})

	t.Run("TestDictReservedWords", func(t *testing.T) {
		var object = map[string]interface{}{
			"_artist": &Album{Name: "Hello"},
			"_type":   "Album",
		}

		var ctx = telepath.NewContext()
		var result, err = ctx.Pack(object)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
			return
		}

		var chk = result.(telepath.TelepathValue)
		if chk.Dict == nil {
			t.Errorf("Expected dict, got nil")
			return
		}

		if chk.Dict["_artist"].(telepath.TelepathValue).Type != "js.funcs.Album" {
			t.Errorf("Expected js.funcs.Album, got %v", chk.Dict["_artist"].(telepath.TelepathValue).Type)
		}

		if chk.Dict["_artist"].(telepath.TelepathValue).Args[0] != "Hello" {
			t.Errorf("Expected Hello, got %v", chk.Dict["_artist"].(telepath.TelepathValue).Args[0])
		}

		if len(chk.Dict["_artist"].(telepath.TelepathValue).Args) != 2 {
			t.Errorf("Expected 2, got %v", len(chk.Dict["_artist"].(telepath.TelepathValue).Args))
		}

		if chk.Dict["_type"] != "Album" {
			t.Errorf("Expected Album, got %v", chk.Dict["_type"])
		}

	})

	t.Run("TestRecursiveArgPacking", func(t *testing.T) {
		var object = &Album{
			Name: "Hello",
			Artists: []*Artist{
				{Name: "Artist 1"},
				{Name: "Artist 2"},
			},
		}

		telepath.Register(ArtistAdapter, &Artist{})

		var ctx = telepath.NewContext()
		var result, err = ctx.Pack(object)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
			return
		}

		var chk = result.(telepath.TelepathValue)
		if chk.Type != "js.funcs.Album" {
			t.Errorf("Expected js.funcs.Album, got %v", chk.Type)
		}

		if chk.Args[0] != "Hello" {
			t.Errorf("Expected Hello, got %v", chk.Args[0])
		}

		if chk.Args[1].(telepath.TelepathValue).List[0].(telepath.TelepathValue).Type != "js.funcs.Artist" {
			t.Errorf("Expected js.funcs.Artist, got %v", chk.Args[1].(telepath.TelepathValue).List[0].(telepath.TelepathValue).Type)
		}

		if chk.Args[1].(telepath.TelepathValue).List[0].(telepath.TelepathValue).Args[0] != "Artist 1" {
			t.Errorf("Expected Artist 1, got %v", chk.Args[1].(telepath.TelepathValue).List[1].(telepath.TelepathValue).Args[0])
		}

		if chk.Args[1].(telepath.TelepathValue).List[1].(telepath.TelepathValue).Type != "js.funcs.Artist" {
			t.Errorf("Expected js.funcs.Artist, got %v", chk.Args[1].(telepath.TelepathValue).List[0].(telepath.TelepathValue).Type)
		}

		if chk.Args[1].(telepath.TelepathValue).List[1].(telepath.TelepathValue).Args[0] != "Artist 2" {
			t.Errorf("Expected Artist 2, got %v", chk.Args[1].(telepath.TelepathValue).List[1].(telepath.TelepathValue).Args[0])
		}
	})
}

type StringLike struct {
	Value string
}

var StringLikeAdapter = &telepath.ObjectAdapter{
	JSConstructor: "js.funcs.StringLike",
	GetJSArgs: func(obj interface{}) []interface{} {
		str := obj.(*StringLike)
		return []interface{}{strings.ToUpper(str.Value)}
	},
}

func TestPackingToString(t *testing.T) {
	var value = []any{
		&StringLike{Value: "hello"},
		"world",
	}

	telepath.Register(StringLikeAdapter, &StringLike{})

	var ctx = telepath.NewContext()
	var result, err = ctx.Pack(value)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	var chk = result.(telepath.TelepathValue)
	if chk.List[0].(telepath.TelepathValue).Type != "js.funcs.StringLike" {
		t.Errorf("Expected js.funcs.StringLike, got %v", chk.List[0].(telepath.TelepathValue).Type)
	}

	if chk.List[0].(telepath.TelepathValue).Args[0] != "HELLO" {
		t.Errorf("Expected HELLO, got %v", chk.List[0].(telepath.TelepathValue).Args[0])
	}

	if chk.List[1] != "world" {
		t.Errorf("Expected world, got %v", chk.List[1])
	}
}

var _ telepath.AdapterGetter = (*TelepathAdapterGetterStruct)(nil)

type TelepathAdapterGetterStruct struct {
	name string
}

func (m *TelepathAdapterGetterStruct) Adapter() telepath.Adapter {
	return &telepath.ObjectAdapter{
		JSConstructor: "js.funcs." + m.name,
		GetJSArgs: func(obj interface{}) []interface{} {
			return []interface{}{obj.(*TelepathAdapterGetterStruct).name}
		},
	}
}

func TestAdapterGetter(t *testing.T) {
	var value = &TelepathAdapterGetterStruct{name: "Getter"}

	telepath.Register(value)

	var ctx = telepath.NewContext()
	var result, err = ctx.Pack(value)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	var chk = result.(telepath.TelepathValue)
	if chk.Type != "js.funcs.Getter" {
		t.Errorf("Expected js.funcs.Getter, got %v", chk.Type)
	}

	if chk.Args[0] != "Getter" {
		t.Errorf("Expected %v, got %v", value, chk.Args[0])
	}
}

//go:embed fixtures/telepath-0.0.3.index.js
var telepath_js string

const (
	vm_js = `
class Album {
	constructor(name, artists) {
		this.name = name;
		this.artists = artists;
	}
}

class Artist {
	constructor(name) {
		this.name = name;
	}
}
TELEPATH.register("js.funcs.Album", Album);
TELEPATH.register("js.funcs.Artist", Artist);`
)

func TestTelepathUnpack(t *testing.T) {
	var value = &Album{
		Name: "Hello",
		Artists: []*Artist{
			{Name: "Artist 1"},
			{Name: "Artist 2"},
		},
	}
	var ctx = telepath.NewContext()
	var result, err = ctx.Pack(value)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	var resultJSON, _ = json.Marshal(result)

	vm := goja.New()
	_, err = vm.RunString(telepath_js)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	_, err = vm.RunString(vm_js)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	vm.Set("testData", string(resultJSON))

	_, err = vm.RunString(`testData = JSON.parse(testData);`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	_, err = vm.RunString(`var data = TELEPATH.unpack(testData);`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	isData, err := vm.RunString(`data instanceof Album`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if !isData.ToBoolean() {
		t.Errorf("Expected true, got %v", isData.ToBoolean())
	}

	name, err := vm.RunString(`data.name`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if name.String() != "Hello" {
		t.Errorf("Expected Hello, got %v", name.String())
	}

	isArtists, err := vm.RunString(`data.artists instanceof Array`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if !isArtists.ToBoolean() {
		t.Errorf("Expected true, got %v", isArtists.ToBoolean())
	}

	isArtist1, err := vm.RunString(`data.artists[0] instanceof Artist`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if !isArtist1.ToBoolean() {
		t.Errorf("Expected true, got %v", isArtist1.ToBoolean())
	}

	artist1Name, err := vm.RunString(`data.artists[0].name`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if artist1Name.String() != "Artist 1" {
		t.Errorf("Expected Artist 1, got %v", artist1Name.String())
	}

	isArtist2, err := vm.RunString(`data.artists[1] instanceof Artist`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if !isArtist2.ToBoolean() {
		t.Errorf("Expected true, got %v", isArtist2.ToBoolean())
	}

	artist2Name, err := vm.RunString(`data.artists[1].name`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if artist2Name.String() != "Artist 2" {
		t.Errorf("Expected Artist 2, got %v", artist2Name.String())
	}
}

func TestReadmeExample(t *testing.T) {
	const albumJS = `class Album {
	constructor(name, artists) {
		this.name = name;
		this.artists = artists;
	}
}

class Artist {
	constructor(name) {
		this.name = name;
	}
}

// If you haven't already instantiated the telepath object
// window.telepath = new Telepath();

TELEPATH.register("js.funcs.Album", Album);
TELEPATH.register("js.funcs.Artist", Artist);

// Now you can use the go values

// Lets assume they are stored in a variable called 'telepathJSON'
var telepathValue = JSON.parse(telepathJSON);

var album = TELEPATH.unpack(telepathValue);`

	var album = &Album{
		Name: "Hello",
		Artists: []*Artist{
			{Name: "Artist 1"},
			{Name: "Artist 2"},
		},
	}

	telepath.Register(
		AlbumAdapter, &Album{},
	)
	telepath.Register(
		ArtistAdapter, &Artist{},
	)

	var ctx = telepath.NewContext()
	var result, err = ctx.Pack(album)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	var resultJSON, _ = json.Marshal(result)

	vm := goja.New()
	_, err = vm.RunString(telepath_js)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	vm.Set("telepathJSON", string(resultJSON))

	_, err = vm.RunString(albumJS)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	isData, err := vm.RunString(`album instanceof Album`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if !isData.ToBoolean() {
		t.Errorf("Expected true, got %v", isData.ToBoolean())
	}

	name, err := vm.RunString(`album.name`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if name.String() != "Hello" {
		t.Errorf("Expected Hello, got %v", name.String())
	}

	isArtists, err := vm.RunString(`album.artists instanceof Array`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if !isArtists.ToBoolean() {
		t.Errorf("Expected true, got %v", isArtists.ToBoolean())
	}

	isArtist1, err := vm.RunString(`album.artists[0] instanceof Artist`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if !isArtist1.ToBoolean() {
		t.Errorf("Expected true, got %v", isArtist1.ToBoolean())
	}

	artist1Name, err := vm.RunString(`album.artists[0].name`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if artist1Name.String() != "Artist 1" {
		t.Errorf("Expected Artist 1, got %v", artist1Name.String())
	}

	isArtist2, err := vm.RunString(`album.artists[1] instanceof Artist`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if !isArtist2.ToBoolean() {
		t.Errorf("Expected true, got %v", isArtist2.ToBoolean())
	}

	artist2Name, err := vm.RunString(`album.artists[1].name`)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if artist2Name.String() != "Artist 2" {
		t.Errorf("Expected Artist 2, got %v", artist2Name.String())
	}
}
