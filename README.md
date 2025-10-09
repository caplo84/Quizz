# 🎯 Quizz - Interactive Quiz Platform

> A modern full-stack quiz application featuring **code syntax highlighting**, **image support**, and **random question generation** from LinkedIn skill assessments.

## ✨ Key Features

- 🧠 **70+ Topics** - CSS, JavaScript, Python, Java, Go, React, and more
- 🎲 **Smart Random Mode** - Never repeat questions in the same session
- 💻 **Code Block Rendering** - Syntax-highlighted code with copy functionality
- 🖼️ **Rich Media Support** - Images with zoom and responsive display
- 🌙 **Dark/Light Mode** - Seamless theme switching
- 📱 **Mobile Responsive** - Perfect experience on all devices
- 🚀 **Real-time Sync** - Auto-sync from GitHub skill assessments

## 🚀 Quick Start

```bash
# Clone and run with Docker
git clone https://github.com/caplo84/Quizz.git
cd Quizz
docker-compose -f deployment/docker-compose.development.yml up
```

**That's it!** Access at `http://localhost:5173` 🎉

## 🛠️ Tech Stack

**Frontend:** React + Vite + Tailwind CSS + Redux Toolkit  
**Backend:** Go + Gin + PostgreSQL + Redis  
**Content:** Auto-synced from [LinkedIn Skill Assessments](https://github.com/Ebazhanov/linkedin-skill-assessments-quizzes)

## 📊 Architecture Highlights

- **Content Separation**: Code blocks and images stored separately for optimal rendering
- **Batch Processing**: Smart question batching to avoid repeats
- **Image Pipeline**: Auto-download and standardization from GitHub
- **Performance**: Redis caching + optimized database queries

## 🎮 Demo

Try the **Random Quiz Mode** - it intelligently selects questions you haven't seen yet, making each session unique!