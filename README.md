# Welcome to Gavin Robson's AI Python Tutor!
## Overview
This AI Python Tutor is a great way for new programmers to get started on their Python journey! For experienced programmers, this tutor can be a tool to freshen up your Python skills. This Python Tutor can:
1. Answer concept questions about the Python language.
2. Provide simple, helpful code examples.
3. Debug the user's own code.
4. Create exercises for the user to complete.
5. Gives constructive, motivating feedback.

## Installation
You can install the program easily by pasting this to the terminal:
```bash
git clone https://github.com/GavinRobson/ai-final GavinRobsonAiFinal
cd GavinRobsonAiFinal
```

## Starting Server
### Windows
* To run the program on windows:
  ```bash
  ./builds/server-windows-amd64.exe
  ```
### MacOS / Linux
* To run the program on MacOS or Linux:
  ```bash
  ./run.sh
  ```
* **NOTE: If you get error `OPENAI_API_KEY must be set`, make sure to `cd` into the newly created `GavinRobsonAiFinal` directory and make sure the `OPENAI_API_KEY` variable is set to a valid key in the `.env` file.**
* Now, the server should be running in the terminal.
* Open a browser and navigate to http://localhost:3000 in order to see the page.
