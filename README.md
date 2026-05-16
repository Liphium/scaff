# Scaff

> [!WARNING]
> Scaff is an experimental project not intended for production use cases.

Scaff (shortened from scaffold) is quite literally a scaffold for point-and-click games that don't require more than UI. We provide a minimal framework capable of everything you need for the kind of game that doesn't require a bloated game engine with collisions, 3D and more.

The vision is to make a really simple and efficient game engine that uses a giant tree of nodes to make development really fast and predictable. We've seen that trees work really well for UI (looking at stuff like React), so why not for everything in your 2D environment too?

## Features

> [!NOTE]
> We're still working on most of this, especially documentation and the canvas package. This is just what should be there in the future.

Scaff provides various features that make your life a lot easier:

- Raw access at every point in case you need it (and with good documentation on how too, hopefully)
- Signals for state tracking and only refreshing what's really needed
- `scaffui` (Scaff User Interface) package for rich declarative UI based on state updates
- `scaffcv` (Scaff Canvas) package for rich game environments reacting to state updates
- Loggers, based on `slog`, that make logging a breeze and actually look really nice
- A scene stack for easy scene management

This is just a really small overview. A lot of the functionality that makes Scaff special are the little packages and how they are made.

## Credits

| Name   | Description                        | License | Link                                | Note                                                                                     |
| ------ | ---------------------------------- | ------- | ----------------------------------- | ---------------------------------------------------------------------------------------- |
| V      | Vector math and position utilities | MIT     | https://github.com/setanarut/v      | Thanks for helping out with vector math!                                                 |
| Kamera | Camera movement and math           | CC0     | https://github.com/setanarut/kamera | Thanks for making this amazing package, making it myself would've been really difficult. |
