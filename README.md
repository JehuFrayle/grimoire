# 📓 Grimoire

Grimoire is a fantasy-themed personal notes web application. It allows users to create, read, update, and delete notes as if they were entries in a magical spellbook.

This project was built to showcase my skills as a fullstack developer, using:

- ⚙️ **Go** (Golang) for the REST API backend  
- 🌐 **Angular** for the single-page frontend  
- 🐘 **PostgreSQL** as the relational database  
- 🐳 **Docker** with `docker-compose` for local development  
- 📦 A well-organized monorepo structure

---

## 🔮 Features

- Full CRUD for personal notes
- Each note includes a title, content (supports markdown), and tags
- Tag-based and title-based filtering
- Optional autosave
- Intuitive UI with magical aesthetics
- Decoupled, RESTful backend
- Ready for containerized deployment

---

## 🗂️ Project structure

```

grimoire/
├── backend/       # Go REST API
├── frontend/      # Angular SPA
├── db/            # SQL migrations and seed data
├── docker/        # Dockerfiles and custom setup
├── docker-compose.yml
└── README.md

````

---

## ⚙️ Technologies used

| Layer        | Technology            |
|--------------|------------------------|
| Backend      | Go (Gin/Fiber), GORM/sqlc |
| Database     | PostgreSQL             |
| Frontend     | Angular                |
| Containers   | Docker + docker-compose |
| Extras       | JWT (optional), Markdown, Git

---

## 🚀 Getting started locally

> Make sure you have [Docker](https://www.docker.com/) installed.

```bash
git clone https://github.com/your-username/grimoire.git
cd grimoire
docker-compose up --build
````

* Frontend: [http://localhost:4200](http://localhost:4200)
* Backend: [http://localhost:8080](http://localhost:8080)
* PostgreSQL: accessible on port `5432`

---

## 📐 Roadmap

* [ ] Basic note CRUD functionality
* [ ] Angular–Go integration
* [ ] JWT-based authentication
* [ ] Shareable public note links
* [ ] Admin dashboard
* [ ] Production deployment (Railway, Render, or Netlify)

---

## ✍️ Author

Built by Jehú Frayle — [@jehufrayle](https://github.com/JehuFrayle)
This is a personal learning and portfolio project.

---

## 🧙‍♂️ License

This project is open-source and available under the MIT License.