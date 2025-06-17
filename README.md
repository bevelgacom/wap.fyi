# WAP.FYI - The Retro Mobile URL Shortener

```
 __          ___  _____    ________     ___     __
 \ \        / / \|  __ \  |  ____\ \   / /     / /
  \ \  /\  / / _ \| |__) | | |__   \ \_/ /    / /
   \ \/  \/ / /_\ \  ___/  |  __|   \   /    / /
    \  /\  /|  _  | |      | |       | |    / /
     \/  \/ |_| |_|_|      |_|       |_|   /_/

    The Link Shortener for the New Millennium!
```

## 🌟 Welcome to WAP.FYI! 🌟

*Shorten your URLs to the power of WAP technology!*

### 📱 What is WAP.FYI?

WAP.FYI is a cutting-edge URL shortener designed for the mobile web revolution! Whether you're browsing on your Nokia 7110, Siemens S35, or that fancy new Windows Mobile device, WAP.FYI has got you covered.

**Features:**
- ✨ **WAP Compatible** - Works on ANY WAP-enabled device!
- 🔒 **Proof-of-Work Security** - Gosh the internet got dark and depressing that we need this
- 🎨 **UI** - Beautiful HTML 4.01 Transitional design with classic inset borders
- 📶 **Mobile First** - Auto-detects WAP browsers and serves WML content

### 🚀 Getting Started

Just go to [wap.fyi](http://wap.fyi) on your Windows XP machine.

Or you can run it at home...

#### Prerequisites
- Go 1.24+ (because we're living in the future!)
- A terminal (preferably green text on black background)
- Netscape Navigator 4.0+ or Internet Explorer 5.0+ for the full experience

#### Installation

1. **Clone this repository:**
   ```bash
   git clone https://github.com/bevelgacom/wap.fyi.git
   cd wap.fyi
   ```

2. **Build the application:**
   ```bash
   go build -o server ./
   ```

3. **Run the server:**
   ```bash
   ./server
   ```

4. **Visit http://localhost:8080** in your browser!

#### Docker Installation (For the Docker Revolution!)
```bash
docker build -t wap.fyi .
docker run -p 8080:8080 wap.fyi
```

### 🎮 How to Use

1. **Visit wap.fyi** - Marvel at the retro design!
1. **Enter your long URL** - No more typing on that silly keypad!
1. **Choose a custom path** - Use letters, numbers, underscores, and hyphens only
1. **Solve the challenge** - Because simple math questions aren't hard enough
1. **Get your shortened URL** - Share it with your friends over SMS!

### 📱 WAP Support

WAP.FYI automatically detects WAP browsers by checking for `text/vnd.wap.wml` in the Accept header. WAP users get:
- Custom WML 404 pages
- Automatic redirect to the Bevelgacom WAP portal
- As short as we can get without spending a million on a domain name

### 🎨 Browser Compatibility

**Tested and working on:**
- ✅ Netscape Navigator 4.0+
- ✅ Internet Explorer 5.0+
- ✅ Opera 3.5+
- ✅ Nokia WAP browsers


---

*"Best viewed in Netscape Navigator at 800x600 resolution"*

**© 1999-2025 Bevelgacom - Connecting Generations!**

---

*Visit our Geocities page for more information!*
