# Tononkira with Tones — Server

API REST en Go pour la gestion de paroles de chants (*tononkira*) avec leurs tons, auteurs et système d'authentification JWT.

## Fonctionnalités

- Authentification JWT (login par nom d'utilisateur / mot de passe)
- Liste paginée des paroles, filtrable par auteur et date
- Gestion des rôles et permissions (superAdmin, admin, etc.)
- Script de seed pour peupler MongoDB depuis des fichiers JSON statiques
- Déploiement Docker / Docker Compose

## Stack technique

- **Go 1.21+** avec [Gin](https://github.com/gin-gonic/gin)
- **MongoDB** via le driver officiel Go
- **JWT** (`golang-jwt/jwt/v5`)
- **bcrypt** pour le hachage des mots de passe

## Prérequis

- Go 1.21 ou supérieur
- MongoDB (local ou distant)
- Fichier `.env` (voir `.env.example`)

## Configuration

Copier `.env.example` vers `.env` et adapter les valeurs :

```bash
cp .env.example .env
```

Variables requises :

| Variable | Description |
|----------|-------------|
| `PORT` | Port d'écoute de l'API |
| `JWT_SECRET` | Clé secrète pour signer les tokens JWT |
| `BCRYPT_SECRET` | Réservé pour usage futur |
| `ADMIN_USERNAME` | Nom d'utilisateur admin créé par le seed |
| `ADMIN_PASSWORD` | Mot de passe admin créé par le seed |
| `DB_NAME` | Nom de la base MongoDB |
| `LOCAL_SCRIPT_DB_HOST` | Hôte MongoDB pour le script de seed |
| `DB_URI` ou `DB_USER`/`DB_PASSWORD`/`DB_HOST`/`DB_PORT` | Connexion MongoDB |

## Démarrage en local

```bash
# 1. Lancer MongoDB (via Docker)
docker compose up mongo -d

# 2. Peupler la base de données
go run scripts/seed.go

# 3. Démarrer le serveur
go run cmd/server/main.go
```

Le serveur écoute sur le port défini par `PORT` (par défaut `8080`).

### Rechargement à chaud (développement)

Le projet inclut une config [Air](https://github.com/air-verse/air) (`.air.toml`) pour le live reload.

## Démarrage avec Docker Compose

```bash
docker compose up --build
```

Le service `tononkira` charge les variables depuis `.env`.

## API

### `POST /login`

Authentification et obtention d'un token JWT.

**Corps :**
```json
{
  "userName": "admin",
  "password": "your-password"
}
```

**Réponse :**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### `GET /lyrics/list`

Liste paginée de toutes les paroles. **Authentification requise.**

**Headers :** `Authorization: Bearer <token>`

**Query params :**

| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Numéro de page (0-based) |
| `limit` | int | Nombre d'éléments par page |
| `fromDate` | date | Filtrer les paroles créées après cette date |

### `GET /lyrics/list/authors/:id`

Liste paginée des paroles d'un auteur donné. **Authentification requise.**

Mêmes query params que ci-dessus.

## Structure du projet

```
.
├── auth/           # JWT service et middleware
├── cmd/server/     # Point d'entrée
├── config/         # Configuration (.env, MongoDB)
├── domain/         # Modèles et interfaces métier
├── helpers/        # Utilitaires (auth, fichiers, dates…)
├── http/           # Handlers HTTP
├── mongodb/        # Services d'accès MongoDB
├── routes/         # Déclaration des routes
├── scripts/        # Script de seed
├── static/         # Données statiques (paroles JSON, auteurs CSV)
├── Dockerfile
└── docker-compose.yml
```

## Seed des données

Le script `scripts/seed.go` :

1. Crée les rôles (superAdmin, etc.)
2. Crée un utilisateur admin (identifiants définis dans `.env`)
3. Importe les auteurs depuis `static/authors.csv`
4. Importe les paroles depuis `static/lyrics/**/*.json`

```bash
go run scripts/seed.go
```

## CI/CD

Un workflow GitHub Actions (`.github/workflows/go.yml`) compile le projet et exécute les tests. Actuellement déclenché manuellement (`workflow_dispatch`).

## Licence

Non spécifiée.
