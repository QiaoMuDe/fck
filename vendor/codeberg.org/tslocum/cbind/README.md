# cbind
[![GoDoc](https://codeberg.org/tslocum/godoc-static/raw/branch/main/badge.svg)](https://docs.rocket9labs.com/codeberg.org/tslocum/cbind)
[![Donate](https://img.shields.io/liberapay/receives/rocket9labs.com.svg?logo=liberapay)](https://liberapay.com/rocket9labs.com)

Key event handling library for tcell

## Features

- Set `KeyEvent` handlers
- Encode and decode `KeyEvent`s as human-readable strings

## Usage

```go
// Create a new input configuration to store the key bindings.
c := NewConfiguration()

// Define save event handler.
handleSave := func(ev *tcell.EventKey) *tcell.EventKey {
    return nil
}

// Define open event handler.
handleOpen := func(ev *tcell.EventKey) *tcell.EventKey {
    return nil
}

// Define exit event handler.
handleExit := func(ev *tcell.EventKey) *tcell.EventKey {
    return nil
}

// Bind Alt+s.
if err := c.Set("Alt+s", handleSave); err != nil {
    log.Fatalf("failed to set keybind: %s", err)
}

// Bind Alt+o.
c.SetRune(tcell.ModAlt, 'o', handleOpen)

// Bind Escape.
c.SetKey(tcell.ModNone, tcell.KeyEscape, handleExit)

// Capture input. This will differ based on the framework in use (if any).
// When using tview or cview, call Application.SetInputCapture before calling
// Application.Run.
app.SetInputCapture(c.Capture)
```

## Documentation

Documentation is available via [gdooc](https://docs.rocket9labs.com/codeberg.org/tslocum/cbind).

The utility program `whichkeybind` is available to determine and validate key combinations.

```bash
go install codeberg.org/tslocum/cbind/whichkeybind@latest
```

## Support

Please share issues and suggestions [here](https://codeberg.org/tslocum/cbind/issues).
