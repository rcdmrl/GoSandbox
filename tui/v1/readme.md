Inspired by this: https://www.youtube.com/watch?v=gXmznGEW9vo&t=588s

Creates a TUI to run all the examples on this project. 

Uses:
- https://github.com/charmbracelet/huh (recommended by the video above)

Details:
- Starts by asking the user to choose a project
- Then, for each project, it asks to choose one of the available versions
- Based on that, it runs the selected project and version
- There's also an option to call quits

The "meh" stuff:
- Each project can have its own set of versions
- This means the second form depends on the selection of the first one
- Not sure how to do this with Huh? so I've created this "branching" logic manually
