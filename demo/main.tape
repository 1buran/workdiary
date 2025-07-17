Output demo/main.webm

Set Shell "bash"
Set FontSize 14
Set Padding 2
Set Margin 2
Set Width 720
Set Framerate 60
Set TypingSpeed 20ms
Set PlaybackSpeed 0.5
Set CursorBlink false

Require workdiary
Type "workdiary -h" Enter
Sleep 1
Type "clear" Enter
Type "workdiary cal -c ~/.config/workdiary/config-test.json" Enter
Sleep 3
Type "clear" Enter
Type "workdiary cal -d -c ~/.config/workdiary/config-test.json" Enter
Sleep 3
Type "clear" Enter
Type "workdiary cal -m 4 -c ~/.config/workdiary/config-test.json" Enter
Sleep 3
Type "clear" Enter
Type "workdiary cal -m 5 -y 2024 -c ~/.config/workdiary/config-test.json" Enter
Sleep 3
