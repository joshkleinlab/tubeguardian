# 🛡️ TubeGuardian

**TubeGuardian** is a CLI tool written in Go that automatically moderates YouTube comments using the **YouTube Data API**.  
It filters out spam, toxic, or unwanted comments based on a keyword list and hides them from your channel — helping creators maintain a safe and positive community.  

## ✨ Features
- ✅ **Custom keyword filtering** – define your own banned keywords  
- ✅ **Automated moderation** – checks new comments every 5 minutes  
- ✅ **Safe authentication** – uses Google OAuth2 for YouTube API access  
- ✅ **Cross-platform** – supports Windows, macOS, and Linux  

## 🛠️ TubeGuardian Installation Guide
TubeGuardian is a lightweight CLI tool that helps YouTube creators automatically filter spam, scams, and toxic comments.

Follow this guide to install and set it up.

### 📦 1. Download TubeGuardian
- Go to the Releases page
- Download the binary for your system:
```
tubeguardian-windows-amd64.exe → Windows
tubeguardian-darwin-arm64 → macOS (Apple Silicon)
tubeguardian-linux-amd64 → Linux
```
Move the file into a folder (e.g., ~/tubeguardian)

### ⚙️ 2. Setup Configs
Inside the project folder, create a configs/ directory with these files:
```yaml
CHANNEL_ID: "YOUR_CHANNEL_ID"
MODE_RATION: "heldForReview"
LOG_DIR: "./logs"
CREDENTIALS_FILE: "configs/credentials.json"
BANNED_WORDS_FILE: "configs/banned_words.txt"

You can also place your keyword list in configs/banned_words.txt.

credentials.json
Download this file from the Google Cloud Console (OAuth2 client secret).
```

### 🔑 3. Authenticate with YouTube
On first run, TubeGuardian will:
- Open a browser window → Google OAuth2 login
- Ask for permissions (youtube.force-ssl)
- Save your token at:
```
Linux/macOS → ./configs/token.json
Windows → ./configs/token.json
```
This step is only required once.

### ▶️ 4. Run TubeGuardian
Run the binary from the command line:
Windows
```
./tubeguardian-windows-amd64.exe
```
macOS/Linux
`
./tubeguardian-linux-amd64

`
The program will:
📥 First run → fetch all past comments & filter them
🔄 Every 5 minutes → fetch latest comments & auto-hide spam
Logs are stored in logs/.

### 🧪 5. Verify
Check your YouTube Studio → Comments → Held for review.
You should see spam/toxic comments hidden automatically.


## 📖 Usage
- On first run: TubeGuardian performs a full scan of all comments.
- Subsequent runs: TubeGuardian checks incremental new comments every 5 minutes.
- Exit: The program runs continuously until terminated manually (CTRL+C).


🔗 **Let’s connect:**
- [Email](mailto:gigacoderx@gmail.com)
- [TubeGuardian Community](https://www.reddit.com/r/TubeGuardian/)

Whether you need help with a project or want to collaborate — I'm always open to new ideas and challenges.

<a href="https://ko-fi.com/lunar_flowing">
  <img src="https://cdn.ko-fi.com/cdn/kofi3.png?v=6" width="150" alt="Buy me a coffee">
</a>

---