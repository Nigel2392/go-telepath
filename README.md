GO-Telepath
========

A simple library to easily communicate your go values to javascript.

Library based on [Telepath](https://github.com/wagtail/telepath) for Wagtail.

## Installation

```bash
go get github.com/Nigel2392/go-telepath@v0.0.3
```

## Usage

First we must define an adapter.

An adapter is an object used to serialize and deserialize go values to javascript.

```go

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

func main() {
    telepath.Register(AlbumAdapter, &Album{})
    telepath.Register(ArtistAdapter, &Artist{})

    album := &Album{
        Name: "The Dark Side of the Moon",
        Artists: []*Artist{
            &Artist{Name: "Pink Floyd"},
        },
    }

    var context = telepath.NewContext()

	var ctx = telepath.NewContext()
	var result, err = ctx.Pack(value)
	if err != nil {
        log.Fatal(err)
    }

    b, err := json.Marshal(result)
    if err != nil {
        log.Fatal(err)
    }
    var telepathJSON = string(b)

    // Pass to JS somehow...
    // ...
}
```

Then we must define the javascript side.

This will deserialize the go values and create the objects.

We will create some classes to represent the go objects.

The classes must also be registered to the `window.telepath` object.

We do ship a `telepath.js` file, but it is only used for testing.

Telepath can be installed via [npm](https://www.npmjs.com/package/telepath-unpack).

```javascript
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

// If you haven't already instantiated the telepath object
// window.telepath = new Telepath();

window.telepath.register("js.funcs.Album", Album);
window.telepath.register("js.funcs.Artist", Artist);

// Now you can use the go values

// Lets assume they are stored in a variable called `telepathJSON`
var telepathValue = JSON.parse(telepathJSON);

var album = telepath.unpack(telepathValue);

console.log(album);

// Output:
// Album {
//   name: 'The Dark Side of the Moon',
//   artists: [ Artist { name: 'Pink Floyd' } ]
// }
```
