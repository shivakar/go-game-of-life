## go-game-of-life

Conway's Game of Life<sup>1</sup> written in Go using Pixel 2D Game Library<sup>2</sup>.

* <sup>1</sup> https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life
* <sup>2</sup> https://github.com/faiface/pixel

### Demo
---

![Game Of Life Run](./GameOfLife.gif?raw=true)

### Installation
---

```
go get -u github.com/shivakar/go-game-of-life
```

### Usage
---

To run a Game of Life simulation:

```
go-game-of-life
```

Features:

* Left-click on a pixel to make it alive
* Left-click and drag to continuously add live pixels
* Hold `Ctrl` key when left-click (and dragging) to reinitialize 11x11 grid
  around the mouse position
* Press `r` to reinitialize
* Press `Space` to pause and restart


### License
---

`go-game-of-life` is licensed under a MIT license.
