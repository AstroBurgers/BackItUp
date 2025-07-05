# 📂Back It Up

A minimal terminal app built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) that lets you:

- ✅ Select file extensions to include in backups  
- ✅ Scan your current directory for those files  
- ✅ Zip them into a structured archive (preserving folder paths)  
- ✅ Save backups as `backitup-YYYY-MM-DD_HH-MM-SS.zip`

---

## ✨ Features

- **Smart Scanning**: Recursively finds files by extension in the current working directory  
- **Structure-Preserving Zip**: Files are zipped relative to where you run the app  
- **Config Editor**: Easily add or remove file extensions with keyboard input  
- **Progress Bar**: Simple visual feedback as files are backed up  
- **Offline-first**: No external APIs, just runs clean in your terminal

---

## 🛠️ Built With
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) – Terminal UI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) – Reusable UI components
- Go standard library (zip, filepath, os, etc.)

---

For any questions, comments, or issues, shoot me a message on discord @`astro.1181`!
