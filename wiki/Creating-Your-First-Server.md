# 🎮 Creating Your First Server

This guide will walk you through creating your first FiveM server using InkWash. It's super easy with our interactive wizard!

## 📋 What You'll Need

Before you start, make sure you have:

- ✅ InkWash installed ([Installation Guide](Installation-Guide))
- ✅ A FiveM license key ([How to get one](#getting-a-license-key))
- ✅ About 500MB of free disk space
- ✅ Internet connection

## 🔑 Getting a License Key

1. Go to **[https://keymaster.fivem.net/](https://keymaster.fivem.net/)**
2. Log in with your FiveM account
3. Click **"New Server"**
4. Fill in the server details (you can use anything, it's just for registration)
5. Click **"Generate"**
6. Copy the key (it starts with `cfxk_`)

**Save this key somewhere safe - you'll need it!**

## 🚀 Step-by-Step: Create a Server

### Step 1: Open InkWash

1. Open Command Prompt (Press `Win + R`, type `cmd`, press Enter)
2. Navigate to your InkWash folder:
   ```
   cd C:\path\to\inkwash
   ```
3. Type:
   ```
   inkwash.exe create
   ```

### Step 2: Follow the Wizard

The wizard will ask you several questions. Here's what each one means:

#### 📝 Server Name
**Question:** "What do you want to name your server?"

**Example:** `my-first-server`

**Tips:**
- Use lowercase letters and hyphens
- No spaces (InkWash will fix this automatically if you do)
- Choose something memorable!

#### 🏗️ FiveM Build
**Question:** "Which FiveM build do you want to use?"

**What is this?** The version of FiveM server software.

**Recommended:** Just press Enter to use the recommended version!

**Advanced:** If you need a specific version, you can choose from the list.

#### 🔑 License Key
**Question:** "Which license key do you want to use?"

**Options:**
1. If you haven't added a key yet: Choose **"Add new key"**
2. If you have keys saved: Select one from the list

**Adding a new key:**
1. Select **"Add new key"**
2. Give it a label (like "Main Key" or "Test Key")
3. Paste your license key (the one starting with `cfxk_`)

#### 🔌 Server Port
**Question:** "What port should your server use?"

**Default:** 30120

**Tips:**
- Just press Enter to use the default
- Only change this if you know what you're doing
- Make sure no other program is using this port!

#### 📁 Installation Path
**Question:** "Where do you want to install the server?"

**Default:** A folder next to InkWash

**Tips:**
- The default is usually fine!
- Make sure you have enough space (about 500MB)
- Don't use OneDrive or Dropbox folders

#### ✅ Confirmation
**Final step:** Review your choices!

The wizard will show you everything you selected. If it looks good, press Enter to continue!

### Step 3: Wait for Installation

InkWash will now:
1. ✅ Download FiveM server files (~200MB)
2. ✅ Extract the files
3. ✅ Clone default resources
4. ✅ Generate server.cfg
5. ✅ Create launch scripts

This takes about 2-5 minutes depending on your internet speed.

**You'll see a progress bar for each step!**

### Step 4: Server Created!

🎉 Success! Your server is ready!

InkWash will show you:
- Server name
- Installation path
- Next steps

## 🎮 Starting Your Server

To start your new server:

```bash
inkwash.exe start my-first-server
```

Replace `my-first-server` with whatever you named your server!

**Your server is now running!**

To connect:
1. Open FiveM
2. Press F8 to open console
3. Type: `connect localhost:30120`

## 📊 Managing Your Server

### View All Servers
```bash
inkwash.exe list
```

### Stop Your Server
```bash
inkwash.exe stop my-first-server
```

### View Server Logs
```bash
inkwash.exe logs my-first-server
```

### View Real-Time Logs
```bash
inkwash.exe logs my-first-server --follow
```

## 🎨 Customizing Your Server

Your server's configuration is in `server.cfg`. You can edit this file to:

- Change server name
- Add resources
- Configure onesync
- Set max players
- And more!

**Location:** `YourServerFolder/server.cfg`

## 🔄 Next Steps

Now that you have a server, you might want to:

- **[Add mods](Converting-GTA5-Mods)** - Convert GTA5 mods to FiveM resources
- **[Install resources](Installing-Resources)** - Add custom resources
- **[Configure server](Server-Configuration)** - Advanced settings

## ❓ Troubleshooting

### Server won't start

1. Check your license key is correct: `inkwash.exe key list`
2. Make sure port 30120 isn't in use
3. Check server logs: `inkwash.exe logs my-first-server`

### Can't connect to server

1. Make sure the server is running: `inkwash.exe list`
2. Check you're using the right port (default: 30120)
3. Try `connect localhost:30120` in FiveM console (F8)

### Installation failed

1. Check your internet connection
2. Make sure you have enough disk space (~500MB)
3. Try again - sometimes downloads fail

**Still having issues?** Check the [Troubleshooting Guide](Troubleshooting) or [open an issue](https://github.com/VexoaXYZ/InkWash/issues)!

---

<div align="center">

**[⬅️ Installation Guide](Installation-Guide)** | **[Next: Converting GTA5 Mods ➡️](Converting-GTA5-Mods)**

</div>
