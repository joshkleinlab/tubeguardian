# 🛡️ TubeGuardian

**TubeGuardian** is a CLI tool written in Go that automatically moderates YouTube comments using the **YouTube Data API**.  
It filters out spam, toxic, or unwanted comments based on a keyword list and hides them from your channel — helping creators maintain a safe and positive community.  

## ✨ Features
- ✅ **Custom keyword filtering** – define your own banned keywords  
- ✅ **Automated moderation** – checks new comments every 5 minutes  
- ✅ **Safe authentication** – uses Google OAuth2 for YouTube API access  
- ✅ **Cross-platform** – supports Windows, macOS, and Linux  

## 🔑 Configuration
TubeGuardian reads settings from a **YAML config file** (`configs/config.yaml`):

```yaml
channel_id: "YOUR_CHANNEL_ID"
channel_size: 50
log_dir: "configs/logs"

You can also place your keyword list in configs/banned_words.txt.
```

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