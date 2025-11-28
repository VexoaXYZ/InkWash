# 🔧 Troubleshooting Guide

Having problems with InkWash? This guide will help you solve common issues!

## 📋 Quick Checks

Before diving into specific problems, try these:

1. ✅ Make sure you're using the **latest version** of InkWash
2. ✅ Check your **internet connection**
3. ✅ **Restart Command Prompt** and try again
4. ✅ Make sure you have **enough disk space** (~500MB minimum)

---

## 🚀 Installation Issues

### "Windows protected your PC"

**Problem:** Windows SmartScreen is blocking InkWash

**Solution:**
1. Click **"More info"**
2. Click **"Run anyway"**
3. This is normal for new downloads!

**Why this happens:** InkWash isn't signed with a Windows certificate (costs money), but it's completely safe!

### Can't extract the ZIP file

**Problem:** ZIP file won't extract or is corrupted

**Solutions:**
1. Try downloading again (file might be corrupted)
2. Use a different extraction tool (7-Zip, WinRAR)
3. Make sure you have enough disk space
4. Check your antivirus isn't blocking it

### InkWash.exe won't run

**Problem:** Double-clicking inkwash.exe does nothing

**Solutions:**
1. Make sure you **extracted** the ZIP first (don't run from inside ZIP!)
2. Check you downloaded the **Windows version** (inkwash-windows-amd64.zip)
3. Try running as Administrator (right-click → "Run as administrator")
4. Check your antivirus isn't blocking it

---

## 🎮 Server Creation Issues

### "Failed to download FiveM build"

**Problem:** Can't download server files

**Solutions:**
1. Check your internet connection
2. Check if https://runtime.fivem.net is accessible
3. Try again later (servers might be busy)
4. Check your firewall isn't blocking InkWash

### "License key validation failed"

**Problem:** Your license key doesn't work

**Solutions:**
1. Check the key is correct (starts with `cfxk_`)
2. Make sure you copied the **entire key**
3. No spaces before or after the key
4. Get a new key from https://keymaster.fivem.net

### "Server name already exists"

**Problem:** You already have a server with that name

**Solutions:**
1. Choose a different name
2. Or delete the old server first: `inkwash.exe remove old-server-name`

### "Path already exists"

**Problem:** The install folder already exists

**Solutions:**
1. Choose a different path
2. Or delete the existing folder
3. Or use the existing folder (InkWash will ask)

---

## 🎛️ Server Management Issues

### Server won't start

**Problem:** Server fails to start or crashes immediately

**Solutions:**

1. **Check license key:**
   ```
   inkwash.exe key list
   ```
   Make sure your key is valid!

2. **Check port 30120 isn't in use:**
   ```
   netstat -ano | findstr :30120
   ```
   If something is using it, change your server port

3. **Check server logs:**
   ```
   inkwash.exe logs your-server-name
   ```
   Look for error messages

4. **Common errors:**
   - "Invalid license key" → Add correct key
   - "Address already in use" → Port 30120 is taken
   - "Missing resources" → Server files are corrupted, recreate server

### Can't connect to server

**Problem:** Server is running but you can't connect

**Solutions:**

1. **Make sure server is actually running:**
   ```
   inkwash.exe list
   ```
   Should show "Running" status

2. **Try connecting with:**
   ```
   connect localhost:30120
   ```
   In FiveM console (press F8)

3. **Check firewall:**
   - Windows Firewall might be blocking port 30120
   - Add an exception for FXServer.exe

4. **Check port:**
   - Default is 30120
   - If you changed it, use the correct port: `connect localhost:YOUR_PORT`

### Server keeps crashing

**Problem:** Server starts but crashes after a while

**Solutions:**
1. Check server logs for errors
2. Remove recently added resources (might be broken)
3. Check RAM usage (server might be out of memory)
4. Try a different FiveM build

### "Server not found" error

**Problem:** InkWash can't find your server

**Solutions:**
1. Check server name is correct: `inkwash.exe list`
2. Server folder might have been deleted or moved
3. Re-create the server if needed

---

## 🚗 Mod Converter Issues

### "URL must be from gta5-mods.com"

**Problem:** Trying to convert from wrong website

**Solution:**
- Only gta5-mods.com URLs work
- Find the mod on gta5-mods.com
- Copy the full URL from your browser

### "Conversion failed: 404"

**Problem:** Mod page doesn't exist or was removed

**Solutions:**
1. Check the URL is correct
2. Try visiting the URL in your browser
3. Mod might have been deleted - find a different one
4. Make sure it's a direct mod page (not a category)

### "Download failed: 404"

**Problem:** Converted file can't be downloaded

**Solutions:**
1. Wait a few minutes and try again
2. Converting too many at once - try fewer mods (5 at a time)
3. convert.cfx.rs might be down - try later

### Mod converted but not in game

**Problem:** Mod is downloaded but doesn't appear

**Solutions:**

1. **Add to server.cfg:**
   ```
   ensure mod-name
   ```
   Check the exact folder name in `resources/[category]/`

2. **Restart server:**
   ```
   inkwash.exe stop your-server
   inkwash.exe start your-server
   ```

3. **Check mod requirements:**
   - Some mods need dependencies
   - Check the original mod page for instructions
   - Some mods are "replace" not "addon" (different installation)

### UI is scrolling too much

**Problem:** Progress display is hard to read

**Solution:**
- This was fixed in v2.0!
- Update to the latest version
- New version has throttled updates (500ms)

---

## 🔑 License Key Issues

### Can't add license key

**Problem:** Key add command fails

**Solutions:**
1. Make sure key starts with `cfxk_`
2. Copy the entire key (no spaces before/after)
3. Get a new key from https://keymaster.fivem.net

### "License key already exists"

**Problem:** You already added this key

**Solutions:**
1. Use the existing key (it's saved!)
2. List your keys: `inkwash.exe key list`
3. No need to add it again

### Can't see my license keys

**Problem:** `key list` shows nothing

**Solutions:**
1. You haven't added any keys yet
2. Add one: `inkwash.exe key add`
3. Keys are encrypted and stored securely

---

## 💻 Command Prompt Issues

### "inkwash is not recognized as an internal or external command"

**Problem:** Command Prompt can't find inkwash.exe

**Solutions:**

1. **Navigate to InkWash folder first:**
   ```
   cd C:\path\to\inkwash
   ```

2. **Use full filename:**
   ```
   inkwash.exe create
   ```

3. **Or add to PATH** (see [Installation Guide](Installation-Guide))

### Commands aren't working

**Problem:** Commands do nothing or show errors

**Solutions:**
1. Check spelling: `inkwash.exe --help`
2. Use correct syntax: `inkwash.exe command server-name`
3. Make sure you're in the right folder
4. Try running as Administrator

---

## 🆘 Still Having Problems?

If none of these solutions work:

1. 📖 **Check the [Wiki](https://github.com/VexoaXYZ/InkWash/wiki)** for more guides
2. 🔍 **Search [existing issues](https://github.com/VexoaXYZ/InkWash/issues)** - someone might have had the same problem!
3. 🐛 **[Open a new issue](https://github.com/VexoaXYZ/InkWash/issues/new)** with:
   - Your Windows version
   - InkWash version (`inkwash.exe --version`)
   - What you were trying to do
   - Error messages (full text)
   - Screenshots if possible

4. 💬 **Ask in [Discussions](https://github.com/VexoaXYZ/InkWash/discussions)**

---

## 📚 Related Guides

- [Installation Guide](Installation-Guide)
- [Creating Your First Server](Creating-Your-First-Server)
- [Converting GTA5 Mods](Converting-GTA5-Mods)
- [Command Reference](Command-Reference)

---

<div align="center">

**[⬅️ Back to Home](Home)** | **[Report a Bug 🐛](https://github.com/VexoaXYZ/InkWash/issues)**

</div>
