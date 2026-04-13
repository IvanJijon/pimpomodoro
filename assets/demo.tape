Output assets/demo.gif
Set Width 800
Set Height 500
Set FontSize 16
Set Shell "bash"
Set Theme "nord"
Set TypingSpeed 80ms

# Launch pimpomodoro
Type "pimpom --work 25 --break 5"
Enter
Sleep 2s

# Start the timer
Type "s"
Sleep 4s

# Pause
Type "s"
Sleep 2s

# Open task list
Type "t"
Sleep 1.5s

# Add a task
Type "a"
Sleep 1s
Type "Write README"
Tab
Type "3"
Enter
Sleep 1.5s

# Select task as WIP
Enter
Sleep 1.5s

# Confirm task switch
Type "y"
Sleep 1.5s

# Back to timer
Type "t"
Sleep 2s

# Resume timer
Type "s"
Sleep 3s

# Skip to next phase
Type "n"
Sleep 1.5s
Type "y"
Sleep 1.5s

# Start break timer
Type "s"
Sleep 3s

# Quit
Type "q"
Sleep 1s
Type "y"
Sleep 1s
