<h1 align="center">Pat</h1>

<h3 align="center">
    <code>$PATH</code> Manager
</h1>

https://user-images.githubusercontent.com/62389790/195914752-6da84193-e39d-4296-9ef4-f94a95cf250e.mp4

## What is this?

Have you ever found yourself searching something like *"Adding a new entry to the PATH variable in ZSH"*, *"How to delete entry from PATH in bash"*, *"How to change order of the PATH entries in fish"*...? And then you forget how to do it and search it again, again and again... ughh, this is 
annoying 💀

***Try Pat!*** It's a `$PATH` environment variable manager. It can reorder, delete and add new path entries with a fancy and easy TUI, so you don't have to remember any commands.

Okay, let's move on!

## How to use?

Let's assume you're a *[Zsh](https://en.wikipedia.org/wiki/Z_shell) user and want to add `~/playground/scripts` to your `$PATH`. *though, fish and bash are also supported, it's just an example*

No need for any extra arguments or preconfiguration just run
```shell
pat
```
It will welcome you with this screen. Select `zsh` (you can move up by pressing arrow up `↑` or `h` and `↓/j` to move down)

<img width="912" alt="Screenshot 2022-10-14 at 21 35 23" src="https://user-images.githubusercontent.com/62389790/195917180-9e20004f-8a75-4e62-977a-a71da1a39186.png">

You will you paths. You can preview each (check what executables it contains) by pressing `p` (acronym for `preview`)

<img width="912" alt="Screenshot 2022-10-14 at 21 37 02" src="https://user-images.githubusercontent.com/62389790/195917418-5ce27dac-33ec-44ca-831f-3e66af6ac899.png">

Press `a` (for `add`) to add a new path. It will open a textinput with smart hints and autocompletion so you won't be lost or misspell something while typing

See, **pat** is smart! You will provide autocompletions while typing

<img width="912" alt="Screenshot 2022-10-14 at 21 41 24" src="https://user-images.githubusercontent.com/62389790/195918178-0f636fe5-d1c2-473e-951c-3fc6a72ccbcc.png">

Press `tab` to accept completion

<img width="912" alt="Screenshot 2022-10-14 at 21 41 41" src="https://user-images.githubusercontent.com/62389790/195918215-947763d6-72b6-46fe-99f0-6885fb8b585f.png">

Press `enter` to confirm

<img width="912" alt="Screenshot 2022-10-14 at 21 42 43" src="https://user-images.githubusercontent.com/62389790/195918384-23ec4722-9f70-4fc3-be7f-1bda46fa42c2.png">

Almost there! Now you will need to *save* it. Simply press `s` or `enter` to save. Don't worry it will ask you to confirm your actions before applying anything

<img width="912" alt="Screenshot 2022-10-14 at 21 43 54" src="https://user-images.githubusercontent.com/62389790/195918581-45898bf2-433a-4325-8135-a3d121825483.png">

Press `Y` to confirm. This keybind is inconvenient on purpose so you won't do anything 
accidentally

**LAST STEP**. You will need to do it only once and forget forever! 

```sh
echo "source $(pat where --zsh) &>/dev/null" >> ~/.zshenv
```

```zsh
# This line should be the last one in your ~/.zshenv
source $(pat where --zsh) &>/dev/null
```

For any other shell steps are the same, except for the last one

For **Fish**
```fish
# If you know what you're doing, you can change config.fish to other file of course
echo "source $(pat where --fish) &>/dev/null" >> ~/.config/fish/config.fish
```

For **Bash**
```bash
echo "source $(path where --bash) &>/dev/null" >> ~/.profile
```

