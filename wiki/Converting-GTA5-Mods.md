# 🚗 Converting GTA5 Mods to FiveM

InkWash makes it super easy to convert GTA5 mods from gta5-mods.com into FiveM resources! No manual downloading or converting needed.

## 🎯 What You'll Need

- ✅ InkWash installed
- ✅ At least one server created
- ✅ URLs of mods from gta5-mods.com
- ✅ Internet connection

## 🚀 Quick Start

### Step 1: Find Mods

1. Go to **[https://www.gta5-mods.com/](https://www.gta5-mods.com/)**
2. Browse or search for the mods you want (vehicles, weapons, scripts, etc.)
3. Copy the URL of the mod page

**Example URLs:**
- `https://www.gta5-mods.com/vehicles/1995-mclaren-f1-lm-addon`
- `https://www.gta5-mods.com/weapons/tactical-rifle`

### Step 2: Run the Converter

1. Open Command Prompt
2. Navigate to InkWash folder
3. Type:
   ```
   inkwash.exe convert
   ```

### Step 3: Follow the Wizard

#### Choose Your Server

The wizard will ask where to install the mods:

1. **Select an existing server** - Choose from your servers
2. **External Server (Current Directory)** - Install to `./resources/`
3. **External Server (Custom Path)** - Specify a custom path

**Tip:** For most users, just select your server from the list!

#### Add Mod URLs

Now you can add as many mod URLs as you want!

1. Paste a mod URL
2. Press **Enter** to add it to the list
3. Repeat for more mods
4. Press **Enter** on an empty line (or **Ctrl+Enter**) when done

**Example:**
```
Add GTA5 Mod URL: https://www.gta5-mods.com/vehicles/1995-mclaren-f1-lm-addon
Added: 1. https://www.gta5-mods.com/vehicles/1995-mclaren-f1-lm-addon

Add GTA5 Mod URL: https://www.gta5-mods.com/vehicles/lamborghini-huracan
Added: 2. https://www.gta5-mods.com/vehicles/lamborghini-huracan

Add GTA5 Mod URL: [Press Enter to continue]
```

### Step 4: Watch the Magic!

InkWash will now:

1. ✅ **Queue your mods** - Organizes them for processing
2. ✅ **Convert each mod** - Uses convert.cfx.rs API (max 2 at a time)
3. ✅ **Download converted files** - Parallel downloads for speed
4. ✅ **Extract to correct folders** - Automatically places in [vehicles], [weapons], etc.
5. ✅ **Clean up** - Removes temporary ZIP files

**You'll see clear progress for each mod:**
- ⏳ Queued - Waiting to start
- 🔄 Converting - Being converted right now
- ⬇️ Downloading - Being downloaded
- ✅ Complete - Ready to use!
- ✗ Failed - Something went wrong (with error message)

## 📁 Where Do Mods Go?

Mods are automatically sorted into category folders:

- **Vehicles** → `[vehicles]/mod-name/`
- **Weapons** → `[weapons]/mod-name/`
- **Scripts** → `[scripts]/mod-name/`
- **Maps** → `[maps]/mod-name/`
- **Misc** → `[misc]/mod-name/`

**Example:**
```
your-server/
└── resources/
    └── [vehicles]/
        ├── mclaren-f1-lm/
        └── lamborghini-huracan/
```

## ⚙️ Activating Mods

After converting, you need to activate the mods in your server.cfg:

1. Open `your-server/server.cfg`
2. Find the resources section
3. Add your mods:
   ```
   ensure mclaren-f1-lm
   ensure lamborghini-huracan
   ```
4. Save the file
5. Restart your server

**Tip:** Some mods might require additional setup - check the mod's original page for instructions!

## 🎯 Pro Tips

### Converting Multiple Mods

You can convert up to **10+ mods at once!** InkWash will:
- Process 2 conversions at a time (API rate limit)
- Download all converted files in parallel
- Show clear progress for each one

### Rate Limiting

InkWash respects the convert.cfx.rs API limits:
- Maximum 2 concurrent conversions
- 200ms delay between starting conversions
- Prevents 404 errors and bans

### External Mode

Use external mode to:
- Convert mods for servers not managed by InkWash
- Download to a specific folder
- Batch convert for multiple servers

## ❓ Common Issues

### "URL must be from gta5-mods.com"

**Problem:** You tried to convert from a different website.

**Solution:** Only gta5-mods.com URLs work. Find the mod on gta5-mods.com instead.

### "Conversion failed with status 404"

**Problem:** The mod might be deleted or the link is wrong.

**Solution:**
1. Check the URL is correct
2. Try visiting the URL in your browser
3. Make sure it's a direct mod page (not a category or search)

### "Download failed: unexpected status code: 404"

**Problem:** The converted file couldn't be downloaded (temporary server issue or rate limit).

**Solution:**
1. Wait a few minutes and try again
2. Don't convert too many mods at once (try 5 at a time)

### Mod doesn't appear in game

**Problem:** Mod is converted but not showing up.

**Solution:**
1. Make sure you added it to server.cfg: `ensure mod-name`
2. Check the mod's folder name in `resources/[category]/`
3. Some mods need additional setup - check the original mod page
4. Restart your server: `inkwash.exe stop server-name` then `inkwash.exe start server-name`

### "Which vehicle is processing?" UI is confusing

This was fixed in v2.0! Update to the latest version for clear progress:
- Each mod shows its name
- Clear status indicators
- Progress counter (X/Y completed)

## 🔧 Advanced Usage

### Custom Resource Paths

When using external mode with custom path:
1. Specify the full path to your resources folder
2. Make sure the path exists
3. InkWash will create category subfolders automatically

### Batch Processing

For bulk mod conversion:
1. Collect all your URLs in a text file
2. Copy and paste them one by one into the wizard
3. Let InkWash handle everything automatically

### Monitoring Progress

The converter shows:
- **Overall progress** - "Progress: 5/10 completed"
- **Queue status** - "3 queued • 2/2 active"
- **Individual status** - Each mod's current state
- **Error details** - If something fails, you'll know why

## 🎨 Mod Categories

InkWash automatically detects categories from the URL:

| Category | Folder | Example |
|----------|--------|---------|
| Vehicles | `[vehicles]/` | Cars, bikes, planes |
| Weapons | `[weapons]/` | Guns, melee weapons |
| Scripts | `[scripts]/` | Gameplay scripts |
| Maps | `[maps]/` | Map mods |
| Player | `[player]/` | Clothing, skins |
| Misc | `[misc]/` | Everything else |

---

## 🎉 You're Done!

Now you know how to convert GTA5 mods! Your server will look amazing with all these custom vehicles and weapons.

**Next steps:**
- Try different mod types (weapons, scripts, maps)
- Build your perfect server
- Share your server with friends!

---

<div align="center">

**[⬅️ Creating Your First Server](Creating-Your-First-Server)** | **[Troubleshooting ➡️](Troubleshooting)**

</div>
