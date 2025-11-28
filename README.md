# 🎨 InkWash - Easy FiveM Server Manager

> The simplest way to create and manage FiveM servers. Just download, click, and go!

[![Download](https://img.shields.io/github/v/release/VexoaXYZ/InkWash?label=Download&style=for-the-badge&logo=github)](https://github.com/VexoaXYZ/InkWash/releases/latest)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows-blue?style=for-the-badge&logo=windows)](https://github.com/VexoaXYZ/InkWash)

---

## 🚀 Super Quick Start (3 Steps!)

### Step 1: Download InkWash
👉 **[Click here to download the latest version](https://github.com/VexoaXYZ/InkWash/releases/latest)**

Look for the file called `inkwash-windows-amd64.zip` and download it.

### Step 2: Extract the File
1. Right-click on the downloaded ZIP file
2. Click "Extract All..."
3. Choose where you want to extract it
4. Click "Extract"

### Step 3: Run InkWash
1. Open the folder where you extracted the files
2. Double-click on `inkwash.exe`
3. Done! The program is now running

---

## 📦 What Can InkWash Do?

✅ **Create FiveM Servers** - Set up a new server in minutes with a step-by-step wizard
✅ **Convert GTA5 Mods** - Turn GTA5 mods into FiveM resources automatically
✅ **Manage Servers** - Start, stop, and monitor all your servers easily
✅ **Beautiful Interface** - Modern, colorful terminal UI that's easy to understand
✅ **Safe & Secure** - Your license keys are encrypted and stored safely

---

## 🎮 How to Use InkWash

### Creating Your First Server

1. Open Command Prompt (Press `Win + R`, type `cmd`, press Enter)
2. Navigate to where you extracted InkWash
3. Type: `inkwash.exe create`
4. Follow the wizard! It will ask you simple questions like:
   - What do you want to name your server?
   - Which FiveM version do you want?
   - What's your license key?

The wizard guides you through everything step-by-step!

### Converting GTA5 Mods to FiveM

1. Open Command Prompt
2. Navigate to your InkWash folder
3. Type: `inkwash.exe convert`
4. The wizard will ask you:
   - Which server do you want to add mods to?
   - What's the GTA5-mods.com URL?
5. InkWash will automatically download, convert, and install the mod!

### Managing Your Servers

**Start a Server:**
```
inkwash.exe start my-server-name
```

**Stop a Server:**
```
inkwash.exe stop my-server-name
```

**See All Servers:**
```
inkwash.exe list
```

**View Server Logs:**
```
inkwash.exe logs my-server-name
```

---

## 💡 Common Questions

### Where do I get a FiveM license key?
1. Go to https://keymaster.fivem.net/
2. Log in with your FiveM account
3. Create a new server key
4. Copy the key (it starts with `cfxk_`)

### How do I add my license key to InkWash?
```
inkwash.exe key add
```
Then paste your license key when it asks!

### My server isn't starting, what do I do?
1. Check your license key is correct: `inkwash.exe key list`
2. Check server logs: `inkwash.exe logs your-server-name`
3. Make sure no other program is using port 30120

### Can I use InkWash without knowing how to code?
**Yes!** InkWash is designed to be super easy. You don't need to know any coding. Just follow the wizard and answer the questions!

---

## 🛠️ All Commands

### Server Commands

| Command | What it does |
|---------|-------------|
| `inkwash.exe create` | Create a new server (opens wizard) |
| `inkwash.exe start <name>` | Start a server |
| `inkwash.exe stop <name>` | Stop a server |
| `inkwash.exe list` | Show all your servers |
| `inkwash.exe logs <name>` | View server logs |

### Mod Converter Commands

| Command | What it does |
|---------|-------------|
| `inkwash.exe convert` | Convert GTA5 mods (opens wizard) |

### License Key Commands

| Command | What it does |
|---------|-------------|
| `inkwash.exe key add` | Add a license key |
| `inkwash.exe key list` | Show all your keys (hidden) |
| `inkwash.exe key remove <id>` | Delete a key |

---

## 🎯 Features in Detail

### 🧙 Interactive Wizards
Never get lost! Our wizards guide you through every step with helpful hints and tips.

### 🎨 Beautiful UI
InkWash looks good and is easy to read with colors that help you understand what's happening.

### ⚡ Super Fast
InkWash is built in Go, making it lightning fast. Servers start in seconds!

### 🔐 Secure
Your license keys are encrypted with military-grade AES-256 encryption. They're safe!

### 🌐 Auto-Updates
Download new versions from our [Releases page](https://github.com/VexoaXYZ/InkWash/releases) whenever they come out!

---

## 📚 Need More Help?

- 📖 [Check out our Wiki](https://github.com/VexoaXYZ/InkWash/wiki) for detailed guides
- 🐛 [Report a Bug](https://github.com/VexoaXYZ/InkWash/issues)
- 💬 Join our Discord: [Coming Soon]
- 📧 Email: [Your Email]

---

## 🎓 For Advanced Users

### Install with Go
If you're a developer and have Go installed:
```bash
go install github.com/VexoaXYZ/inkwash@latest
```

### Build from Source
```bash
git clone https://github.com/VexoaXYZ/InkWash.git
cd InkWash
go build -o inkwash.exe .
```

### Add to PATH (Windows)
So you can use `inkwash` from anywhere:
1. Press `Win + R`, type `sysdm.cpl`, press Enter
2. Go to "Advanced" tab
3. Click "Environment Variables"
4. Under "System Variables", find "Path"
5. Click "Edit" → "New"
6. Paste the folder path where `inkwash.exe` is located
7. Click OK on everything
8. Restart Command Prompt

Now you can just type `inkwash` from anywhere!

---

## 📜 License

InkWash is free and open source under the [MIT License](LICENSE).

---

## ❤️ Made with Love

Created by [Vexoa](https://github.com/VexoaXYZ) to make FiveM server management easy for everyone.

**Version 2.0** - Complete rewrite with better everything!

---

<div align="center">

### Ready to get started?

**[📥 Download InkWash Now](https://github.com/VexoaXYZ/InkWash/releases/latest)**

</div>
